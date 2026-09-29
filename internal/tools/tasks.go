package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

// RegisterTaskTools registers all task-related MCP tools onto the server.
func RegisterTaskTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_tasks",
			mcp.WithDescription("List tasks, optionally filtered. Task search matches title/description and task keys like DO-123 or 123."),
			mcp.WithString("project_id", mcp.Description("Filter by project ID; includes subprojects")),
			mcp.WithString("status", mcp.Description("Comma-separated statuses; matches any"), mcp.Enum("todo", "doing", "blocked", "done")),
			mcp.WithString("priority", mcp.Description("Comma-separated priorities; matches any"), mcp.Enum("low", "medium", "high")),
			mcp.WithString("tag", mcp.Description("Comma-separated tags; task must have all tags")),
			mcp.WithString("assigned_to", mcp.Description("User id, me, or none")),
			mcp.WithString("due_after", mcp.Description("Inclusive lower due date bound, YYYY-MM-DD")),
			mcp.WithString("due_before", mcp.Description("Inclusive upper due date bound, YYYY-MM-DD")),
			mcp.WithBoolean("has_due_date", mcp.Description("Filter tasks by whether they have a due date")),
			mcp.WithString("search", mcp.Description("Search query; task keys like DO-123 or 123 match seq")),
			mcp.WithNumber("limit", mcp.Description("Max rows, 1-500")),
			mcp.WithNumber("offset", mcp.Description("Offset for pagination")),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			q := url.Values{}
			args := getArgs(req)
			addStringQuery(q, args, "project_id", "status", "priority", "tag", "assigned_to", "due_after", "due_before", "search")
			if v, ok := args["has_due_date"].(bool); ok {
				q.Set("has_due_date", strconv.FormatBool(v))
			}
			addNumberQuery(q, args, "limit", "offset")

			path := "/api/tasks"
			if enc := q.Encode(); enc != "" {
				path += "?" + enc
			}
			data, err := c.Get(path)
			if err != nil {
				return nil, fmt.Errorf("list_tasks: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("search_tasks",
			mcp.WithDescription("Type-ahead search for open tasks, max 8 results"),
			mcp.WithString("q", mcp.Description("Search text, at least 2 chars"), mcp.Required()),
			mcp.WithString("exclude", mcp.Description("Task ID to exclude")),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			query, ok := args["q"].(string)
			if !ok || query == "" {
				return nil, fmt.Errorf("search_tasks: q is required")
			}
			q := url.Values{"q": []string{query}}
			if v, ok := args["exclude"].(string); ok && v != "" {
				q.Set("exclude", v)
			}
			data, err := c.Get("/api/tasks/search?" + q.Encode())
			if err != nil {
				return nil, fmt.Errorf("search_tasks: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_task",
			mcp.WithDescription("Get a single task by UUID, task key like DO-142, or numeric seq like 142"),
			mcp.WithString("id", mcp.Description("Task UUID, task key, or seq"), mcp.Required()),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, ok := getArgs(req)["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("get_task: id is required")
			}
			data, err := c.Get("/api/tasks/" + url.PathEscape(id))
			if err != nil {
				return nil, fmt.Errorf("get_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_task",
			mcp.WithDescription("Create a new task. Title supports NLP unless skip_nlp is true."),
			mcp.WithString("title", mcp.Description("Task title"), mcp.Required()),
			mcp.WithString("description", mcp.Description("Task description")),
			mcp.WithString("status", mcp.Description("Task status"), mcp.Enum("todo", "doing", "blocked", "done")),
			mcp.WithString("due_date", mcp.Description("Due date in YYYY-MM-DD format")),
			mcp.WithString("due_time", mcp.Description("Due time in HH:MM format")),
			mcp.WithString("project_id", mcp.Description("Project ID to assign the task to")),
			mcp.WithArray("tags", mcp.Description("List of tags")),
			mcp.WithString("recurrence", mcp.Description("Recurrence rule")),
			mcp.WithArray("links", mcp.Description("List of links; each item an object with label and url keys")),
			mcp.WithString("priority", mcp.Description("Task priority"), mcp.Enum("low", "medium", "high")),
			mcp.WithString("assigned_to", mcp.Description("User ID to assign the task to")),
			mcp.WithString("recurrence_end", mcp.Description("Date to stop recurrence, YYYY-MM-DD")),
			mcp.WithString("timezone", mcp.Description("IANA timezone used when syncing to GCal")),
			mcp.WithBoolean("skip_nlp", mcp.Description("Take payload verbatim and require project_id to already exist")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Post("/api/tasks", buildTaskBody(getArgs(req)))
			if err != nil {
				return nil, fmt.Errorf("create_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_task",
			mcp.WithDescription("Update an existing task by UUID, task key like DO-142, or numeric seq like 142"),
			mcp.WithString("id", mcp.Description("Task UUID, task key, or seq"), mcp.Required()),
			mcp.WithString("title", mcp.Description("Task title")),
			mcp.WithString("description", mcp.Description("Task description")),
			mcp.WithString("status", mcp.Description("Task status"), mcp.Enum("todo", "doing", "blocked", "done")),
			mcp.WithString("due_date", mcp.Description("Due date in YYYY-MM-DD format; pass empty/null via raw client to clear")),
			mcp.WithString("due_time", mcp.Description("Due time in HH:MM format")),
			mcp.WithString("project_id", mcp.Description("Project ID")),
			mcp.WithArray("tags", mcp.Description("List of tags")),
			mcp.WithString("recurrence", mcp.Description("Recurrence rule")),
			mcp.WithArray("links", mcp.Description("List of links; each item an object with label and url keys")),
			mcp.WithString("priority", mcp.Description("Task priority"), mcp.Enum("low", "medium", "high")),
			mcp.WithString("assigned_to", mcp.Description("User ID to assign the task to")),
			mcp.WithString("recurrence_end", mcp.Description("Date to stop recurrence, YYYY-MM-DD")),
			mcp.WithNumber("position", mcp.Description("Task position")),
			mcp.WithString("timezone", mcp.Description("IANA timezone")),
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
			data, err := c.Patch("/api/tasks/"+url.PathEscape(id), buildTaskBody(args))
			if err != nil {
				return nil, fmt.Errorf("update_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_task",
			mcp.WithDescription("Delete a task by UUID, task key, or seq. Soft-delete by default; set purge=true for permanent delete."),
			mcp.WithString("id", mcp.Description("Task UUID, task key, or seq"), mcp.Required()),
			mcp.WithBoolean("purge", mcp.Description("Hard delete permanently instead of soft-delete")),
			mcp.WithDestructiveHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			id, ok := args["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("delete_task: id is required")
			}
			path := "/api/tasks/" + url.PathEscape(id)
			if purge, ok := args["purge"].(bool); ok && purge {
				path += "?purge=true"
			}
			data, err := c.Delete(path)
			if err != nil {
				return nil, fmt.Errorf("delete_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("restore_task",
			mcp.WithDescription("Restore a soft-deleted task by UUID, task key, or seq"),
			mcp.WithString("id", mcp.Description("Task UUID, task key, or seq"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, ok := getArgs(req)["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("restore_task: id is required")
			}
			data, err := c.Post("/api/tasks/"+url.PathEscape(id)+"/restore", nil)
			if err != nil {
				return nil, fmt.Errorf("restore_task: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("reorder_tasks",
			mcp.WithDescription("Bulk update task positions and statuses for drag-and-drop"),
			mcp.WithArray("tasks", mcp.Description("Array of objects: {id, position, status}"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			tasks, ok := getArgs(req)["tasks"]
			if !ok || tasks == nil {
				return nil, fmt.Errorf("reorder_tasks: tasks is required")
			}
			data, err := c.Post("/api/tasks/reorder", tasks)
			if err != nil {
				return nil, fmt.Errorf("reorder_tasks: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("bulk_update_tasks",
			mcp.WithDescription("Apply project_id and/or status updates to up to 200 tasks"),
			mcp.WithArray("ids", mcp.Description("Task ids/keys to update, 1-200"), mcp.Required()),
			mcp.WithString("project_id", mcp.Description("Project ID to move tasks to")),
			mcp.WithString("status", mcp.Description("Status to set"), mcp.Enum("todo", "doing", "blocked", "done")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			ids, ok := args["ids"]
			if !ok || ids == nil {
				return nil, fmt.Errorf("bulk_update_tasks: ids is required")
			}
			updates := map[string]any{}
			if v, ok := args["project_id"].(string); ok && v != "" {
				updates["project_id"] = v
			}
			if v, ok := args["status"].(string); ok && v != "" {
				updates["status"] = v
			}
			if len(updates) == 0 {
				return nil, fmt.Errorf("bulk_update_tasks: project_id or status is required")
			}
			data, err := c.Patch("/api/tasks/bulk", map[string]any{"ids": ids, "updates": updates})
			if err != nil {
				return nil, fmt.Errorf("bulk_update_tasks: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	addSimpleTaskListTool(s, c, "get_today_tasks", "Get all tasks due today", "/api/tasks/today")
	addSimpleTaskListTool(s, c, "get_overdue_tasks", "Get all overdue tasks", "/api/tasks/overdue")
	addSimpleTaskListTool(s, c, "get_upcoming_tasks", "Get tasks due within the next 7 days", "/api/tasks/upcoming")
}

func addSimpleTaskListTool(s *server.MCPServer, c *client.Client, name, desc, path string) {
	s.AddTool(mcp.NewTool(name, mcp.WithDescription(desc), mcp.WithReadOnlyHintAnnotation(true)), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		data, err := c.Get(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return mcp.NewToolResultText(prettyJSON(data)), nil
	})
}

// buildTaskBody constructs a map for task create/update payloads from tool arguments.
func buildTaskBody(args map[string]any) map[string]any {
	fields := []string{"title", "description", "status", "due_date", "due_time", "project_id", "tags", "recurrence", "recurrence_end", "links", "priority", "assigned_to", "position", "timezone", "skip_nlp"}
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
	if len(raw) == 0 {
		return "OK"
	}
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

func addStringQuery(q url.Values, args map[string]any, keys ...string) {
	for _, k := range keys {
		if v, ok := args[k].(string); ok && v != "" {
			q.Set(k, v)
		}
	}
}

func addNumberQuery(q url.Values, args map[string]any, keys ...string) {
	for _, k := range keys {
		if v, ok := args[k].(float64); ok {
			q.Set(k, strconv.Itoa(int(v)))
		}
	}
}
