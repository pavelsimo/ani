---
title: TUI Guide
description: Interactive terminal UI — tabs, navigation, and search
---

Run `ani` with no arguments to open the interactive TUI.

## Tabs

| Tab | Content |
|-----|---------|
| Trending | Currently trending anime |
| Popular | Popular this season |
| Upcoming | Not yet airing |
| All Time | Most popular of all time |
| Top 100 | Highest scored |
| Search | Search with filters |

Tabs are lazy-loaded — data fetches only when you first visit each tab.
Switching tabs is instant after the first load.

## Keybindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `→` / `Tab` | Next tab |
| `←` / `Shift+Tab` | Previous tab |
| `/` | Open search input |
| `Enter` | Confirm search |
| `Esc` | Cancel search |
| `r` | Refresh current tab |
| `q` / `Ctrl+C` | Quit |

## Search

Press `/` from any tab to open the inline search input.
Type a query and press `Enter` — results populate the Search tab.
Press `Esc` to cancel without searching.

## Title language

Launch with `--lang english` to display English titles instead of romaji:

```bash
ani --lang english
```
