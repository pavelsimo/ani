# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- The TUI now respects `--no-color` and the `NO_COLOR` environment variable
- `--per-page` and `top --limit` are clamped to AniList's 50-result cap instead of failing at the API
- `--page 0` and negative pages are rejected with a clear error
- `ani search` with no query or filters now errors instead of silently querying everything
- API errors are printed once instead of twice
- Rate-limit (429) responses now produce a clear message including the Retry-After value
- `ani info` on an unknown ID reports "no media found" instead of rendering a blank page
- Episodes that already aired show "aired" instead of "airing soon"
- `make tools` works on Go 1.25 (updated gofumpt/goimports pins, goreleaser v2 module path)

### Removed
- Unused `--verbose` and `--quiet` flags

### Changed
- Reduced allocations in the TUI render loop (cached width-derived and genre tag styles)

## [0.4.0] - 2026-05-24

### Added
- install on Windows via `winget install pavelsimo.ani`

## [0.3.0] - 2026-05-16

### Added
- Show community recommendations in the anime/manga detail view, listed after the synopsis

## [0.2.0] - 2026-05-16

### Added
- Full anime detail view with synopsis, genres, tags, studios, relations, and streaming links via `ani info <id>`
- In-app detail panel accessible from any list by pressing Enter
- Browse and search manga titles alongside anime from any list view
- Navigate to additional results pages in the TUI with next-page pagination
- Clickable streaming links in the detail panel
- Native language title support in the detail view

### Changed
- Improve search tab interface for faster, more intuitive filtering

### Fixed
- Fix column misalignment and duplicate episode count label in TUI list views

## [0.1.0] - 2026-05-14

### Added
- Browse trending, popular, upcoming, all-time, and top-100 anime directly from the terminal
- Interactive TUI with tab navigation, live search, and keyboard shortcuts
- JSON output mode for scripting via `--json`
- Title language selection (romaji / english) via `--lang`
- Install via Homebrew tap

[Unreleased]: https://github.com/pavelsimo/ani/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/pavelsimo/ani/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/pavelsimo/ani/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/pavelsimo/ani/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/pavelsimo/ani/releases/tag/v0.1.0
