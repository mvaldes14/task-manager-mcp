package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

// RegisterProjectTools registers all project-related MCP tools onto the server.
func RegisterProjectTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("List all projects"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/projects")
			if err != nil {
				return nil, fmt.Errorf("list_projects: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_project",
			mcp.WithDescription("Create a new project"),
			mcp.WithString("name",
				mcp.Description("Project name"),
				mcp.Required(),
			),
			mcp.WithString("color", mcp.Description("Project color as a hex string (e.g. #ff0000)")),
			mcp.WithString("icon", mcp.Description("Project icon identifier")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			body := buildProjectBody(getArgs(req))
			data, err := c.Post("/api/projects", body)
			if err != nil {
				return nil, fmt.Errorf("create_project: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_project",
			mcp.WithDescription("Update an existing project by ID"),
			mcp.WithString("id",
				mcp.Description("Project ID"),
				mcp.Required(),
			),
			mcp.WithString("name", mcp.Description("Project name")),
			mcp.WithString("color", mcp.Description("Project color as a hex string")),
			mcp.WithString("icon", mcp.Description("Project icon identifier")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			all := getArgs(req)
			id, ok := all["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("update_project: id is required")
			}
			args := make(map[string]any)
			for k, v := range all {
				if k != "id" {
					args[k] = v
				}
			}
			body := buildProjectBody(args)
			data, err := c.Patch("/api/projects/"+id, body)
			if err != nil {
				return nil, fmt.Errorf("update_project: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_project",
			mcp.WithDescription("Delete a project by ID"),
			mcp.WithString("id",
				mcp.Description("Project ID"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, ok := getArgs(req)["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("delete_project: id is required")
			}
			data, err := c.Delete("/api/projects/" + id)
			if err != nil {
				return nil, fmt.Errorf("delete_project: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}

// buildProjectBody constructs a map for project create/update payloads.
func buildProjectBody(args map[string]any) map[string]any {
	fields := []string{"name", "color", "icon"}
	body := make(map[string]any, len(fields))
	for _, f := range fields {
		if v, ok := args[f]; ok && v != nil {
			body[f] = v
		}
	}
	return body
}
