package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

// RegisterTaskTools registers all task-related MCP tools onto the server.
func RegisterTaskTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_tasks",
			mcp.WithDescription("List tasks, optionally filtered by project, status, or search query"),
			mcp.WithString("project_id", mcp.Description("Filter by project ID")),
			mcp.WithString("status",
				mcp.Description("Filter by status"),
				mcp.Enum("todo", "doing", "done"),
			),
			mcp.WithString("search", mcp.Description("Search query to filter tasks")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			path := "/api/tasks"
			sep := "?"
			if v, ok := args["project_id"].(string); ok && v != "" {
				path += sep + "project_id=" + v
				sep = "&"
			}
			if v, ok := args["status"].(string); ok && v != "" {
				path += sep + "status=" + v
				sep = "&"
			}
			if v, ok := args["search"].(string); ok && v != "" {
				path += sep + "search=" + v
			}

			data, err := c.Get(path)
			if err != nil {
				return nil, fmt.Errorf("list_tasks: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_task",
			mcp.WithDescription("Get a single task by ID"),
			mcp.WithString("id",
				mcp.Description("Task ID"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			id, ok := args["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("get_task: id is required")
			}
			data, err := c.Get("/api/tasks/" + id)
			if err != nil {
				return nil, fmt.Errorf("get_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_task",
			mcp.WithDescription("Create a new task"),
			mcp.WithString("title",
				mcp.Description("Task title"),
				mcp.Required(),
			),
			mcp.WithString("description", mcp.Description("Task description")),
			mcp.WithString("status",
				mcp.Description("Task status"),
				mcp.Enum("todo", "doing", "done"),
			),
			mcp.WithString("due_date", mcp.Description("Due date in YYYY-MM-DD format")),
			mcp.WithString("due_time", mcp.Description("Due time in HH:MM format")),
			mcp.WithString("project_id", mcp.Description("Project ID to assign the task to")),
			mcp.WithArray("tags", mcp.Description("List of tags")),
			mcp.WithString("recurrence", mcp.Description("Recurrence rule")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			body := buildTaskBody(getArgs(req))
			data, err := c.Post("/api/tasks", body)
			if err != nil {
				return nil, fmt.Errorf("create_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_task",
			mcp.WithDescription("Update an existing task by ID"),
			mcp.WithString("id",
				mcp.Description("Task ID"),
				mcp.Required(),
			),
			mcp.WithString("title", mcp.Description("Task title")),
			mcp.WithString("description", mcp.Description("Task description")),
			mcp.WithString("status",
				mcp.Description("Task status"),
				mcp.Enum("todo", "doing", "done"),
			),
			mcp.WithString("due_date", mcp.Description("Due date in YYYY-MM-DD format")),
			mcp.WithString("due_time", mcp.Description("Due time in HH:MM format")),
			mcp.WithString("project_id", mcp.Description("Project ID")),
			mcp.WithArray("tags", mcp.Description("List of tags")),
			mcp.WithString("recurrence", mcp.Description("Recurrence rule")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			all := getArgs(req)
			id, ok := all["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("update_task: id is required")
			}
			args := make(map[string]any)
			for k, v := range all {
				if k != "id" {
					args[k] = v
				}
			}
			body := buildTaskBody(args)
			data, err := c.Patch("/api/tasks/"+id, body)
			if err != nil {
				return nil, fmt.Errorf("update_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_task",
			mcp.WithDescription("Delete a task by ID"),
			mcp.WithString("id",
				mcp.Description("Task ID"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			id, ok := args["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("delete_task: id is required")
			}
			data, err := c.Delete("/api/tasks/" + id)
			if err != nil {
				return nil, fmt.Errorf("delete_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_today_tasks",
			mcp.WithDescription("Get all tasks due today"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/tasks/today")
			if err != nil {
				return nil, fmt.Errorf("get_today_tasks: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_overdue_tasks",
			mcp.WithDescription("Get all overdue tasks"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/tasks/overdue")
			if err != nil {
				return nil, fmt.Errorf("get_overdue_tasks: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}

// buildTaskBody constructs a map for task create/update payloads from tool arguments.
func buildTaskBody(args map[string]any) map[string]any {
	fields := []string{"title", "description", "status", "due_date", "due_time", "project_id", "tags", "recurrence"}
	body := make(map[string]any, len(fields))
	for _, f := range fields {
		if v, ok := args[f]; ok && v != nil {
			body[f] = v
		}
	}
	return body
}

// prettyJSON returns a human-readable JSON string from raw JSON bytes.
// If the input is not valid JSON it is returned as-is.
func prettyJSON(raw []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return string(raw)
	}
	return buf.String()
}

// getArgs type-asserts req.Params.Arguments to map[string]any.
func getArgs(req mcp.CallToolRequest) map[string]any {
	args, _ := req.Params.Arguments.(map[string]any)
	if args == nil {
		return map[string]any{}
	}
	return args
}
