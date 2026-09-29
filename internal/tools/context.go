package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

// RegisterContextTools registers aggregate + discovery MCP tools onto the server.
func RegisterContextTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("get_agent_context",
			mcp.WithDescription("One-shot planning snapshot: today's tasks, overdue, next 7 days, and the project list in a single call. Prefer this over separate calls when building context."),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			sections := []struct {
				key  string
				path string
			}{
				{"today", "/api/tasks/today"},
				{"overdue", "/api/tasks/overdue"},
				{"upcoming", "/api/tasks/upcoming"},
				{"projects", "/api/projects"},
				{"settings", "/api/settings"},
			}
			out := make(map[string]json.RawMessage, len(sections))
			errs := make(map[string]string)
			for _, sec := range sections {
				data, err := c.Get(sec.path)
				if err != nil {
					errs[sec.key] = err.Error()
					continue
				}
				out[sec.key] = json.RawMessage(data)
			}
			payload := map[string]any{"context": out}
			if len(errs) > 0 {
				payload["errors"] = errs
			}
			merged, err := json.Marshal(payload)
			if err != nil {
				return nil, fmt.Errorf("get_agent_context: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(merged)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_all_tags",
			mcp.WithDescription("List all distinct tags in use across tasks; use to ground tag values instead of inventing them"),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/tags")
			if err != nil {
				return nil, fmt.Errorf("get_all_tags: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_dashboard_stats",
			mcp.WithDescription("Productivity stats over a rolling window: counts, completion trend, status breakdown, project progress, top tags, streaks"),
			mcp.WithNumber("days", mcp.Description("Window size in days (e.g. 7, 30, 90); defaults to 30")),
			mcp.WithString("today", mcp.Description("Client-supplied today date, YYYY-MM-DD, to avoid timezone issues")),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			path := "/api/dashboard/stats"
			q := url.Values{}
			args := getArgs(req)
			if v, ok := args["days"].(float64); ok {
				q.Set("days", fmt.Sprintf("%d", int(v)))
			}
			if v, ok := args["today"].(string); ok && v != "" {
				q.Set("today", v)
			}
			if enc := q.Encode(); enc != "" {
				path += "?" + enc
			}
			data, err := c.Get(path)
			if err != nil {
				return nil, fmt.Errorf("get_dashboard_stats: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
