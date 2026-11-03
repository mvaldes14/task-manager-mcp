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
			mcp.WithDescription("Merge-patch application settings. Use typed fields for common settings or settings_json for arbitrary keys/null values. Read-only keys like gcal_enabled are stripped by the server."),
			mcp.WithString("settings_json", mcp.Description("Optional JSON object with settings to merge; use null here to clear nullable values")),
			mcp.WithString("obsidian_vault", mcp.Description("Obsidian vault path")),
			mcp.WithString("obsidian_inbox", mcp.Description("Obsidian inbox path")),
			mcp.WithString("ai_webhook_url", mcp.Description("Bot assignment webhook URL; payload includes task id, seq, and formatted key")),
			mcp.WithString("otel_frontend_endpoint", mcp.Description("Frontend OTLP endpoint")),
			mcp.WithString("otel_frontend_service", mcp.Description("Frontend OpenTelemetry service name")),
			mcp.WithBoolean("reminder_enabled", mcp.Description("Enable task and project deadline reminders")),
			mcp.WithNumber("reminder_minutes_before", mcp.Description("Minutes before timed tasks to send reminders")),
			mcp.WithString("reminder_allday_time", mcp.Description("Local time for all-day reminders, HH:MM")),
			mcp.WithString("reminder_timezone", mcp.Description("IANA timezone for reminders")),
			mcp.WithString("task_key_prefix", mcp.Description("Task key prefix: 2-6 uppercase letters/digits starting with a letter, e.g. DO"), mcp.Pattern("^[A-Z][A-Z0-9]{1,5}$")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			body := map[string]any{}
			if raw, ok := args["settings_json"].(string); ok && raw != "" {
				if err := json.Unmarshal([]byte(raw), &body); err != nil {
					return nil, fmt.Errorf("update_settings: invalid settings_json: %w", err)
				}
				if body == nil {
					return nil, fmt.Errorf("update_settings: settings_json must be a JSON object")
				}
			}

			fields := []string{
				"obsidian_vault", "obsidian_inbox", "ai_webhook_url",
				"otel_frontend_endpoint", "otel_frontend_service",
				"reminder_enabled", "reminder_minutes_before", "reminder_allday_time",
				"reminder_timezone", "task_key_prefix",
			}
			for _, field := range fields {
				if value, ok := args[field]; ok && value != nil {
					body[field] = value
				}
			}
			if len(body) == 0 {
				return nil, fmt.Errorf("update_settings: at least one setting is required")
			}

			data, err := c.Patch("/api/settings", body)
			if err != nil {
				return nil, fmt.Errorf("update_settings: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
