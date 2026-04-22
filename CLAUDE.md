# task-manager-mcp

MCP server for the doit task manager. Written in Go using `github.com/mark3labs/mcp-go`.

## Build & run

```bash
go mod tidy
go build -o task-manager-mcp .
DOIT_BASE_URL=http://localhost:5001 DOIT_API_KEY=<key> ./task-manager-mcp
```

## Architecture

- `main.go` — reads env, wires client + server, calls all Register* functions, serves stdio
- `internal/client/client.go` — thin HTTP client; all methods return `([]byte, error)`; non-2xx = error
- `internal/tools/*.go` — one file per domain; each exposes a `Register*Tools(s *server.MCPServer, c *client.Client)` function
- `prettyJSON` helper lives in `tools/tasks.go` and is shared across the package

## Adding a new tool

1. Pick the right file in `internal/tools/` (or create one for a new domain)
2. Add an `s.AddTool(mcp.NewTool(...), handler)` call inside the `Register*Tools` function
3. If it's a new file, add `tools.Register*Tools(s, c)` in `main.go`

## Conventions

- Tool names use `snake_case`
- Required params get `mcp.Required()`
- Handlers return raw API JSON, pretty-printed via `prettyJSON`
- Errors are wrapped with `fmt.Errorf("tool_name: %w", err)`
- No global state — client is passed explicitly

## Environment vars

| Var | Default | Purpose |
|---|---|---|
| `DOIT_BASE_URL` | `http://localhost:5001` | doit instance URL |
| `DOIT_API_KEY` | _(none)_ | Bearer token for auth |
