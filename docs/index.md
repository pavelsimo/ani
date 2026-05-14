---
title: ani
description: browse and search AniList anime from your terminal
---

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

Download the latest release from the [GitHub Releases](https://github.com/pavelsimo/ani/releases) page.
Pre-built binaries are available for macOS (amd64/arm64), Linux (amd64/arm64), and Windows (amd64).

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

# Show top 50 anime by score
ani top --limit 50

# JSON output for scripting
ani trending --json | jq '.[0].title.romaji'
```

## Commands

| Command | Description |
|---------|-------------|
| `ani` | Open interactive TUI (tab bar: Trending, Popular, Upcoming, All Time, Top 100, Search) |
| `ani trending` | Currently trending anime |
| `ani popular` | Popular anime this season |
| `ani upcoming` | Upcoming anime not yet airing |
| `ani all-time` | Most popular anime of all time |
| `ani top` | Highest-scored anime |
| `ani search <query>` | Search by title, genre, year, season, or format |
| `ani version` | Print version and exit |

## Search Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--genre` | Filter by genre (repeatable) | `--genre Action --genre Drama` |
| `--year` | Filter by year | `--year 2024` |
| `--season` | Filter by season | `--season winter` |
| `--format` | Filter by format | `--format tv` |
| `--status` | Filter by status | `--status airing` |
| `--min-score` | Minimum average score (0-100) | `--min-score 75` |

## TUI Keybindings

| Key | Action |
|-----|--------|
| `↑` / `↓` / `j` / `k` | Navigate list |
| `Tab` / `Shift+Tab` | Switch tabs |
| `/` | Open search input |
| `r` | Refresh current tab |
| `q` / `Ctrl+C` | Quit |

## Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--help` | `-h` | — | Show help |
| `--verbose` | `-v` | `false` | Verbose output |
| `--quiet` | `-q` | `false` | Suppress non-essential output |
| `--json` | — | `false` | Output as JSON |
| `--no-color` | — | `false` | Disable color output |
| `--lang` | — | `romaji` | Title language: `romaji` or `english` |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `NO_COLOR` | Set to any value to disable color output |

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Runtime failure |
| `2` | Invalid usage |

## Shell Completions

```bash
# bash
ani completion bash > /etc/bash_completion.d/ani

# zsh
ani completion zsh > "${fpath[1]}/_ani"

# fish
ani completion fish > ~/.config/fish/completions/ani.fish
```
