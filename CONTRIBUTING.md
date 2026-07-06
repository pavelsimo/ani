# Contributing to ani

Thanks for your interest in contributing!

## Prerequisites

- Go 1.25 or later
- `make tools` — installs gofumpt, goimports, golangci-lint, lefthook, and goreleaser, and registers the git hooks

## Workflow

1. Fork and clone the repository.
2. Create a branch from `main`.
3. Make your changes. Keep the conventions in [AGENTS.md](AGENTS.md) in mind — it documents the project layout, coding style, and TUI architecture.
4. Run the full gate before pushing:

   ```bash
   make ci   # fmt-check + lint + test (70% coverage gate) + build
   ```

   The lefthook pre-commit hook runs `fmt-check` and `lint` automatically.

## Tests

- Table-driven tests live in `_test.go` files adjacent to the code under test.
- Mock the HTTP client — see `internal/anilist/client_test.go` and `cmd/cmd_test.go` for the pattern. Tests must not hit the real AniList API.
- CI enforces ≥ 70% total coverage (`make coverage-check`).

## Commit messages

Commits follow `<emoji> <lowercase imperative summary>`:

- ✨ new feature · 🐛 bug fix · ♻️ refactor · 📝 docs
- 👷 CI/CD · 🔧 config · ⬆️ dependency bump · 🔥 remove code

Example: `🐛 clamp per-page to the anilist cap`

## Pull requests

- Keep PRs focused — one logical change per PR.
- Describe what changed and why; link related issues.
- Make sure `make ci` passes; CI runs the same gate.
