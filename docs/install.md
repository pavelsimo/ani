---
title: Install
description: Install ani on macOS, Linux, or Windows
---

## Homebrew (macOS / Linux)

```bash
brew tap pavelsimo/homebrew-tap
brew install ani
```

This is the recommended method. Tap updates are automatic with `brew upgrade`.

## Go install

```bash
go install github.com/pavelsimo/ani@latest
```

Requires Go 1.21 or later. The binary lands in `$GOPATH/bin`.

## Prebuilt binary

Download the latest release from the [GitHub Releases](https://github.com/pavelsimo/ani/releases) page.
Pre-built binaries are available for:

- macOS amd64 / arm64
- Linux amd64 / arm64
- Windows amd64

Extract the archive and place `ani` (or `ani.exe`) somewhere on your `PATH`.

## Verify

```bash
ani version
```

You should see the current version string printed to stdout.
