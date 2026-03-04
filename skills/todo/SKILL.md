---
name: todo
description: Manages TODO tasks using the todo CLI. Use this skill when the user wants to create, list, update, complete, delete, or search TODOs, manage task priorities and due dates, check overdue items, or configure the todo notification daemon. Handles all task management workflows.
---

# todo Task Manager

You are an expert at managing TODO tasks using the `todo` CLI tool.

## Setup

Before running any commands, verify `todo` is available:

```bash
which todo || todo version
```

If `todo` is **not found**, help the user install it:

1. **From GitHub releases (recommended):**
   ```bash
   # Determine platform
   OS=$(uname -s | tr '[:upper:]' '[:lower:]')
   ARCH=$(uname -m); [ "$ARCH" = "x86_64" ] && ARCH="amd64"; [ "$ARCH" = "aarch64" ] && ARCH="arm64"
   VERSION=$(curl -s https://api.github.com/repos/vivek/todod/releases/latest | grep tag_name | cut -d '"' -f 4)

   # Download and install
   curl -L "https://github.com/vivek/todod/releases/download/${VERSION}/todo-${VERSION}-${OS}-${ARCH}.tar.gz" -o /tmp/todo.tar.gz
   tar xzf /tmp/todo.tar.gz -C /tmp todo
   sudo mv /tmp/todo /usr/local/bin/todo
   rm /tmp/todo.tar.gz
   ```

2. **Build from source (requires Go 1.24+):**
   ```bash
   git clone https://github.com/vivek/todod.git
   cd todod
   go build -o todo .
   sudo mv todo /usr/local/bin/todo
   ```

After installation, confirm it works: `todo version`

## Core Principles

1. **Always confirm destructive actions** — before deleting tasks, ask the user to confirm unless they explicitly said "delete".
2. **Show context after mutations** — after adding, editing, completing, or deleting a task, run `list` so the user sees the current state.
3. **Use appropriate priorities** — if the user's language implies urgency (e.g. "urgent", "ASAP", "critical", "blocker"), use `-p high`. If casual (e.g. "whenever", "someday", "nice to have"), use `-p low`. Default to `-p medium`.
4. **Parse natural dates** — convert the user's natural language dates into the formats todo accepts. See the date format reference below.
5. **Be proactive** — if the user adds a task without a due date but the context implies one (e.g. "before Friday's meeting"), suggest setting one.

## Command Reference

### Add a task
```bash
todo add "<title>" [-d "<description>"] [-p <low|medium|high>] [-D "<due_date>"]
```

Examples:
```bash
todo add "Deploy v2.0 to production" -p high -D "2026-03-10T14:00:00Z"
todo add "Write unit tests for auth" -d "Cover login, logout, and token refresh" -p medium
todo add "Read Go blog post" -p low
```

### List tasks
```bash
todo list                    # Active tasks (pending + in-progress), sorted by priority then due date
todo list -a                 # All tasks including done
todo list -p high            # Only high priority
todo list -s in-progress     # Only in-progress
todo list -s pending -p high # Pending + high priority
```

### Show task details
```bash
todo show <id>
```

### Edit a task
```bash
todo edit <id> [-t "<new_title>"] [-d "<new_description>"] [-p <priority>] [-s <status>] [-D "<due_date>"]
todo edit <id> -D ""         # Clear the due date
```

### Mark a task as done
```bash
todo done <id>
```

### Delete a task
```bash
todo delete <id>
```

### Configuration
```bash
todo config set ntfy_url "<ntfy_topic_url>"    # ntfy push notifications
todo config set postback_url "<webhook_url>"    # generic webhook
todo config get ntfy_url
todo config get postback_url
```

### Daemon
```bash
todo daemon                  # Start daemon (checks every 60s)
todo daemon -i 30            # Check every 30 seconds
```

The daemon can also be run via Docker:
```bash
# From the cloned todod repo directory:
docker compose up -d    # Start
docker compose logs -f  # View logs
docker compose down     # Stop
```

## Date Formats

todo accepts these date formats for the `-D` flag:

