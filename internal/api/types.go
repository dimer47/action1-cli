package api

import "encoding/json"

// OAuthResponse is the response from the OAuth token endpoint.
type OAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// APIError is the standard error response from the API.
type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Details string `json:"details"`
}

// PageResponse is the standard paginated response wrapper.
type PageResponse struct {
	Items      []json.RawMessage `json:"items"`
	TotalCount int               `json:"total_count"`
	NextPage   string            `json:"next_page"`
	PrevPage   string            `json:"prev_page"`
}
