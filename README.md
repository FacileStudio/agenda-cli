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
agenda login --password <secret>
agenda calendars list
agenda events create --calendar <id> --title "Meeting" --start 2026-09-02T14:00:00Z --end 2026-09-02T15:00:00Z
agenda caldav url
```

## Commands

- `auth login`, `auth logout`, `auth whoami`
- `calendars list`, `create`, `show`, `update`, `delete`, `share`, `members`, `unshare`
- `events list`, `show`, `create`, `update`, `delete`
- `caldav url`, `caldav sync-token <calendar>`, `caldav test`

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
