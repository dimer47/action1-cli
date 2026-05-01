package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"github.com/dimer47/action1-cli/internal/auth"
	"github.com/dimer47/action1-cli/internal/config"
)

// NewCmdMcpServe returns the mcp-serve cobra command.
func NewCmdMcpServe() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp-serve",
		Short: "Start the MCP server (for AI assistants)",
		Long: `Starts a Model Context Protocol (MCP) server over stdio.

This allows AI assistants (Claude Code, VS Code, JetBrains) to interact
with the Action1 API through structured tools.

Example configuration for Claude Code:

  {
    "mcpServers": {
      "action1": {
        "command": "action1",
        "args": ["mcp-serve"]
      }
    }
  }`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMcpServer()
		},
		SilenceUsage: true,
	}
}

func runMcpServer() error {
	s := server.NewMCPServer(
		"action1-cli",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	registerTools(s)

	return server.ServeStdio(s)
}

// --- API client for MCP ---

type mcpClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
	store      auth.Store
	profile    string
	creds      auth.Credentials
	region     config.Region
}

func resolveClient() (*mcpClient, error) {
	cfg, err := config.Load("")
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	profile := cfg.CurrentProfile
	p := cfg.ActiveProfile()
	region := p.Region
	if region == "" {
		region = config.RegionNA
	}

	store := auth.NewStore("", false)
	creds, err := store.Load(profile)
	if err != nil {
		return nil, fmt.Errorf("not authenticated — run 'action1 auth login' first")
	}

	c := &mcpClient{
		baseURL:    region.BaseURL(),
		token:      creds.AccessToken,
		httpClient: &http.Client{Timeout: 120 * time.Second},
		store:      store,
		profile:    profile,
		creds:      creds,
		region:     region,
	}

	// If no access token but we have refresh credentials, refresh now
	if creds.AccessToken == "" {
		if err := c.refreshToken(); err != nil {
			return nil, fmt.Errorf("not authenticated — run 'action1 auth login' first")
		}
	}

	return c, nil
}

func (c *mcpClient) refreshToken() error {
	if c.creds.ClientID == "" || c.creds.ClientSecret == "" {
		return fmt.Errorf("no client credentials available")
	}

	form := url.Values{
		"client_id":     {c.creds.ClientID},
		"client_secret": {c.creds.ClientSecret},
	}

	req, err := http.NewRequest("POST", c.baseURL+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	var oauthResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&oauthResp); err != nil {
		return err
	}

	c.token = oauthResp.AccessToken
	c.creds.AccessToken = oauthResp.AccessToken
	if oauthResp.RefreshToken != "" {
		c.creds.RefreshToken = oauthResp.RefreshToken
	}

	// Persist the new token
	_ = c.store.Save(c.profile, c.creds)

	return nil
}

func (c *mcpClient) do(method, path string, body interface{}) ([]byte, error) {
	return c.doWithRetry(method, path, body, true)
}

func (c *mcpClient) doWithRetry(method, path string, body interface{}, canRetry bool) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = strings.NewReader(string(data))
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Auto-refresh on 401 and retry once
	if resp.StatusCode == http.StatusUnauthorized && canRetry {
		if refreshErr := c.refreshToken(); refreshErr == nil {
			return c.doWithRetry(method, path, body, false)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// --- Helpers ---

func resolveOrg() string {
	cfg, _ := config.Load("")
	if cfg != nil {
		p := cfg.ActiveProfile()
		return p.Org
	}
	return ""
}

func getParam(r mcp.CallToolRequest, key, defaultVal string) string {
	if v := r.GetString(key, ""); v != "" {
		return v
	}
	return defaultVal
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}
}

func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: msg,
			},
		},
	}
}

type pathBuilder func(mcp.CallToolRequest) (string, interface{})

func makeHandler(method string, buildPath pathBuilder) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		client, err := resolveClient()
		if err != nil {
			return errorResult(err.Error()), nil
		}

		path, body := buildPath(request)
		data, err := client.do(method, path, body)
		if err != nil {
			return errorResult(err.Error()), nil
		}

		// Pretty-print JSON
		var pretty interface{}
		if json.Unmarshal(data, &pretty) == nil {
			formatted, _ := json.MarshalIndent(pretty, "", "  ")
			return textResult(string(formatted)), nil
		}

		return textResult(string(data)), nil
	}
}

// --- Tool registration ---

