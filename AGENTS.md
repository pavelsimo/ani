# Agent Instructions — ani

This file is the single canonical source of agent instructions for this repository.
`CLAUDE.md` is a symlink to this file. Do not create separate `CLAUDE.md` or `CODEX.md` variants.

## Project Structure

```
ani/
├── cmd/                  # Cobra commands (one file per subcommand)
│   ├── root.go           # Root command, global flags (--lang, --json, --no-color)
│   ├── version.go        # version subcommand
│   ├── tui.go            # default (no args) → launches Bubble Tea TUI
│   ├── trending.go       # ani trending
│   ├── popular.go        # ani popular [--season --year]
│   ├── upcoming.go       # ani upcoming
│   ├── alltime.go        # ani all-time
│   ├── top.go            # ani top [--limit]
│   └── search.go         # ani search [query] [filters]
├── internal/
│   ├── anilist/          # AniList GraphQL client + types + query constants
│   │   ├── client.go     # HTTP client, Query(ctx, query, vars) → *Page
│   │   ├── types.go      # Media, Title, Page, AiringEpisode, FuzzyDate, SearchParams
│   │   └── queries.go    # GraphQL query string constants
│   ├── display/          # Rendering layer (lipgloss tables, genre tags, formatters)
│   │   ├── table.go      # Render(media, lang, noColor) → string
│   │   ├── tags.go       # RenderTags(genres) → colored pill string
│   │   └── format.go     # Title, Score, Popularity, Format, Status, Episodes, NextEp
│   └── tui/              # Bubble Tea TUI
│       ├── model.go      # Root Model, Init/Update/View, Start(client, lang)
│       ├── keys.go       # Key bindings (keys.Quit, keys.Search, etc.)
│       └── styles.go     # All lipgloss styles in one place
├── docs/                 # Markdown documentation source
├── scripts/              # Build tooling (docs site builder)
├── .github/workflows/    # CI, release, and pages workflows
├── Makefile              # All developer tasks
├── .golangci.yml         # Linter configuration
├── .goreleaser.yaml      # Release configuration
└── .lefthook.yml         # Git hook configuration
```

## Build / Test / Dev Commands

```bash
make build        # compile binary to bin/ani
make test         # run tests with coverage
make lint         # run golangci-lint
make fmt          # format code (gofumpt + goimports)
make fmt-check    # CI-safe format check (exits non-zero if dirty)
make ci           # full gate: fmt-check + lint + test + build
make docs         # build documentation site to dist/docs-site/
make tools        # install all development tools
```

## AniList API

- Endpoint: `POST https://graphql.anilist.co`
- No authentication needed for read queries
- Rate limit: 90 requests/minute
- All query strings live in `internal/anilist/queries.go`
- `internal/anilist/client.go` — `New()` returns `*Client`; `Query(ctx, query, vars)` returns `*Page`

## Coding Style

- All errors returned from `RunE` are printed to stderr; do not use `fmt.Println` for errors
- Primary data to stdout; diagnostics, progress, and errors to stderr
- `SilenceUsage: true` on every `cobra.Command` with a `RunE`
- Flag names are lowercase hyphenated (never camelCase)
- Respect `NO_COLOR` env and `--no-color` flag everywhere color is used
- `--lang romaji|english` controls title display (default: romaji)
- `--json` produces JSON array to stdout; human table is a rendering layer on the same structs

## TUI Architecture

- Entry point: `tui.Start(client, lang)` in `internal/tui/model.go`
- Tabs: Trending (0), Popular (1), Upcoming (2), All Time (3), Top 100 (4), Search (5)
- Data is lazy-loaded per tab on first activation
- Search modal activated with `/`; executes on Enter, cancels on Esc
- All lipgloss styles are defined in `internal/tui/styles.go` — edit there, nowhere else
- Key bindings are defined in `internal/tui/keys.go`

## Testing Guidelines

- Table-driven tests in `_test.go` files adjacent to the code under test
- Mock the HTTP client for `internal/anilist` tests
- Aim for ≥ 70% coverage on `internal/`

## Commit & PR Guidelines

Commits follow the [/commit skill](https://github.com/pavelsimo/commit) convention:
`<emoji> <lowercase imperative summary>`

Common emoji:
- ✨ new feature · 🐛 bug fix · ♻️ refactor · 📝 docs
- 👷 CI/CD · 🔧 config · ⬆️ dependency bump · 🔥 remove code
