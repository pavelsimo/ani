---
title: Reference
description: Global flags, environment variables, exit codes, and shell completions
---

## Global flags

These flags apply to every command.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--help` | `-h` | — | Show help for the current command |
| `--verbose` | `-v` | `false` | Enable verbose output |
| `--quiet` | `-q` | `false` | Suppress non-essential output |
| `--json` | — | `false` | Output results as a JSON array |
| `--no-color` | — | `false` | Disable ANSI color output |
| `--lang` | — | `english` | Title language: `romaji`, `english`, or `native` |
| `--type` | — | `anime` | Media type: `anime` or `manga` |

## Environment variables

| Variable | Description |
|----------|-------------|
| `NO_COLOR` | Set to any non-empty value to disable color output (equivalent to `--no-color`) |

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Runtime failure (API error, network error) |
| `2` | Invalid usage (bad flags, missing required argument) |

## Shell completions

ani ships Cobra-generated completions for bash, zsh, and fish.

```bash
# bash — system-wide
ani completion bash > /etc/bash_completion.d/ani

# zsh — user fpath
ani completion zsh > "${fpath[1]}/_ani"

# fish
ani completion fish > ~/.config/fish/completions/ani.fish
```

After installing, reload your shell or source the completion file to activate it.

## Data source

ani queries the [AniList GraphQL API](https://anilist.co) (`https://graphql.anilist.co`).
No API key or account is required. The public rate limit is 90 requests per minute.
