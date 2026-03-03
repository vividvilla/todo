# todo — CLI-Based TODO Manager for Local AI Workflows

A lightweight, Go-based CLI TODO manager designed for local AI workflows — with a background daemon that sends push notifications (via [ntfy](https://ntfy.sh)) and webhooks when tasks are due. All data is stored locally in SQLite.

## Features

- ✅ **Full CRUD** — Add, list, show, edit, mark done, and delete TODOs
- 🔔 **Multi-channel Notifications** — ntfy push notifications + generic webhook postback
- 📦 **SQLite Storage** — No external database needed, everything stored in `~/.local/share/todo/todo.db`
- 🎨 **Pretty Output** — Color-coded priorities and statuses with table formatting
- 🤖 **AI-Agent Ready** — Structured CLI output, simple commands, and a pre-defined agent skill (`@skills/todo/SKILL.md`) make it a perfect tool for AI-driven workflows
- 🐳 **Docker Support** — Run the daemon as a container with shared SQLite volume
- 🚀 **Single Binary** — Pure Go, no CGO, compiles to a single binary

## Built for AI Workflows

`todo` is designed to work seamlessly as a tool in AI agent pipelines. Whether you're using [pi](https://github.com/badlogic/pi), Claude Code, Codex, or your own agent harness:

- **Deterministic CLI** — Every command has predictable output, easy for agents to parse
- **Simple CRUD interface** — `add`, `list`, `show`, `edit`, `done`, `delete` — no ambiguity
- **Filterable output** — Agents can query by status (`-s pending`) or priority (`-p high`) to focus on what matters
- **Notification daemon** — Set it and forget it; the daemon handles alerting humans while agents manage the task lifecycle
- **Included Agent Skill** — The `@skills/todo/SKILL.md` teaches any compatible agent how to use `todo` expertly, including natural date parsing, priority inference, and multi-step workflows
- **Webhook integration** — The postback channel lets agents create tasks that notify external systems (Slack bots, CI pipelines, dashboards) when due

**Example AI workflow:**
```
User: "Remind me to review the PR before Friday's standup"
Agent: todo add "Review PR" -p high -D "2026-03-06 09:00"
       → Daemon sends ntfy push to your phone at 9 AM Friday
       → Postback hits your Slack bot with the reminder
```

Drop `@skills/todo/` into your agent's skills directory and it just works.

**Install via [vercel-labs/skills](https://github.com/vercel-labs/skills):**
```bash
# Install the todo skill into your agent's skills directory
npx @vercel-labs/skills add https://github.com/vividvilla/todo.git/skills
```

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/vividvilla/todo.git
cd todo

# Build
go build -o todo .

# (Optional) Install to your PATH
cp todo ~/.local/bin/

# Or use just
just install
```

### Requirements

- Go 1.24+ (for building)

## Quick Start

```bash
# Add your first TODO
todo add "Buy groceries" -p high -D "2026-03-04T10:00:00Z"

# Add more TODOs
todo add "Write documentation" -d "README and API docs" -p medium
todo add "Fix login bug" -p high -D "2026-03-04"

# List all active TODOs
todo list

# Show details
todo show 1

# Mark as done
todo done 1

# Set up notifications (one or both)
todo config set ntfy_url "https://ntfy.sh/your_topic"
todo config set postback_url "https://your-server.com/webhook"

# Start the notification daemon
todo daemon
```

## Commands

### `todo add <title> [flags]`

Add a new TODO.

```bash
todo add "Deploy v2.0"
todo add "Review PR" -p high -D "2026-03-05T09:00:00Z"
todo add "Write tests" -d "Unit tests for auth module" -p medium
```

**Flags:**

| Flag               | Short | Default  | Description                      |
|--------------------|-------|----------|----------------------------------|
| `--description`    | `-d`  |          | Description text                 |
| `--priority`       | `-p`  | `medium` | Priority: `low`, `medium`, `high`|
| `--due`            | `-D`  |          | Due date/time (see date formats) |

### `todo list [flags]`

List TODOs. By default, hides completed tasks.

```bash
todo list              # Active TODOs (pending + in-progress)
todo list -a           # All TODOs including done
todo list -p high      # Only high priority
todo list -s pending   # Only pending
todo ls                # Alias for list
```

**Flags:**

| Flag         | Short | Description                                       |
|--------------|-------|---------------------------------------------------|
| `--status`   | `-s`  | Filter: `pending`, `in-progress`, `done`          |
| `--priority` | `-p`  | Filter: `low`, `medium`, `high`                   |
| `--all`      | `-a`  | Show all including done                            |

**Output Example:**

```
╭────┬────────────────────┬──────────┬─────────────┬──────────────────────╮
│ ID │ TITLE              │ PRIORITY │ STATUS      │ DUE                  │
├────┼────────────────────┼──────────┼─────────────┼──────────────────────┤
│  3 │ Fix login bug      │ high     │ pending     │ Mar 4, 2026 12:00 AM │
│  1 │ Buy groceries      │ high     │ pending     │ Mar 4, 2026 3:30 PM  │
│  2 │ Write docs         │ medium   │ in-progress │ -                    │
╰────┴────────────────────┴──────────┴─────────────┴──────────────────────╯
```

### `todo show <id>`

Show full details of a TODO.

```bash
todo show 1
```

**Output:**

```
TODO #1
─────────────────────────────────
  Title:       Buy groceries
  Description: -
  Priority:    high
  Status:      pending
  Due:         Wed Mar 4, 2026 3:30 PM
  Notified:    false
  Created:     Wed Mar 4, 2026 12:38 AM
  Updated:     Wed Mar 4, 2026 12:38 AM
```

### `todo edit <id> [flags]`

Edit an existing TODO. Only specified flags are updated.

```bash
todo edit 1 -t "Buy organic groceries"    # Change title
todo edit 1 -p low                         # Change priority
todo edit 1 -s in-progress                 # Change status
todo edit 1 -D "2026-03-10"               # Change due date
todo edit 1 -D ""                          # Clear due date
```

**Flags:**

| Flag            | Short | Description                            |
|-----------------|-------|----------------------------------------|
| `--title`       | `-t`  | New title                              |
| `--description` | `-d`  | New description                        |
| `--priority`    | `-p`  | New priority: `low`, `medium`, `high`  |
| `--status`      | `-s`  | New status: `pending`, `in-progress`, `done` |
| `--due`         | `-D`  | New due date (empty string to clear)   |

### `todo done <id>`

Mark a TODO as done.

```bash
todo done 1
# ✓ TODO #1 marked as done
```

### `todo delete <id>`

Delete a TODO permanently.

```bash
todo delete 3
todo rm 3        # Alias
```

### `todo config set <key> <value>`

Set a configuration value.

```bash
todo config set ntfy_url "https://ntfy.sh/your_topic"
todo config set postback_url "https://your-server.com/hooks/todo"
```

**Supported config keys:**

| Key             | Description                                      |
|-----------------|--------------------------------------------------|
| `postback_url`  | Generic webhook URL (JSON POST)                  |
| `ntfy_url`      | ntfy topic URL (e.g. `https://ntfy.sh/my_topic`) |

### `todo config get <key>`

Get a configuration value.

```bash
todo config get ntfy_url
# ntfy_url: https://ntfy.sh/your_topic
```

### `todo daemon [flags]`

Start the notification daemon. It periodically checks for overdue TODOs and sends notifications via all configured channels (postback, ntfy).

```bash
todo daemon              # Check every 60 seconds (default)
todo daemon -i 30        # Check every 30 seconds
```

**Flags:**

| Flag         | Short | Default | Description              |
|--------------|-------|---------|--------------------------|
| `--interval` | `-i`  | `60`    | Check interval (seconds) |

## Date Formats

The `--due` / `-D` flag accepts the following date formats:

| Format                       | Example                    |
|------------------------------|----------------------------|
| RFC3339                      | `2026-03-04T15:00:00Z`     |
| RFC3339 with offset          | `2026-03-04T20:30:00+05:30`|
| Date only (midnight local)   | `2026-03-04`               |
| Datetime without TZ (local)  | `2026-03-04T15:00:00`      |
| Date + time (local)          | `2026-03-04 15:00`         |

## Notification Daemon

### How It Works

1. The daemon runs a ticker loop (default: every 60 seconds)
2. On each tick, it queries the database for TODOs where:
   - `due_at` is in the past
   - `notified` is `false`
   - `status` is not `done`
3. For each overdue TODO, it sends notifications via **all configured channels**
4. If **at least one** channel delivers successfully (HTTP 2xx), the TODO is marked as `notified`
5. On failure, each channel retries once after 2 seconds, then skips until the next cycle

### Notification Channels

| Channel      | Config Key       | Description                              |
|--------------|------------------|------------------------------------------|
| **Postback** | `postback_url`   | Generic JSON webhook (POST)              |
| **ntfy**     | `ntfy_url`       | Push notifications via [ntfy.sh](https://ntfy.sh) |

You can configure one or both. They fire independently.

### ntfy Channel

[ntfy](https://ntfy.sh) is a simple pub-sub notification service. Install the ntfy app on your phone, subscribe to your topic, and get push notifications when TODOs are overdue.

```bash
todo config set ntfy_url "https://ntfy.sh/your_topic"
```

The ntfy notification includes:
- **Title:** `⏰ TODO #1 Overdue`
- **Markdown body:** Task title, description, priority, status, due date
- **Priority mapping:** `high` → urgent (5), `medium` → default (3), `low` → low (2)
- **Emoji tags:** 🚨 for high, ⚠️ for medium, ℹ️ for low

### Postback Channel

Sends a JSON POST to any webhook URL:

```bash
todo config set postback_url "https://your-server.com/webhook"
```

**Payload format:**

```json
{
  "event": "todo_due",
  "todo": {
    "id": 3,
    "title": "Fix login bug",
    "description": "Authentication failing on mobile",
    "priority": "high",
    "status": "pending",
    "due_at": "2026-03-04T00:00:00Z",
    "notified": false,
    "created_at": "2026-03-03T19:08:00Z",
    "updated_at": "2026-03-03T19:08:00Z"
  },
  "timestamp": "2026-03-04T00:01:00Z"
}
```

### Setup

```bash
# 1. Configure one or more notification channels
todo config set ntfy_url "https://ntfy.sh/your_topic"
todo config set postback_url "https://your-server.com/webhook"

# 2. Start the daemon
todo daemon

# 3. (Optional) Run as a systemd service — see below
```

### Running as a Systemd Service

Create `~/.config/systemd/user/todo.service`:

```ini
[Unit]
Description=todo notification daemon
After=network.target

[Service]
ExecStart=%h/.local/bin/todo daemon
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
```

Then:

```bash
systemctl --user daemon-reload
systemctl --user enable todo
systemctl --user start todo
systemctl --user status todo

# View logs
journalctl --user -u todo -f
```

## Docker

### Running the Daemon with Docker Compose

The daemon can run as a container while you use the CLI on the host. They share the same SQLite database via a volume mount.

```bash
# 1. Configure notification channels (from host CLI)
todo config set ntfy_url "https://ntfy.sh/your_topic"
# and/or
todo config set postback_url "https://your-server.com/webhook"

# 2. Start the daemon container
docker compose up -d

# 3. View logs
docker compose logs -f

# 4. Stop the daemon
docker compose down
```

### Environment Variables

| Variable      | Default                     | Description              |
|---------------|-----------------------------|--------------------------|
| `TODO_DATA`   | `~/.local/share/todo`       | Path to the data directory (SQLite DB) |
| `TZ`          | `Asia/Kolkata`              | Timezone for the container |

### Custom Data Path

```bash
TODO_DATA=/path/to/data docker compose up -d
```

### Build Only

```bash
docker build -t todo .
docker run --rm -v ~/.local/share/todo:/home/todo/.local/share/todo todo daemon
```

## Data Storage

All data is stored in `~/.local/share/todo/`:

```
~/.local/share/todo/
├── todo.db    # SQLite database
└── todo.pid   # PID file (when daemon is running)
```

## Extending Notifications

The notification system uses a `Notifier` interface, making it easy to add new channels:

```go
type Notifier interface {
    Name() string
    Send(todo model.Todo) error
}
```

**Built-in channels:** Postback (webhook), ntfy

**Planned / easy to add:**

- **Email** — Add an SMTP notifier in `notify/`
- **Slack** — Add a Slack webhook notifier
- **Desktop** — Add `notify-send` / OS notification support
- **Telegram** — Add a Telegram bot notifier

To add a new channel: implement `Notifier`, then register it in `buildNotifiers()` in `cmd/daemon.go`.

## License

MIT
