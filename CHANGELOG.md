# Changelog

All notable changes to agenda-cli are documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/), and the project adheres to
[Semantic Versioning](https://semver.org/).

## [Unreleased]

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