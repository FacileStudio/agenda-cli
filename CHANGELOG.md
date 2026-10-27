# Changelog

All notable changes to agenda-cli are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/), and the project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

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