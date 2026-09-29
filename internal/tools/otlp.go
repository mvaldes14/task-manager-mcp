package tools

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
)

func RegisterOTLPTools(s *server.MCPServer, c *client.Client) {
	s.AddTool(
		mcp.NewTool("send_otlp_traces",
			mcp.WithDescription("Proxy OTLP/HTTP traces through the backend. Provide either json_payload or protobuf_base64."),
			mcp.WithString("json_payload", mcp.Description("OTLP JSON payload")),
			mcp.WithString("protobuf_base64", mcp.Description("OTLP protobuf payload, base64 encoded")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := getArgs(req)
			if raw, ok := args["json_payload"].(string); ok && raw != "" {
				data, err := c.PostBytes("/api/otlp/v1/traces", "application/json", []byte(raw))
				if err != nil {
					return nil, fmt.Errorf("send_otlp_traces: %w", err)
				}
				return mcp.NewToolResultText(prettyJSON(data)), nil
			}
			if raw, ok := args["protobuf_base64"].(string); ok && raw != "" {
				body, err := base64.StdEncoding.DecodeString(raw)
				if err != nil {
					return nil, fmt.Errorf("send_otlp_traces: invalid protobuf_base64: %w", err)
				}
				data, err := c.PostBytes("/api/otlp/v1/traces", "application/x-protobuf", body)
				if err != nil {
					return nil, fmt.Errorf("send_otlp_traces: %w", err)
				}
				return mcp.NewToolResultText(prettyJSON(data)), nil
			}
			return nil, fmt.Errorf("send_otlp_traces: json_payload or protobuf_base64 is required")
		},
	)
}
