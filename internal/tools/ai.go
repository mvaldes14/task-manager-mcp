package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

func RegisterAITools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("get_task_ai_result",
			mcp.WithDescription("Get the stored AI-generated result for a task"),
			mcp.WithString("task_id",
				mcp.Description("Task ID"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			taskID, ok := args["task_id"].(string)
			if !ok || taskID == "" {
				return nil, fmt.Errorf("get_task_ai_result: task_id is required")
			}
			data, err := c.Get("/api/tasks/" + taskID + "/ai")
			if err != nil {
				return nil, fmt.Errorf("get_task_ai_result: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("store_task_ai_result",
			mcp.WithDescription("Store (upsert) an AI-generated result for a task"),
			mcp.WithString("task_id",
				mcp.Description("Task ID"),
				mcp.Required(),
			),
			mcp.WithString("content",
				mcp.Description("AI-generated text to store"),
				mcp.Required(),
			),
			mcp.WithString("model",
				mcp.Description("Model identifier e.g. claude-sonnet-4-6"),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			taskID, ok := args["task_id"].(string)
			if !ok || taskID == "" {
				return nil, fmt.Errorf("store_task_ai_result: task_id is required")
			}
			content, ok := args["content"].(string)
			if !ok || content == "" {
				return nil, fmt.Errorf("store_task_ai_result: content is required")
			}

			body := map[string]any{"content": content}
			if model, ok := args["model"].(string); ok && model != "" {
				body["model"] = model
			}

			data, err := c.Put("/api/tasks/"+taskID+"/ai", body)
			if err != nil {
				return nil, fmt.Errorf("store_task_ai_result: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
