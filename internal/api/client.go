package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dimer47/action1-cli/internal/auth"
	"github.com/dimer47/action1-cli/internal/config"
)

// Client is the Action1 API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Store      auth.Store
	Profile    string
	Verbose    bool
	token      string
}

// NewClient creates a new API client for the given region.
func NewClient(region config.Region, store auth.Store, profile string, verbose bool) *Client {
	return &Client{
		BaseURL: region.BaseURL(),
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		Store:   store,
		Profile: profile,
		Verbose: verbose,
	}
}

// Authenticate obtains an OAuth token using client credentials.
func (c *Client) Authenticate(clientID, clientSecret string) (*OAuthResponse, error) {
	form := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if c.Verbose {
		fmt.Printf("→ POST %s/oauth2/token\n", c.BaseURL)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp)
	}

	var oauthResp OAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return nil, fmt.Errorf("decoding auth response: %w", err)
	}

	c.token = oauthResp.AccessToken
	return &oauthResp, nil
}

// RefreshAuth refreshes the OAuth token using the refresh token.
func (c *Client) RefreshAuth(clientID, clientSecret, refreshToken string) (*OAuthResponse, error) {
	form := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseError(resp)
	}

	var oauthResp OAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return nil, fmt.Errorf("decoding refresh response: %w", err)
	}

	c.token = oauthResp.AccessToken
	return &oauthResp, nil
}

// EnsureAuth loads credentials and sets the token, refreshing if needed.
func (c *Client) EnsureAuth() error {
	if c.token != "" {
		return nil
	}

	creds, err := c.Store.Load(c.Profile)
	if err != nil {
		return fmt.Errorf("no credentials found — run 'action1 auth login' first")
	}

	if creds.AccessToken == "" && creds.RefreshToken == "" {
		return fmt.Errorf("no credentials found — run 'action1 auth login' first")
	}

	// Try the existing access token first
	if creds.AccessToken != "" {
		c.token = creds.AccessToken
		return nil
	}

	// Try to refresh
	if creds.RefreshToken != "" {
		oauthResp, err := c.RefreshAuth(creds.ClientID, creds.ClientSecret, creds.RefreshToken)
		if err != nil {
			return fmt.Errorf("token refresh failed — run 'action1 auth login' again: %w", err)
		}

		creds.AccessToken = oauthResp.AccessToken
		creds.RefreshToken = oauthResp.RefreshToken
		if err := c.Store.Save(c.Profile, creds); err != nil {
			return fmt.Errorf("saving refreshed credentials: %w", err)
		}

		return nil
	}

	return fmt.Errorf("no valid credentials — run 'action1 auth login' first")
}

// Get performs a GET request.
func (c *Client) Get(path string, query url.Values) (json.RawMessage, error) {
	return c.do("GET", path, query, nil)
}

// Post performs a POST request.
func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
	return c.do("POST", path, nil, body)
}

// Patch performs a PATCH request.
func (c *Client) Patch(path string, body interface{}) (json.RawMessage, error) {
	return c.do("PATCH", path, nil, body)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) (json.RawMessage, error) {
	return c.do("DELETE", path, nil, nil)
}

// Put performs a PUT request with raw body.
func (c *Client) Put(path string, body io.Reader, contentType string) (json.RawMessage, error) {
	return c.doRaw("PUT", path, body, contentType)
}

// GetAll performs paginated GET requests, collecting all results.
func (c *Client) GetAll(path string, query url.Values) ([]json.RawMessage, error) {
	var all []json.RawMessage

	for {
		raw, err := c.Get(path, query)
		if err != nil {
			return nil, err
		}

		var page PageResponse
		if err := json.Unmarshal(raw, &page); err != nil {
			// Not paginated — return single result
			all = append(all, raw)
			return all, nil
		}

		if page.Items != nil {
			all = append(all, page.Items...)
		} else {
			all = append(all, raw)
		}

		if page.NextPage == "" {
			break
		}

		// Parse next_page cursor into query params
		if query == nil {
			query = url.Values{}
		}
		query.Set("next_page", page.NextPage)
	}

	return all, nil
}

func (c *Client) do(method, path string, query url.Values, body interface{}) (json.RawMessage, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	return c.doRaw(method, path, bodyReader, "application/json")
}

func (c *Client) doRaw(method, path string, body io.Reader, contentType string) (json.RawMessage, error) {
	if err := c.EnsureAuth(); err != nil {
		return nil, err
	}

	fullURL := c.BaseURL + path
	if c.Verbose {
		fmt.Printf("→ %s %s\n", method, fullURL)
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseErrorFromBody(resp.StatusCode, respBody)
	}

	if len(respBody) == 0 {
		return json.RawMessage("{}"), nil
	}

	return json.RawMessage(respBody), nil
}

func parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	return parseErrorFromBody(resp.StatusCode, body)
}

func parseErrorFromBody(statusCode int, body []byte) error {
	var apiErr APIError
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Message != "" {
		return fmt.Errorf("API error %d: %s", statusCode, apiErr.Message)
	}
	return fmt.Errorf("API error %d: %s", statusCode, string(body))
}
