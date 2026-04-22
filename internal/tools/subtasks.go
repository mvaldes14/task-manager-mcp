package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

// RegisterSubtaskTools registers all subtask-related MCP tools onto the server.
func RegisterSubtaskTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("add_subtask",
			mcp.WithDescription("Add a subtask to an existing task"),
			mcp.WithString("task_id",
				mcp.Description("Parent task ID"),
				mcp.Required(),
			),
			mcp.WithString("title",
				mcp.Description("Subtask title"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			taskID, ok := args["task_id"].(string)
			if !ok || taskID == "" {
				return nil, fmt.Errorf("add_subtask: task_id is required")
			}
			title, ok := args["title"].(string)
			if !ok || title == "" {
				return nil, fmt.Errorf("add_subtask: title is required")
			}
			body := map[string]any{"title": title}
			data, err := c.Post("/api/tasks/"+taskID+"/subtasks", body)
			if err != nil {
				return nil, fmt.Errorf("add_subtask: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_subtask",
			mcp.WithDescription("Update a subtask on a task"),
			mcp.WithString("task_id",
				mcp.Description("Parent task ID"),
				mcp.Required(),
			),
			mcp.WithString("subtask_id",
				mcp.Description("Subtask ID"),
				mcp.Required(),
			),
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

			path := "/api/tasks/" + taskID + "/subtasks/" + subtaskID
			data, err := c.Patch(path, body)
			if err != nil {
				return nil, fmt.Errorf("update_subtask: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_subtask",
			mcp.WithDescription("Delete a subtask from a task"),
			mcp.WithString("task_id",
				mcp.Description("Parent task ID"),
				mcp.Required(),
			),
			mcp.WithString("subtask_id",
				mcp.Description("Subtask ID"),
				mcp.Required(),
			),
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

			path := "/api/tasks/" + taskID + "/subtasks/" + subtaskID
			data, err := c.Delete(path)
			if err != nil {
				return nil, fmt.Errorf("delete_subtask: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
