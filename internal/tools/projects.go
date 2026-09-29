package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

// RegisterProjectTools registers all project-related MCP tools onto the server.
func RegisterProjectTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("List all projects"),
			mcp.WithReadOnlyHintAnnotation(true),
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
			mcp.WithString("name", mcp.Description("Project name"), mcp.Required()),
			mcp.WithString("color", mcp.Description("Project color as a hex string (e.g. #ff0000)")),
			mcp.WithString("icon", mcp.Description("Project icon identifier")),
			mcp.WithString("parent_id", mcp.Description("Root project to nest under; one level only")),
			mcp.WithString("description", mcp.Description("Project description")),
			mcp.WithString("due_date", mcp.Description("Deadline as YYYY-MM-DD")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Post("/api/projects", buildProjectBody(getArgs(req)))
			if err != nil {
				return nil, fmt.Errorf("create_project: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_project",
			mcp.WithDescription("Update an existing project by ID"),
			mcp.WithString("id", mcp.Description("Project ID"), mcp.Required()),
			mcp.WithString("name", mcp.Description("Project name")),
			mcp.WithString("color", mcp.Description("Project color as a hex string")),
			mcp.WithString("icon", mcp.Description("Project icon identifier")),
			mcp.WithBoolean("shared", mcp.Description("Whether project is shared")),
			mcp.WithString("parent_id", mcp.Description("Root project id to nest, or empty/null via raw client to promote")),
			mcp.WithString("description", mcp.Description("Project description")),
			mcp.WithString("due_date", mcp.Description("Deadline as YYYY-MM-DD; empty/null via raw client clears it")),
			mcp.WithBoolean("archived", mcp.Description("true archives project and hides its tasks; false restores")),
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
			data, err := c.Patch("/api/projects/"+url.PathEscape(id), buildProjectBody(args))
			if err != nil {
				return nil, fmt.Errorf("update_project: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("delete_project",
			mcp.WithDescription("Delete a project by ID; tasks fall back to inbox and subprojects promote to top level"),
			mcp.WithString("id", mcp.Description("Project ID"), mcp.Required()),
			mcp.WithDestructiveHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, ok := getArgs(req)["id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("delete_project: id is required")
			}
			data, err := c.Delete("/api/projects/" + url.PathEscape(id))
			if err != nil {
				return nil, fmt.Errorf("delete_project: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("reorder_projects",
			mcp.WithDescription("Reorder projects. Inbox is ignored; positions are scoped per sibling group; send flat depth-first IDs."),
			mcp.WithArray("order", mcp.Description("Project IDs in desired order"), mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			order, ok := getArgs(req)["order"]
			if !ok || order == nil {
				return nil, fmt.Errorf("reorder_projects: order is required")
			}
			data, err := c.Post("/api/projects/reorder", map[string]any{"order": order})
			if err != nil {
				return nil, fmt.Errorf("reorder_projects: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}

// buildProjectBody constructs a map for project create/update payloads.
func buildProjectBody(args map[string]any) map[string]any {
	fields := []string{"name", "color", "icon", "shared", "parent_id", "description", "due_date", "archived"}
	body := make(map[string]any, len(fields))
	for _, f := range fields {
		if v, ok := args[f]; ok && v != nil {
			body[f] = v
		}
	}
	return body
}
