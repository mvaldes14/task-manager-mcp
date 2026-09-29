package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

func RegisterSettingsTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("get_settings",
			mcp.WithDescription("Get application settings, including task_key_prefix used to display task keys like DO-123"),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/settings")
			if err != nil {
				return nil, fmt.Errorf("get_settings: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_settings",
			mcp.WithDescription("Merge-patch settings. Pass a JSON object string; read-only keys like gcal_enabled are stripped by the server."),
			mcp.WithString("settings_json", mcp.Description("JSON object with settings to merge"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			raw, ok := getArgs(req)["settings_json"].(string)
			if !ok || raw == "" {
				return nil, fmt.Errorf("update_settings: settings_json is required")
			}
			var body map[string]any
			if err := json.Unmarshal([]byte(raw), &body); err != nil {
				return nil, fmt.Errorf("update_settings: invalid settings_json: %w", err)
			}
			data, err := c.Patch("/api/settings", body)
			if err != nil {
				return nil, fmt.Errorf("update_settings: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
