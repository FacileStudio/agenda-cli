# agenda-cli

Terminal client for the Facile Agenda calendar server.

## Install

```bash
# via facile installer
facile install agenda

# or via GoReleaser release
curl -fsSL https://github.com/FacileStudio/agenda-cli/releases/latest/download/install.sh | sh
```

## Usage

```bash
agenda config server https://agenda.example.com
agenda login
agenda whoami --json
agenda calendars list --json
agenda events create --calendar <id> --title "Meeting" --start 2026-09-02T14:00:00Z --end 2026-09-02T15:00:00Z
agenda events create --calendar <id> --title "Standup" --start "2026-09-07 09:00" --end "2026-09-07 09:15" \
  --timezone Europe/Paris --recurrence "FREQ=WEEKLY;BYDAY=MO" --status tentative \
  --conference-url https://meet.example/x --conference-provider Jitsi --location Room --description "Weekly"
agenda events update <id> --location "" --all-day=false --status confirmed
agenda token status
agenda caldav url
```

## Commands

- `login`, `logout`, `whoami [--json]`
- `auth login`, `auth logout`, `auth whoami`
- `token status [--json]`, `token create [-n <name>] [--json]`, `token revoke [--json]` (alias `delete`)
- `calendars list [--json]`, `create [name] [--name --color --description --echo-url --space]`, `show <id> [--json]`, `update <id> [--name --color --description --echo-url]`, `delete`, `share`, `members`, `unshare`
- `events list --calendar <id> [--from --to --search --timezone] [--json]`, `show <id> [--ics] [--json]`
- `events create` and `events update <id>` take `--title --start --end --all-day --location --description --recurrence --status --conference-url --conference-provider --timezone`, plus `--calendar`; `delete <id> -f`
- `caldav url`, `caldav sync-token <calendar>`, `caldav test`
- `config server <url>`, `config get`

`update` keeps every field you do not pass, and an empty value (`--location ""`) clears a text field.
`--start` and `--end` accept RFC 3339 or `YYYY-MM-DD[ HH:MM]`; the second form is read in `--timezone`,
UTC when omitted. The server stores no attendees, reminders, event colors or per-event URL besides
the conference link, so the CLI has no flags for them. `--search` filters the listing in the CLI.

Use `--json` on supported commands for machine-readable output without decorative styling.

## Config

`~/.config/agenda/config.yml`:

```yaml
server_url: https://agenda.example.com
token: <bearer-token>
user_id: <user-id>
email: <user@example.com>
```

Environment variables: `AGENDA_SERVER_URL`, `AGENDA_TOKEN`

## License

MIT
