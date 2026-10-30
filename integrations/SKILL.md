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
agenda calendars list [--json]                                       List all calendars
agenda calendars create <name> [--color <hex>] [--json]              Create a calendar
agenda calendars show <id> [--json]                                  Show calendar details
agenda calendars update <id> [--name] [--color] [--description] [--echo-url] [--json]  Update calendar (unpassed fields kept)
agenda calendars delete <id> [-f] [--json]                           Delete calendar
agenda calendars members <id> [--json]                               List calendar members
agenda calendars share <id> <email> <role> [--json]                  Share calendar
agenda calendars unshare <id> <email-or-id> [--json]                 Remove member from calendar
```

### Events
```
agenda events list --calendar <id> [--from YYYY-MM-DD] [--to YYYY-MM-DD] [--search <q>] [--timezone <iana>] [--json]
agenda events show <id> [--ics] [--json]
agenda events create --calendar <id> --title <t> --start <dt> --end <dt> [--all-day] [--location <l>] [--description <d>]
    [--recurrence "FREQ=WEEKLY;BYDAY=MO"] [--status confirmed|tentative|cancelled]
    [--conference-url <url>] [--conference-provider <name>] [--timezone Europe/Paris] [--json]
agenda events update <id> [any create flag] [--calendar <id>] [--all-day=false] [--location ""] [--json]
agenda events delete <id> [-f] [--json]
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
- `events update` and `calendars update` keep unpassed fields; an empty value clears a text field
- Naive times (`2026-09-07 09:00`) are read in `--timezone`, UTC when omitted; pass `--timezone` with `--recurrence`
- The server has no attendees, reminders or event color: do not look for flags for them
- Run `agenda <cmd> -h` for exact syntax when unsure
