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
agenda token status
agenda caldav url
```

## Commands

- `login`, `logout`, `whoami [--json]`
- `auth login`, `auth logout`, `auth whoami`
- `token status [--json]`, `token create [-n <name>] [--json]`, `token revoke [--json]` (alias `delete`)
- `calendars list [--json]`, `create`, `show <id> [--json]`, `update`, `delete`, `share`, `members`, `unshare`
- `events list --calendar <id> [...] [--json]`, `show <id> [--ics] [--json]`, `create`, `update`, `delete`
- `caldav url`, `caldav sync-token <calendar>`, `caldav test`
- `config server <url>`, `config get`

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
