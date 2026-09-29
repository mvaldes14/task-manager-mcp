package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

func RegisterCalendarTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("get_gcal_status",
			mcp.WithDescription("Get Google Calendar integration status"),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/gcal/status")
			if err != nil {
				return nil, fmt.Errorf("get_gcal_status: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("sync_gcal",
			mcp.WithDescription("Push all open tasks with due dates to Google Calendar"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Post("/api/gcal/sync", nil)
			if err != nil {
				return nil, fmt.Errorf("sync_gcal: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("list_ics_calendars",
			mcp.WithDescription("List ICS calendars"),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/ics")
			if err != nil {
				return nil, fmt.Errorf("list_ics_calendars: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_ics_calendar",
			mcp.WithDescription("Add an ICS calendar by URL or uploaded file_base64"),
			mcp.WithString("url", mcp.Description("ICS calendar URL")),
			mcp.WithString("file_base64", mcp.Description("ICS file bytes, base64 encoded")),
			mcp.WithString("filename", mcp.Description("Filename for uploaded ICS, defaults to calendar.ics")),
			mcp.WithString("name", mcp.Description("Calendar name")),
			mcp.WithString("color", mcp.Description("Calendar color")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			fields := map[string]string{}
			for _, k := range []string{"name", "color"} {
				if v, ok := args[k].(string); ok && v != "" {
					fields[k] = v
				}
			}
			if rawFile, ok := args["file_base64"].(string); ok && rawFile != "" {
				file, err := base64.StdEncoding.DecodeString(rawFile)
				if err != nil {
					return nil, fmt.Errorf("create_ics_calendar: invalid file_base64: %w", err)
				}
				filename := "calendar.ics"
				if v, ok := args["filename"].(string); ok && v != "" {
					filename = v
				}
				data, err := c.PostMultipart("/api/ics", fields, "file", filename, file)
				if err != nil {
					return nil, fmt.Errorf("create_ics_calendar: %w", err)
				}
				return mcp.NewToolResultText(prettyJSON(data)), nil
			}
			icsURL, ok := args["url"].(string)
			if !ok || icsURL == "" {
				return nil, fmt.Errorf("create_ics_calendar: url or file_base64 is required")
			}
			body := map[string]any{"url": icsURL}
			for k, v := range fields {
				body[k] = v
			}
			data, err := c.Post("/api/ics", body)
			if err != nil {
				return nil, fmt.Errorf("create_ics_calendar: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_ics_calendar",
			mcp.WithDescription("Delete an ICS calendar"),
			mcp.WithString("calendar_id", mcp.Description("ICS calendar ID"), mcp.Required()),
			mcp.WithDestructiveHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, ok := getArgs(req)["calendar_id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("delete_ics_calendar: calendar_id is required")
			}
			data, err := c.Delete("/api/ics/" + url.PathEscape(id))
			if err != nil {
				return nil, fmt.Errorf("delete_ics_calendar: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_ics_events",
			mcp.WithDescription("Get events for an ICS calendar month; defaults to current month"),
			mcp.WithString("calendar_id", mcp.Description("ICS calendar ID"), mcp.Required()),
			mcp.WithNumber("year", mcp.Description("Year")),
			mcp.WithNumber("month", mcp.Description("Month 1-12")),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			id, ok := args["calendar_id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("get_ics_events: calendar_id is required")
			}
			q := url.Values{}
			if v, ok := args["year"].(float64); ok {
				q.Set("year", strconv.Itoa(int(v)))
			}
			if v, ok := args["month"].(float64); ok {
				q.Set("month", strconv.Itoa(int(v)))
			}
			path := "/api/ics/" + url.PathEscape(id) + "/events"
			if enc := q.Encode(); enc != "" {
				path += "?" + enc
			}
			data, err := c.Get(path)
			if err != nil {
				return nil, fmt.Errorf("get_ics_events: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
