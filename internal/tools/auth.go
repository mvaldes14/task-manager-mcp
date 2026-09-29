package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

func RegisterAuthTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("get_auth_status",
			mcp.WithDescription("Get current auth and user state"),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/auth/status")
			if err != nil {
				return nil, fmt.Errorf("get_auth_status: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("auth_login",
			mcp.WithDescription("Log in with username/password. API-key auth is preferred for MCP; this does not persist cookies across server restarts."),
			mcp.WithString("username", mcp.Description("Username"), mcp.Required()),
			mcp.WithString("password", mcp.Description("Password"), mcp.Required()),
			mcp.WithBoolean("remember", mcp.Description("Remember login")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			username, ok := args["username"].(string)
			if !ok || username == "" {
				return nil, fmt.Errorf("auth_login: username is required")
			}
			password, ok := args["password"].(string)
			if !ok || password == "" {
				return nil, fmt.Errorf("auth_login: password is required")
			}
			body := map[string]any{"username": username, "password": password}
			if v, ok := args["remember"].(bool); ok {
				body["remember"] = v
			}
			data, err := c.Post("/auth/login", body)
			if err != nil {
				return nil, fmt.Errorf("auth_login: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("auth_logout",
			mcp.WithDescription("Destroy the current session cookie if session auth is in use"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Post("/auth/logout", nil)
			if err != nil {
				return nil, fmt.Errorf("auth_logout: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
