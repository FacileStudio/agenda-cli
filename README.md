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
agenda events create --calendar <id> --title "Review" --start "2026-09-08 15:00" --end "2026-09-08 16:00" \
  --timezone Europe/Paris --color "#ff8800" --attendee "bob@example.com,name=Bob,role=optional" \
  --attendee eve@example.com --reminder 15 --reminder 5:audio
agenda events update <id> --location "" --all-day=false --status confirmed
agenda events update <id> --clear-attendees --clear-reminders
agenda events cancel-occurrence <id> --date 2026-09-14
agenda events update <id> --occurrence 2026-09-21 --title "Moved" --start "2026-09-21 11:00"
agenda events restore-occurrence <id> --date 2026-09-14
agenda token status
agenda caldav url
```

## Commands

- `login`, `logout`, `whoami [--json]`
- `auth login`, `auth logout`, `auth whoami`
- `token status [--json]`, `token create [-n <name>] [--json]`, `token revoke [--json]` (alias `delete`)
- `calendars list [--json]`, `create [name] [--name --color --description --echo-url --space]`, `show <id> [--json]`, `update <id> [--name --color --description --echo-url]`, `delete`, `share`, `members`, `unshare`
- `events list --calendar <id> [--from --to --search --timezone] [--json]`, `show <id> [--ics] [--json]`
- `events create` and `events update <id>` take `--title --start --end --all-day --location --description --recurrence --status --conference-url --conference-provider --timezone --color --attendee --reminder`, plus `--calendar`; `update` also takes `--clear-attendees --clear-reminders`; `delete <id> -f`
- `events cancel-occurrence <id> --date`, `events restore-occurrence <id> --date`, `events update <id> --occurrence <date> [--title --description --location --start --end --status]`
- `caldav url`, `caldav sync-token <calendar>`, `caldav test`
- `config server <url>`, `config get`

`update` keeps every field you do not pass, and an empty value (`--location ""`) clears a text field.
`--start` and `--end` accept RFC 3339 or `YYYY-MM-DD[ HH:MM]`; the second form is read in `--timezone`,
UTC when omitted; on `update` the event's stored timezone is the default. `--attendee EMAIL[,name=,role=,response=]`
and `--reminder MINUTES[:method]` are repeatable and replace the stored list; `events show --json` reads
back each attendee's `response`. `--occurrence` and `--date` take `YYYY-MM-DD` in the event's timezone
or an RFC 3339 start. `--search` filters the listing in the CLI.

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
