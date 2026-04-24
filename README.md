# task-manager-mcp

MCP server for [doit](https://github.com/mvaldes14/task-manager) — exposes tasks, projects, subtasks, and NLP parsing as MCP tools so any MCP-compatible client (Claude Code, Claude Desktop, etc.) can interact with your self-hosted doit instance.

## Requirements

- Go 1.22+
- A running [doit](https://github.com/mvaldes14/task-manager) instance
- `DOIT_API_KEY` set in your environment

## Build

```bash
go mod tidy
go build -o task-manager-mcp .
```

## Configuration

| Env var | Default | Description |
|---|---|---|
| `DOIT_BASE_URL` | _(required)_ | Base URL of your doit instance (e.g. `http://localhost:5001`) |
| `DOIT_API_KEY` | _(required)_ | Bearer token (`TD_API_KEY` from doit's `.env`) |
| `MCP_ADDR` | `:8080` | Address the HTTP server binds to |

## Transport

The server uses the [MCP streamable HTTP transport](https://modelcontextprotocol.io/docs/concepts/transports) — plain HTTP POST, no SSE keepalive. The endpoint is `POST /mcp`.

```bash
DOIT_BASE_URL=http://localhost:5001 DOIT_API_KEY=<key> ./task-manager-mcp
# listening on :8080/mcp
```

## Claude Code setup

Add to `~/.claude/settings.json` (or your project-level `.mcp.json`):

```json
{
  "mcpServers": {
    "doit": {
      "type": "http",
      "url": "http://localhost:8080/mcp",
      "env": {
        "DOIT_BASE_URL": "http://localhost:5001",
        "DOIT_API_KEY": "your-key"
      }
    }
  }
}
```

Then restart Claude Code — the tools appear automatically.

## Tools

### Tasks

| Tool | Description |
|---|---|
| `list_tasks` | List tasks; optional filters: `project_id`, `status` (todo\|doing\|done), `search` |
| `get_task` | Get a single task by `id` |
| `create_task` | Create a task (`title` required; optional: `description`, `status`, `due_date`, `due_time`, `project_id`, `tags`, `recurrence`) |
| `update_task` | Update any fields on a task by `id` |
| `delete_task` | Delete a task by `id` |
| `get_today_tasks` | Tasks due today |
| `get_overdue_tasks` | Past-due tasks |

### Projects

| Tool | Description |
|---|---|
| `list_projects` | List all projects (includes task counts) |
| `create_project` | Create a project (`name` required; optional: `color`, `icon`) |
| `update_project` | Update a project by `id` |
| `delete_project` | Delete a project by `id` (tasks moved to inbox) |

### Subtasks

| Tool | Description |
|---|---|
| `add_subtask` | Add a subtask to a task (`task_id` + `title` required) |
| `update_subtask` | Update a subtask (`task_id`, `subtask_id` required; optional: `title`, `completed`) |
| `delete_subtask` | Delete a subtask (`task_id`, `subtask_id` required) |

### NLP

| Tool | Description |
|---|---|
| `parse_nlp` | Parse natural language (`text` required) → returns structured task fields |

### AI Results

| Tool | Description |
|---|---|
| `get_task_ai_result` | Get the stored AI-generated result for a task (`task_id` required) |
| `store_task_ai_result` | Store (upsert) an AI result for a task (`task_id` + `content` required; optional: `model`) |

## Project structure

```
.
├── main.go
└── internal/
    ├── client/
    │   └── client.go       # HTTP client with Get/Post/Patch/Delete/Put
    └── tools/
        ├── tasks.go
        ├── projects.go
        ├── subtasks.go
        ├── nlp.go
        └── ai.go
```