func registerTools(s *server.MCPServer) {
	defaultOrg := resolveOrg()

	// ============================================================
	// Organizations
	// ============================================================
	s.AddTool(mcp.NewTool("org-list",
		mcp.WithDescription("List all organizations in the enterprise"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/organizations", nil
	}))

	s.AddTool(mcp.NewTool("org-create",
		mcp.WithDescription("Create a new organization"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Organization name")),
		mcp.WithString("description", mcp.Description("Organization description")),
	), makeHandler("POST", func(r mcp.CallToolRequest) (string, interface{}) {
		body := map[string]interface{}{"name": r.GetString("name", "")}
		if d := r.GetString("description", ""); d != "" {
			body["description"] = d
		}
		return "/organizations", body
	}))

	// ============================================================
	// Endpoints
	// ============================================================
	s.AddTool(mcp.NewTool("endpoint-list",
		mcp.WithDescription("List all managed endpoints in an organization"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/endpoints/managed/%s", org), nil
	}))

	s.AddTool(mcp.NewTool("endpoint-get",
		mcp.WithDescription("Get details of a specific endpoint"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("endpointId", mcp.Required(), mcp.Description("Endpoint ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/endpoints/managed/%s/%s", org, r.GetString("endpointId", "")), nil
	}))

	s.AddTool(mcp.NewTool("endpoint-status",
		mcp.WithDescription("Get endpoint status (online/offline counts)"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/endpoints/status/%s", org), nil
	}))

	s.AddTool(mcp.NewTool("endpoint-update",
		mcp.WithDescription("Update endpoint name or comment"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("endpointId", mcp.Required(), mcp.Description("Endpoint ID")),
		mcp.WithString("name", mcp.Description("New endpoint name")),
		mcp.WithString("comment", mcp.Description("New comment")),
	), makeHandler("PATCH", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		body := map[string]interface{}{}
		if n := r.GetString("name", ""); n != "" {
			body["name"] = n
		}
		if c := r.GetString("comment", ""); c != "" {
			body["comment"] = c
		}
		return fmt.Sprintf("/endpoints/managed/%s/%s", org, r.GetString("endpointId", "")), body
	}))

	s.AddTool(mcp.NewTool("endpoint-delete",
		mcp.WithDescription("Delete an endpoint"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("endpointId", mcp.Required(), mcp.Description("Endpoint ID")),
	), makeHandler("DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/endpoints/managed/%s/%s", org, r.GetString("endpointId", "")), nil
	}))

	s.AddTool(mcp.NewTool("endpoint-missing-updates",
		mcp.WithDescription("List missing updates for a specific endpoint"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("endpointId", mcp.Required(), mcp.Description("Endpoint ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/endpoints/managed/%s/%s/missing-updates", org, r.GetString("endpointId", "")), nil
	}))

	// ============================================================
	// Endpoint Groups
	// ============================================================
	s.AddTool(mcp.NewTool("endpoint-group-list",
		mcp.WithDescription("List endpoint groups"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/endpoints/groups/%s", org), nil
	}))

	s.AddTool(mcp.NewTool("endpoint-group-get",
		mcp.WithDescription("Get details of an endpoint group"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("groupId", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/endpoints/groups/%s/%s", org, r.GetString("groupId", "")), nil
	}))

	s.AddTool(mcp.NewTool("endpoint-group-members",
		mcp.WithDescription("List endpoints in a group"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("groupId", mcp.Required(), mcp.Description("Group ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/endpoints/groups/%s/%s/contents", org, r.GetString("groupId", "")), nil
	}))

	// ============================================================
	// Automations - Schedules
	// ============================================================
	s.AddTool(mcp.NewTool("automation-schedule-list",
		mcp.WithDescription("List automation schedules"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/automations/schedules/%s", org), nil
	}))

	s.AddTool(mcp.NewTool("automation-schedule-get",
		mcp.WithDescription("Get details of an automation schedule"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("automationId", mcp.Required(), mcp.Description("Automation schedule ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/automations/schedules/%s/%s", org, r.GetString("automationId", "")), nil
	}))

	s.AddTool(mcp.NewTool("automation-schedule-delete",
		mcp.WithDescription("Delete an automation schedule"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("automationId", mcp.Required(), mcp.Description("Automation schedule ID")),
	), makeHandler("DELETE", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/automations/schedules/%s/%s", org, r.GetString("automationId", "")), nil
	}))

	// ============================================================
	// Automations - Instances
	// ============================================================
	s.AddTool(mcp.NewTool("automation-instance-list",
		mcp.WithDescription("List automation instances (running or completed)"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/automations/instances/%s", org), nil
	}))

	s.AddTool(mcp.NewTool("automation-instance-get",
		mcp.WithDescription("Get details of an automation instance"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("instanceId", mcp.Required(), mcp.Description("Automation instance ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/automations/instances/%s/%s", org, r.GetString("instanceId", "")), nil
	}))

	s.AddTool(mcp.NewTool("automation-instance-results",
		mcp.WithDescription("List endpoint results for an automation instance"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("instanceId", mcp.Required(), mcp.Description("Instance ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/automations/instances/%s/%s/endpoint-results", org, r.GetString("instanceId", "")), nil
	}))

	s.AddTool(mcp.NewTool("automation-instance-stop",
		mcp.WithDescription("Stop a running automation instance"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("instanceId", mcp.Required(), mcp.Description("Instance ID")),
	), makeHandler("POST", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/automations/instances/%s/%s/stop", org, r.GetString("instanceId", "")), nil
	}))

	s.AddTool(mcp.NewTool("automation-template-list",
		mcp.WithDescription("List available action templates"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/automations/action-templates", nil
	}))

	// ============================================================
	// Reports
	// ============================================================
	s.AddTool(mcp.NewTool("report-list",
		mcp.WithDescription("List all available reports"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/reports/all", nil
	}))

	s.AddTool(mcp.NewTool("report-data",
		mcp.WithDescription("Get report data rows"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("reportId", mcp.Required(), mcp.Description("Report ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/reportdata/%s/%s/data", org, r.GetString("reportId", "")), nil
	}))

	s.AddTool(mcp.NewTool("report-export",
		mcp.WithDescription("Export a report"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("reportId", mcp.Required(), mcp.Description("Report ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/reportdata/%s/%s/export", org, r.GetString("reportId", "")), nil
	}))

	s.AddTool(mcp.NewTool("report-requery",
		mcp.WithDescription("Re-run a report to refresh data"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("reportId", mcp.Required(), mcp.Description("Report ID")),
	), makeHandler("POST", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/reportdata/%s/%s/requery", org, r.GetString("reportId", "")), nil
	}))

	// ============================================================
	// Software Repository
	// ============================================================
	s.AddTool(mcp.NewTool("software-list",
		mcp.WithDescription("List software repository packages"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/software-repository/%s", org), nil
	}))

	s.AddTool(mcp.NewTool("software-get",
		mcp.WithDescription("Get software package details"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("packageId", mcp.Required(), mcp.Description("Package ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/software-repository/%s/%s", org, r.GetString("packageId", "")), nil
	}))

	// ============================================================
	// Updates (Patches)
	// ============================================================
	s.AddTool(mcp.NewTool("update-list",
		mcp.WithDescription("List all missing updates across endpoints"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/updates/%s", org), nil
	}))

	s.AddTool(mcp.NewTool("update-get",
		mcp.WithDescription("List updates for a specific package"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("packageId", mcp.Required(), mcp.Description("Package ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/updates/%s/%s", org, r.GetString("packageId", "")), nil
	}))

	// ============================================================
	// Installed Software Inventory
	// ============================================================
	s.AddTool(mcp.NewTool("installed-software-list",
		mcp.WithDescription("List all installed software across endpoints"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/installed-software/%s/data", org), nil
	}))

	s.AddTool(mcp.NewTool("installed-software-get",
		mcp.WithDescription("List installed software on a specific endpoint"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("endpointId", mcp.Required(), mcp.Description("Endpoint ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/installed-software/%s/data/%s", org, r.GetString("endpointId", "")), nil
	}))

	// ============================================================
	// Vulnerabilities
	// ============================================================
	s.AddTool(mcp.NewTool("vulnerability-list",
		mcp.WithDescription("List all vulnerable software in the organization"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/vulnerabilities/%s", org), nil
	}))

	s.AddTool(mcp.NewTool("vulnerability-get",
		mcp.WithDescription("Get details of a specific vulnerability (CVE) in the organization"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("cveId", mcp.Required(), mcp.Description("CVE ID (e.g. CVE-2024-1234)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/vulnerabilities/%s/%s", org, r.GetString("cveId", "")), nil
	}))

	s.AddTool(mcp.NewTool("vulnerability-endpoints",
		mcp.WithDescription("List endpoints affected by a specific CVE"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("cveId", mcp.Required(), mcp.Description("CVE ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/vulnerabilities/%s/%s/endpoints", org, r.GetString("cveId", "")), nil
	}))

	s.AddTool(mcp.NewTool("vulnerability-cve",
		mcp.WithDescription("Get CVE description (global, not org-specific)"),
		mcp.WithString("cveId", mcp.Required(), mcp.Description("CVE ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return fmt.Sprintf("/CVE-descriptions/%s", r.GetString("cveId", "")), nil
	}))

	// ============================================================
	// Scripts
	// ============================================================
	s.AddTool(mcp.NewTool("script-list",
		mcp.WithDescription("List all scripts in the library"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/scripts/all", nil
	}))

	s.AddTool(mcp.NewTool("script-get",
		mcp.WithDescription("Get a specific script"),
		mcp.WithString("scriptId", mcp.Required(), mcp.Description("Script ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return fmt.Sprintf("/scripts/all/%s", r.GetString("scriptId", "")), nil
	}))

	// ============================================================
	// Data Sources
	// ============================================================
	s.AddTool(mcp.NewTool("data-source-list",
		mcp.WithDescription("List all data sources"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/data-sources/all", nil
	}))

	// ============================================================
	// Users
	// ============================================================
	s.AddTool(mcp.NewTool("user-me",
		mcp.WithDescription("Get current user information"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/me", nil
	}))

	s.AddTool(mcp.NewTool("user-list",
		mcp.WithDescription("List all users"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/users", nil
	}))

	s.AddTool(mcp.NewTool("user-get",
		mcp.WithDescription("Get user details"),
		mcp.WithString("userId", mcp.Required(), mcp.Description("User ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return fmt.Sprintf("/users/%s", r.GetString("userId", "")), nil
	}))

	s.AddTool(mcp.NewTool("user-roles",
		mcp.WithDescription("List roles assigned to a user"),
		mcp.WithString("userId", mcp.Required(), mcp.Description("User ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return fmt.Sprintf("/users/%s/roles", r.GetString("userId", "")), nil
	}))

	// ============================================================
	// Roles (RBAC)
	// ============================================================
	s.AddTool(mcp.NewTool("role-list",
		mcp.WithDescription("List all roles"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/roles", nil
	}))

	s.AddTool(mcp.NewTool("role-get",
		mcp.WithDescription("Get role details"),
		mcp.WithString("roleId", mcp.Required(), mcp.Description("Role ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return fmt.Sprintf("/roles/%s", r.GetString("roleId", "")), nil
	}))

	s.AddTool(mcp.NewTool("role-users",
		mcp.WithDescription("List users assigned to a role"),
		mcp.WithString("roleId", mcp.Required(), mcp.Description("Role ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return fmt.Sprintf("/roles/%s/users", r.GetString("roleId", "")), nil
	}))

	s.AddTool(mcp.NewTool("role-permissions",
		mcp.WithDescription("List all available permission templates"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/permissions", nil
	}))

	// ============================================================
	// Enterprise
	// ============================================================
	s.AddTool(mcp.NewTool("enterprise-get",
		mcp.WithDescription("Get enterprise settings"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/enterprise", nil
	}))

	// ============================================================
	// Subscription
	// ============================================================
	s.AddTool(mcp.NewTool("subscription-info",
		mcp.WithDescription("Get enterprise license information"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/subscription/enterprise", nil
	}))

	s.AddTool(mcp.NewTool("subscription-usage",
		mcp.WithDescription("Get enterprise usage statistics"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/subscription/usage/enterprise", nil
	}))

	// ============================================================
	// Search
	// ============================================================
	s.AddTool(mcp.NewTool("search",
		mcp.WithDescription("Search for reports, endpoints, and apps"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		q := url.QueryEscape(r.GetString("query", ""))
		return fmt.Sprintf("/search/%s?q=%s", org, q), nil
	}))

	// ============================================================
	// Audit Trail
	// ============================================================
	s.AddTool(mcp.NewTool("audit-list",
		mcp.WithDescription("List audit trail events"),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return "/audit/events", nil
	}))

	s.AddTool(mcp.NewTool("audit-get",
		mcp.WithDescription("Get a specific audit event"),
		mcp.WithString("eventId", mcp.Required(), mcp.Description("Audit event ID")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		return fmt.Sprintf("/audit/events/%s", r.GetString("eventId", "")), nil
	}))

	// ============================================================
	// Diagnostic Logs
	// ============================================================
	s.AddTool(mcp.NewTool("log-get",
		mcp.WithDescription("Get diagnostic logs"),
		mcp.WithString("orgId", mcp.Description("Organization ID (default: from config)")),
	), makeHandler("GET", func(r mcp.CallToolRequest) (string, interface{}) {
		org := getParam(r, "orgId", defaultOrg)
		return fmt.Sprintf("/logs/%s", org), nil
	}))
}
