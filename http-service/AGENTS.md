# AGENTS

Operating contract for AI agents and automation helpers working in this project.

## Mission

- Keep this service on the stdlib: `net/http` 1.22 patterns until real needs (middleware chains, websockets) earn a dependency.

## Core Rules

- Routes are registered in `internal/api.NewMux` with `METHOD /path` patterns; handlers stay small and return JSON via `writeJSON`.
- `cmd/server` is a thin shell: config from env, no logic.
- Every route gets an httptest-based test.
- Update docs in the same change when behavior or process changes.

## Required Checks Before Finishing

- `go build ./...` and `go test ./...` pass.
- `gofmt` clean (no diff).

## Safe Change Workflow

1. Read the affected files fully before editing.
2. Make the smallest change that solves the task.
3. Build and test, then review the diff with git before committing.
