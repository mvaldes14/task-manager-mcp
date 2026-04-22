package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

// RegisterNLPTools registers natural-language-processing MCP tools onto the server.
func RegisterNLPTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("parse_nlp",
			mcp.WithDescription("Parse a natural language string into structured task fields"),
			mcp.WithString("text",
				mcp.Description("Natural language text describing a task (e.g. 'Buy groceries tomorrow at 3pm')"),
				mcp.Required(),
			),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			text, ok := getArgs(req)["text"].(string)
			if !ok || text == "" {
				return nil, fmt.Errorf("parse_nlp: text is required")
			}
			body := map[string]any{"text": text}
			data, err := c.Post("/api/nlp/parse", body)
			if err != nil {
				return nil, fmt.Errorf("parse_nlp: %w", err)
			}
			return mcp.NewToolResultText(prettyJSON(data)), nil
		},
	)
}
