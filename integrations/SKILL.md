---
name: agenda
description: >
  Facile calendar CLI. Use when the user asks to create, list, manage,
  or inspect calendars and events, or test CalDAV connectivity.
---

# agenda: Facile calendar

Binary: `agenda`
Config: `~/.config/agenda/config.yml`
Token: `AGENDA_TOKEN` in environment or stored via `agenda login`

## When to apply

Use when the user mentions calendars, events, meetings, schedule, CalDAV, or Agenda.
Triggers: "calendar", "event", "meeting", "caldav", "agenda", "schedule"

## Commands

### Setup and auth
```
agenda config server <url>     Set the Agenda server URL
agenda login [--no-browser]    Authenticate via browser OIDC
agenda logout                  Clear stored credentials
agenda whoami [--json]         Show authenticated user profile
```

### Personal API tokens
```
agenda token status [--json]             Check personal API token status
agenda token create [-n <name>] [--json] Create or replace personal API token
agenda token revoke [--json]             Revoke personal API token (alias: delete)
```

### Calendars
```
agenda calendars list [--json]                                List all calendars
agenda calendars create <name> [--color <hex>]                Create a calendar
agenda calendars show <id> [--json]                           Show calendar details
agenda calendars update <id> [--name <name>] [--color <hex>]  Update calendar
agenda calendars delete <id>                                 Delete calendar
agenda calendars members <id>                                List calendar members
agenda calendars share <id> <email> [--role <role>]          Share calendar
agenda calendars unshare <id> <email>                        Remove user from calendar
```

### Events
```
agenda events list --calendar <id> [--from YYYY-MM-DD] [--to YYYY-MM-DD] [--search <q>] [--json]
agenda events show <id> [--ics] [--json]
agenda events create --calendar <id> --title <title> --start <RFC3339> --end <RFC3339>
agenda events update <calendar-id> <event-id> [--title <title>] [--start <RFC3339>] [--end <RFC3339>]
agenda events delete <calendar-id> <event-id>
```

### CalDAV
```
agenda caldav url                          Print CalDAV home set URL
agenda caldav sync-token <calendar-id>     Print calendar sync token
agenda caldav test                         Test CalDAV connectivity
```

## Rules

- Configure server URL first with `agenda config server <url>` or `AGENDA_SERVER_URL`
- In CI and automated contexts, set `AGENDA_TOKEN` and `AGENDA_SERVER_URL`
- Use `--json` for machine-readable JSON output without ANSI styling
- Run `agenda <cmd> -h` for exact syntax when unsure
