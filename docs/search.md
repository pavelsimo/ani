---
title: Search
description: Search and filter anime by title, genre, year, season, format, and score
---

## Basic search

```bash
ani search "attack on titan"
ani search "made in abyss"
```

A query string matches against anime titles (romaji and English).
You must provide at least a query or one filter flag.

## Filter flags

| Flag | Type | Description | Example |
|------|------|-------------|---------|
| `--genre` | string (repeatable) | Filter by genre | `--genre Action` |
| `--year` | int | Filter by season year | `--year 2024` |
| `--season` | string | Filter by season | `--season winter` |
| `--format` | string | Filter by format | `--format tv` |
| `--status` | string | Filter by airing status | `--status airing` |
| `--min-score` | int (0–100) | Minimum average score | `--min-score 75` |
| `--page` | int | Page number | `--page 2` |
| `--per-page` | int | Results per page (max 50) | `--per-page 10` |

## Valid values

**`--season`:** `winter`, `spring`, `summer`, `fall`

**`--format`:** `tv`, `ova`, `ona`, `movie`, `special`

**`--status`:** `airing` (alias: `releasing`), `finished`, `upcoming` (alias: `not_yet_released`), `cancelled`

## Examples

```bash
# All Action anime from 2023
ani search --genre Action --year 2023

# Multiple genres
ani search --genre Action --genre Drama

# Title search with score filter
ani search "jujutsu" --min-score 80

# Finished TV series from winter 2022
ani search --format tv --season winter --year 2022 --status finished

# Pipe into jq
ani search "chainsaw" --json | jq '.[] | .title.romaji + " " + (.averageScore|tostring)'
```

## Pagination

```bash
ani search --genre "Slice of Life" --page 2 --per-page 15
```
