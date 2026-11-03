# task-manager-mcp

MCP server for [doit](https://github.com/mvaldes14/task-manager). Exposes task-manager API operations through MCP streamable HTTP.

## Requirements

- Go 1.22+
- Running doit instance
- `DOIT_API_KEY`

## Build and run

```bash
go mod tidy
go build -o task-manager-mcp .
DOIT_BASE_URL=http://localhost:5001 \
DOIT_API_KEY=<key> \
./task-manager-mcp
```

Endpoint: `POST http://localhost:8080/mcp`

## Configuration

| Variable | Default | Purpose |
|---|---|---|
| `DOIT_BASE_URL` | required | doit base URL |
| `DOIT_API_KEY` | required | Bearer token (`TD_API_KEY` in doit) |
| `MCP_ADDR` | `:8080` | MCP HTTP listen address |

## MCP client configuration

```json
{
  "mcpServers": {
    "doit": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

Server process must receive `DOIT_BASE_URL` and `DOIT_API_KEY`.

## Tools

### Tasks

| Tool | Purpose |
|---|---|
| `list_tasks` | List/filter tasks by project, status, priority, tags, assignee, due-date range, due-date presence, search, limit, and offset |
| `search_tasks` | Type-ahead search for open tasks; results include UUID and permanent numeric `seq` |
| `get_task` | Get task by UUID, formatted key such as `DO-142`, or bare seq such as `142` |
| `create_task` | Create task; supports NLP or verbatim `skip_nlp` mode |
| `update_task` | Update task by UUID, key, or seq |
| `delete_task` | Soft-delete task; optional `purge=true` permanently deletes it |
| `restore_task` | Restore soft-deleted task |
| `reorder_tasks` | Update task positions/statuses for drag-and-drop |
| `bulk_update_tasks` | Move or change status for up to 200 tasks |
| `get_today_tasks` | List incomplete tasks due today |
| `get_overdue_tasks` | List incomplete overdue tasks |
| `get_upcoming_tasks` | List tasks due in next seven days |

Task responses include permanent numeric `seq`. Format human-readable keys as `<task_key_prefix>-<seq>`; retrieve prefix with `get_settings`.

### Projects

| Tool | Purpose |
|---|---|
| `list_projects` | List projects and subprojects |
| `create_project` | Create project or one-level subproject, with description/deadline support |
| `update_project` | Update metadata, parent, deadline, sharing, or archive state |
| `delete_project` | Delete project; tasks move to inbox and children promote to root |
| `reorder_projects` | Reorder flat depth-first project IDs within sibling groups |

### Subtasks

| Tool | Purpose |
|---|---|
| `add_subtask` | Add verbatim-title subtask or link existing task |
| `update_subtask` | Update title/completion state |
| `delete_subtask` | Delete subtask |

Parent task references accept UUID, formatted key, or bare seq.

### Context and discovery

| Tool | Purpose |
|---|---|
| `get_agent_context` | Fetch today, overdue, upcoming, projects, and settings in one call |
| `get_all_tags` | List distinct task tags |
| `get_dashboard_stats` | Fetch counts, trends, project progress, at-risk deadlines, tags, and streaks |
| `parse_nlp` | Parse natural-language task text |

### Settings

| Tool | Purpose |
|---|---|
| `get_settings` | Read settings, including `task_key_prefix` |
| `update_settings` | Merge typed common settings or arbitrary `settings_json`; supports validated task-key prefix |

`task_key_prefix` must be 2-6 uppercase letters/digits and start with a letter.

### AI results

| Tool | Purpose |
|---|---|
| `get_task_ai_result` | Read stored AI result for task |
| `store_task_ai_result` | Upsert AI result and model metadata |

### Users and authentication

| Tool | Purpose |
|---|---|
| `get_auth_status` | Read authentication/current-user state |
| `auth_login` | Log in with username/password |
| `auth_logout` | Destroy session |
| `list_users` | List users |
| `create_user` | Create user; admin only |
| `get_current_user` | Read current user |
| `update_current_user` | Update display name or password |
| `upload_current_user_avatar` | Upload base64-encoded avatar image |
| `get_user_avatar` | Fetch user avatar |

API-key deployments normally use bearer auth and do not need login tools.

### Calendars and telemetry

| Tool | Purpose |
|---|---|
| `get_gcal_status` | Read Google Calendar integration state |
| `sync_gcal` | Push open dated tasks to Google Calendar |
| `list_ics_calendars` | List ICS calendars |
| `create_ics_calendar` | Add ICS calendar by URL or uploaded file |
| `delete_ics_calendar` | Delete ICS calendar |
| `get_ics_events` | List calendar events by month |
| `send_otlp_traces` | Forward JSON or protobuf OTLP traces through doit |

## Project structure

```text
.
├── main.go
└── internal/
    ├── client/client.go
    └── tools/
        ├── tasks.go
        ├── projects.go
        ├── subtasks.go
        ├── settings.go
        ├── context.go
        └── ...
```
