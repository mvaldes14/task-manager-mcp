package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

// RegisterSubtaskTools registers all subtask-related MCP tools onto the server.
func RegisterSubtaskTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("add_subtask",
			mcp.WithDescription("Add a subtask to an existing task. Title is stored verbatim; use linked_task_id to link an existing task."),
			mcp.WithString("task_id", mcp.Description("Parent task UUID, task key, or seq"), mcp.Required()),
			mcp.WithString("title", mcp.Description("Subtask title")),
			mcp.WithString("linked_task_id", mcp.Description("Link this subtask to an existing task; title is taken from that task")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			taskID, ok := args["task_id"].(string)
			if !ok || taskID == "" {
				return nil, fmt.Errorf("add_subtask: task_id is required")
			}
			body := map[string]any{}
			if title, ok := args["title"].(string); ok && title != "" {
				body["title"] = title
			}
			if linked, ok := args["linked_task_id"].(string); ok && linked != "" {
				body["linked_task_id"] = linked
			}
			if len(body) == 0 {
				return nil, fmt.Errorf("add_subtask: title or linked_task_id is required")
			}
			data, err := c.Post("/api/tasks/"+url.PathEscape(taskID)+"/subtasks", body)
			if err != nil {
				return nil, fmt.Errorf("add_subtask: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_subtask",
			mcp.WithDescription("Update a subtask title or completion state"),
			mcp.WithString("task_id", mcp.Description("Parent task UUID, task key, or seq"), mcp.Required()),
			mcp.WithString("subtask_id", mcp.Description("Subtask UUID"), mcp.Required()),
			mcp.WithString("title", mcp.Description("New subtask title")),
			mcp.WithBoolean("completed", mcp.Description("Mark the subtask as completed or not")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			taskID, ok := args["task_id"].(string)
			if !ok || taskID == "" {
				return nil, fmt.Errorf("update_subtask: task_id is required")
			}
			subtaskID, ok := args["subtask_id"].(string)
			if !ok || subtaskID == "" {
				return nil, fmt.Errorf("update_subtask: subtask_id is required")
			}
			body := make(map[string]any)
			if v, ok := args["title"].(string); ok && v != "" {
				body["title"] = v
			}
			if v, ok := args["completed"].(bool); ok {
				body["completed"] = v
			}
			data, err := c.Patch("/api/tasks/"+url.PathEscape(taskID)+"/subtasks/"+url.PathEscape(subtaskID), body)
			if err != nil {
				return nil, fmt.Errorf("update_subtask: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_subtask",
			mcp.WithDescription("Delete a subtask from a task"),
			mcp.WithDestructiveHintAnnotation(true),
			mcp.WithString("task_id", mcp.Description("Parent task UUID, task key, or seq"), mcp.Required()),
			mcp.WithString("subtask_id", mcp.Description("Subtask UUID"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			taskID, ok := args["task_id"].(string)
			if !ok || taskID == "" {
				return nil, fmt.Errorf("delete_subtask: task_id is required")
			}
			subtaskID, ok := args["subtask_id"].(string)
			if !ok || subtaskID == "" {
				return nil, fmt.Errorf("delete_subtask: subtask_id is required")
			}
			data, err := c.Delete("/api/tasks/" + url.PathEscape(taskID) + "/subtasks/" + url.PathEscape(subtaskID))
			if err != nil {
				return nil, fmt.Errorf("delete_subtask: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
