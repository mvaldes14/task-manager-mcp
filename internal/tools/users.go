package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

func RegisterUserTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("list_users",
			mcp.WithDescription("List users"),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/users")
			if err != nil {
				return nil, fmt.Errorf("list_users: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("create_user",
			mcp.WithDescription("Create a user (admin only)"),
			mcp.WithString("username", mcp.Description("Username"), mcp.Required()),
			mcp.WithString("password", mcp.Description("Password"), mcp.Required()),
			mcp.WithString("display_name", mcp.Description("Display name")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			username, ok := args["username"].(string)
			if !ok || username == "" {
				return nil, fmt.Errorf("create_user: username is required")
			}
			password, ok := args["password"].(string)
			if !ok || password == "" {
				return nil, fmt.Errorf("create_user: password is required")
			}
			body := map[string]any{"username": username, "password": password}
			if v, ok := args["display_name"].(string); ok && v != "" {
				body["display_name"] = v
			}
			data, err := c.Post("/api/users", body)
			if err != nil {
				return nil, fmt.Errorf("create_user: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_current_user",
			mcp.WithDescription("Get the current user"),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			data, err := c.Get("/api/users/me")
			if err != nil {
				return nil, fmt.Errorf("get_current_user: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("update_current_user",
			mcp.WithDescription("Update own display name and/or password"),
			mcp.WithString("display_name", mcp.Description("Display name")),
			mcp.WithString("current_password", mcp.Description("Required when setting new_password")),
			mcp.WithString("new_password", mcp.Description("New password")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			body := map[string]any{}
			for _, k := range []string{"display_name", "current_password", "new_password"} {
				if v, ok := getArgs(req)[k].(string); ok && v != "" {
					body[k] = v
				}
			}
			data, err := c.Patch("/api/users/me", body)
			if err != nil {
				return nil, fmt.Errorf("update_current_user: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("upload_current_user_avatar",
			mcp.WithDescription("Upload the current user's avatar. Input is image bytes as base64; server stores a resized 50x50 JPEG."),
			mcp.WithString("file_base64", mcp.Description("Image file bytes, base64 encoded"), mcp.Required()),
			mcp.WithString("filename", mcp.Description("Filename, defaults to avatar.jpg")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			raw, ok := args["file_base64"].(string)
			if !ok || raw == "" {
				return nil, fmt.Errorf("upload_current_user_avatar: file_base64 is required")
			}
			file, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				return nil, fmt.Errorf("upload_current_user_avatar: invalid file_base64: %w", err)
			}
			filename := "avatar.jpg"
			if v, ok := args["filename"].(string); ok && v != "" {
				filename = v
			}
			data, err := c.PostMultipart("/api/users/me/avatar", nil, "file", filename, file)
			if err != nil {
				return nil, fmt.Errorf("upload_current_user_avatar: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)

	s.AddTool(
		mcp.NewTool("get_user_avatar",
			mcp.WithDescription("Fetch a user's avatar JPEG as base64"),
			mcp.WithString("user_id", mcp.Description("User ID"), mcp.Required()),
			mcp.WithReadOnlyHintAnnotation(true),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			id, ok := getArgs(req)["user_id"].(string)
			if !ok || id == "" {
				return nil, fmt.Errorf("get_user_avatar: user_id is required")
			}
			data, err := c.Get("/api/users/" + url.PathEscape(id) + "/avatar")
			if err != nil {
				return nil, fmt.Errorf("get_user_avatar: %w", err)
			}
			return mcp.NewToolResultText(base64.StdEncoding.EncodeToString(data)), nil
		},
	)
}
