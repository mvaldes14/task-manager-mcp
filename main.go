package main

import (
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/mvaldes14/task-manager-mcp/internal/client"
	"github.com/mvaldes14/task-manager-mcp/internal/tools"
)

func main() {
	baseURL := os.Getenv("DOIT_BASE_URL")
	if baseURL == "" {
		log.Fatal("DOIT_BASE_URL is not set")
	}

	apiKey := os.Getenv("DOIT_API_KEY")
	if apiKey == "" {
		log.Fatal("DOIT_API_KEY is not set")
	}

	addr := os.Getenv("MCP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	mcpURL := os.Getenv("MCP_BASE_URL")
	if mcpURL == "" {
		mcpURL = "http://localhost" + addr
	}

	c := client.New(baseURL, apiKey)

	s := server.NewMCPServer("doit-mcp", "0.1.0")

	tools.RegisterTaskTools(s, c)
	tools.RegisterProjectTools(s, c)
	tools.RegisterNLPTools(s, c)
	tools.RegisterSubtaskTools(s, c)

	sse := server.NewSSEServer(s, server.WithBaseURL(mcpURL))
	log.Printf("MCP SSE server listening on %s", addr)
	if err := sse.Start(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
