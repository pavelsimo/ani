---
title: Commands
description: All ani subcommands and their flags
---

## ani (no arguments)

Opens the interactive TUI.

```bash
ani
ani --lang english
```

See the [TUI Guide](tui.html) for keybindings and tab details.

## ani trending

Currently trending anime on AniList.

```bash
ani trending
ani trending --page 2 --per-page 10
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `1` | Page number |
| `--per-page` | `20` | Results per page (max 50) |

## ani popular

Popular anime for the current or a specified season.

```bash
ani popular
ani popular --season winter --year 2024
```

| Flag | Default | Description |
|------|---------|-------------|
| `--season` | current | Season: `winter`, `spring`, `summer`, `fall` |
| `--year` | current | Year |
| `--page` | `1` | Page number |
| `--per-page` | `20` | Results per page (max 50) |

## ani upcoming

Anime that have not yet started airing.

```bash
ani upcoming
ani upcoming --page 2
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `1` | Page number |
| `--per-page` | `20` | Results per page (max 50) |

## ani all-time

Most popular anime of all time by AniList popularity score.

```bash
ani all-time
ani all-time --per-page 50
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `1` | Page number |
| `--per-page` | `20` | Results per page (max 50) |

## ani top

Highest-scored anime by average score.

```bash
ani top
ani top --limit 50 --page 2
```

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `20` | Number of results (max 50 per page) |
| `--page` | `1` | Page number |

## ani version

Print the current version and exit.

```bash
ani version
```
