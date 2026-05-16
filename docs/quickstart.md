---
title: Quick Start
description: Get productive with ani in 60 seconds
---

## Open the TUI

```bash
ani
```

Launches the interactive browser. Six tabs load lazily as you navigate them.
Press `q` or `Ctrl+C` to quit.

## Browse from the command line

```bash
# Trending anime right now
ani trending

# Popular this season
ani popular

# Upcoming — not yet airing
ani upcoming

# Most popular all time
ani all-time

# Highest scored (top 20)
ani top
```

## Search

```bash
# By title
ani search "fullmetal alchemist"

# By genre
ani search --genre Action --genre Drama

# Combined: title + year + minimum score
ani search "demon" --year 2021 --min-score 80
```

## Scripting with JSON

Every command accepts `--json` and writes a JSON array to stdout:

```bash
ani trending --json | jq '.[0].title.romaji'
ani search "spy" --json | jq '[.[] | {title: .title.romaji, score: .averageScore}]'
```

## Pagination

```bash
ani trending --page 2 --per-page 10
ani top --page 3 --limit 50
```

## Change title language

```bash
ani trending --lang romaji
ani trending --lang native
```

Valid values are `romaji`, `english` (default), and `native`.
