# Changelog

All notable changes to agenda-cli are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/), and the project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

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