# ani

Browse and search AniList anime from your terminal.

[![release](https://img.shields.io/github/v/release/pavelsimo/ani?style=flat-square&color=4d9e4d&logoColor=white)](https://github.com/pavelsimo/ani/releases)
[![license MIT](https://img.shields.io/badge/license-MIT-ffd60a?style=flat-square&logoColor=white)](LICENSE)
[![Go 1.25](https://img.shields.io/badge/Go-1.25-2ea44f?style=flat-square&logoColor=white)](https://go.dev)
[![Homebrew](https://img.shields.io/badge/Homebrew-b28f62?style=flat-square&logoColor=white)](https://github.com/pavelsimo/homebrew-tap)
[![DeepWiki](https://img.shields.io/badge/DeepWiki-0088cc?style=flat-square&logoColor=white)](https://deepwiki.com/pavelsimo/ani)

## Installation

### Homebrew (macOS / Linux)

```bash
brew tap pavelsimo/homebrew-tap
brew install ani
```

### Go install

```bash
go install github.com/pavelsimo/ani@latest
```

### Download binary

Pre-built binaries for macOS (amd64/arm64), Linux (amd64/arm64), and Windows (amd64) are available on the [Releases](https://github.com/pavelsimo/ani/releases) page.

## Quick Start

```bash
# Open the interactive TUI
ani

# Show trending anime as a table
ani trending

# Search for an anime by title
ani search "attack on titan"

# Filter by genre and year
ani search --genre Action --genre Drama --year 2024

# Show popular anime this season
ani popular

# JSON output for scripting
ani trending --json | jq '.[0].title.romaji'
```

## Commands

| Command | Description |
|---------|-------------|
| `ani` | Open interactive TUI (tabs: Trending, Popular, Upcoming, All Time, Top 100, Search) |
| `ani trending` | Currently trending anime |
| `ani popular` | Popular anime this season |
| `ani upcoming` | Upcoming anime not yet airing |
| `ani all-time` | Most popular anime of all time |
| `ani top` | Highest-scored anime |
| `ani search <query>` | Search by title, genre, year, season, or format |
| `ani version` | Print version and exit |

## TUI Keybindings

| Key | Action |
|-----|--------|
| `↑` / `↓` / `j` / `k` | Navigate list |
| `Tab` / `Shift+Tab` | Switch tabs |
| `/` | Open search input |
| `r` | Refresh current tab |
| `q` / `Ctrl+C` | Quit |

## Docs

Full documentation at **[pavelsimo.github.io/ani](https://pavelsimo.github.io/ani)**.

## License

MIT
