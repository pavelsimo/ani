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

Currently trending anime or manga on AniList (use `--type manga` for manga).

```bash
ani trending
ani trending --page 2 --per-page 10
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `1` | Page number |
| `--per-page` | `20` | Results per page (max 50) |

## ani popular

Popular anime (or manga with `--type manga`) for the current or a specified season.

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

Titles that have not yet started airing (use `--type manga` for unreleased manga).

```bash
ani upcoming
ani upcoming --page 2
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `1` | Page number |
| `--per-page` | `20` | Results per page (max 50) |

## ani all-time

Most popular anime (or manga with `--type manga`) of all time by AniList popularity score.

```bash
ani all-time
ani all-time --per-page 50
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `1` | Page number |
| `--per-page` | `20` | Results per page (max 50) |

## ani top

Highest-scored anime or manga by average score (use `--type manga` for manga).

```bash
ani top
ani top --limit 50 --page 2
```

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `20` | Number of results (max 50 per page) |
| `--page` | `1` | Page number |

## ani info

Show full details for a single title by its AniList ID.

```bash
ani info 1
ani info 1 --json
ani info 1 --lang romaji
```

The AniList ID is a positive integer visible in the AniList URL (e.g. `anilist.co/anime/1`).

The detail view includes synopsis, score, studios, genres, tags (top non-spoiler), relations (prequel/sequel/etc.), and streaming links. Supports `--json`, `--lang`, and `--no-color`.

## ani version

Print the current version and exit.

```bash
ani version
```
