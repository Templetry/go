# AGENTS

Operating contract for AI agents and automation helpers working in this project.

## Mission

- A REST API over SQLite with no CGO: `database/sql` + `modernc.org/sqlite`, hand-written SQL, thin handlers.

## Core Rules

- SQL lives in `internal/store` repositories; handlers never build queries.
- Schema statements stay **idempotent** (`CREATE TABLE IF NOT EXISTS`) — they run on every start.
- New resources register their routes via `api.Register` and their DDL via `store.Register`, both from `init` — never edit `NewMux` or the base schema to add one.
- Missing rows return `store.ErrNotFound`, which `Fail` maps to 404; never leak driver errors to clients.
- Every route gets an httptest test against an in-memory database.
- Update docs in the same change when behavior or process changes.

## Required Checks Before Finishing

- `go build ./...` and `go test ./...` pass.
- `gofmt` clean (no diff).

## Safe Change Workflow

1. Read the affected files fully before editing.
2. Make the smallest change that solves the task.
3. Build and test, then review the diff with git before committing.