| User says | Convert to | Flag value |
|-----------|-----------|------------|
| "tomorrow" | Next day, midnight local | `2026-03-05` (calculate from current date) |
| "next Monday" | The coming Monday | `2026-03-09` (calculate from current date) |
| "March 10" | That date | `2026-03-10` |
| "March 10 at 2pm" | Date + time local | `2026-03-10 14:00` |
| "in 3 days" | Current date + 3 | Calculate and use `YYYY-MM-DD` |
| "end of week" | Coming Friday | Calculate and use `YYYY-MM-DD` |
| "end of month" | Last day of current month | Calculate and use `YYYY-MM-DD` |
| Explicit ISO | Pass through | `2026-03-10T14:00:00Z` |

Always calculate the actual date from the current date/time before passing to todo. Use `date` command if needed to compute relative dates:
```bash
date -d "+3 days" +%Y-%m-%d          # 3 days from now
date -d "next friday" +%Y-%m-%d      # Next Friday
```

## Workflows

### When the user wants to add a task

1. Extract title, priority, due date, and description from their message
2. Run the `add` command with appropriate flags
3. Run `list` to show the updated task list

### When the user wants to see their tasks

1. Determine if they want all tasks or filtered (by status/priority)
2. Run the appropriate `list` command
3. If there are overdue tasks (due date in the past, still pending), proactively mention them

### When the user wants to update a task

1. If they reference a task by name/description rather than ID, run `list` first to find the ID
2. Run the `edit` command with only the changed fields
3. Run `show <id>` to confirm the changes

### When the user wants to complete a task

1. If referenced by name, find the ID via `list` first
2. Run `done <id>`
3. Run `list` to show remaining tasks

### When the user wants to delete a task

1. If referenced by name, find the ID via `list` first
2. **Confirm with the user before deleting** (unless they were explicit)
3. Run `delete <id>`
4. Run `list` to show remaining tasks

### When the user asks about overdue or due-soon tasks

1. Run `list` and examine the due dates in the output
2. Highlight any tasks that are overdue or due within 24 hours
3. Suggest actions (complete, reschedule, or escalate priority)

### When the user wants a summary or review

1. Run `list -a` to get all tasks
2. Provide a summary: count by status, count by priority, any overdue
3. Suggest next actions (e.g. "You have 3 high-priority tasks pending, 2 are overdue")

### When the user wants to set up notifications

Two channels are supported — configure one or both:
1. **ntfy** (recommended for phone push notifications): `config set ntfy_url "https://ntfy.sh/<topic>"`
2. **postback** (generic webhook): `config set postback_url "<url>"`

Steps:
1. Ask which channel(s) they want
2. For ntfy: run `config set ntfy_url "<url>"`, suggest installing the ntfy app and subscribing to the topic
3. For postback: run `config set postback_url "<url>"`
4. Start the daemon: suggest `todo daemon` or `docker compose up -d` (from the cloned repo directory)
5. Both channels fire independently — a notification succeeds if at least one channel delivers

### Bulk operations

When the user wants to operate on multiple tasks (e.g. "mark all high priority as in-progress"), iterate through the tasks:

1. Run `list` with appropriate filters to get IDs
2. Parse the output to extract IDs
3. Run the operation for each ID
4. Show the final `list` result

## Status Values

| Status | Meaning |
|--------|---------|
| `pending` | Not yet started (default for new tasks) |
| `in-progress` | Currently being worked on |
| `done` | Completed |

## Priority Values

| Priority | When to use |
|----------|-------------|
| `high` | Urgent, blockers, deadlines today/tomorrow, user says "urgent"/"critical"/"ASAP" |
| `medium` | Normal tasks, default |
| `low` | Nice-to-have, "whenever", "someday", low impact |

## Error Handling

- If a command fails, read the error message and explain it to the user in plain language
- If a task ID is not found, run `list -a` to help find the correct task
- If the database doesn't exist yet, any command will auto-create it — no special setup needed

## Important Notes

- The `todo` binary must be in `$PATH` — see Setup section above if not found
- Data is stored at `~/.local/share/todo/todo.db`
- All times displayed are in local timezone
- The `list` command hides done tasks by default — use `-a` to see them
- Task IDs are stable integers that auto-increment — they don't change when other tasks are deleted
