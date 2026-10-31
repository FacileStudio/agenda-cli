# Changelog

All notable changes to agenda-cli are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/), and the project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.5.0] - 2026-10-31

Needs an Agenda server with attendee, reminder, color and occurrence support.

### Added

- `events create` and `events update` take `--attendee` (repeatable, `EMAIL[,name=,role=,response=]`), `--clear-attendees`, `--reminder` (repeatable, `MINUTES[:display|audio|email]`), `--clear-reminders` and `--color`.
- `events show` prints the timezone, color, attendees with their response, reminders and occurrence exceptions; `--ics` carries ATTENDEE, VALARM, COLOR and EXDATE.
- `events cancel-occurrence <id> --date`, `events restore-occurrence <id> --date` and `events update <id> --occurrence <date>` act on a single occurrence of a recurring event.
- `events list` prints the `occurrence_start` of each expanded occurrence.

### Fixed

- `events update` keeps the stored timezone, so a recurring series is no longer re-anchored to UTC, and reads naive `--start` and `--end` in that zone.

## [0.4.0] - 2026-10-30

### Added

- `events create` and `events update` take `--recurrence`, `--status`, `--conference-url`, `--conference-provider` and `--timezone`. Naive `--start` and `--end` are read in `--timezone`.
- `calendars create` and `calendars update` take `--description` and `--echo-url`; `create` takes `--space` and `--name`.
- `events list --timezone` expands recurring events in that zone.
- `events show` prints recurrence, status and conference fields; `--ics` carries RRULE, STATUS and URL.

### Fixed

- `events update` no longer drops recurrence, status or conference fields, and `--all-day=false` and an empty `--location ""` now clear the stored value.
- `calendars update` fetches the calendar first, so passing only `--color` no longer fails or wipes the description and echo URL.
- `events list --search` filters in the CLI; the server ignored it.
- `events show`, `events list` and `--ics` read `start_at`, `end_at` and `is_all_day`, so they no longer print empty times.
- `events delete`, `calendars delete` and `calendars unshare` report a server error instead of success.

## [0.3.2] - 2026-10-27

### Added

- Added `--json` flag to `calendars create`, `update`, `delete`, `members`, `share`, `unshare`, and `events create`, `update`, `delete`.
- Flexible datetime parsing for `events create` and `update` supporting ISO 8601, RFC 3339, and `YYYY-MM-DD HH:MM`.
- Member resolution in `calendars unshare` accepting both email addresses and integer member IDs.

### Fixed

- Fixed event creation and update payload fields to match Agenda API contract (`start_at`, `end_at`, `is_all_day` instead of `start`, `end`, `all_day`).
- Preserved existing event fields during partial updates (`events update`) to avoid overwriting unchanged values.

## [0.3.1] - 2026-10-27

### Fixed

- Default server URL to `https://agenda.facile.studio/api` so running with `AGENDA_TOKEN` requires no extra URL configuration.
- Automatically resolve caller profile from `/users/me` on `agenda whoami` when authenticated via environment variable or token without local config.
- Normalize server URLs pointing to `agenda.facile.studio` to append `/api`.

## [0.3.0] - 2026-10-27

### Added

- `--json` flags across commands (`whoami`, `calendars list`, `calendars show`, `events list`, `events show`, and `token` commands) for raw JSON output.
- `agenda token` command group (`status`, `create`, `revoke`/`delete`) for personal API token management matching Agenda server endpoints.
- Comprehensive unit tests covering token commands, json formatting, whoami, calendars, and events.

### Fixed

- OIDC flow endpoint correction.
- User profile path `/users/me`.
- Config file missing fallback and `AGENDA_SERVER_URL` / `AGENDA_TOKEN` resolution.
- `auth`, `whoami`, and `logout` command registration on root command.

## [v0.2.1] - 2026-09-10

### Fixed

- `--version` now prints the suite-standard `<bin> <semver>` form (`agenda 0.2.1`)
  instead of cobra's default `agenda version 0.2.1`.

## [v0.2.0] - 2026-09-10

### Changed

- Tightened the filet gate to fail on info-level findings and refactored every
  package to pass it: extracted RunE logic into `run*` helpers, split oversized
  files, added doc comments to all exported symbols, replaced `interface{}`
  with `any`, and removed dead code.
- Replaced the loopback browser-open stubs with real `os/exec` calls.
- Moved the UI init out of a package `init()` into the root command constructor.
- Updated goreleaser to version 2 schema and injected the version via ldflags.

### Removed

- Dropped a stray test binary from the repository.