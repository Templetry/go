# AGENTS

Operating contract for AI agents and automation helpers working in this project.

## Mission

- Keep this CLI dependency-free: stdlib `flag` until subcommands genuinely demand more.

## Core Rules

- Library code lives under `internal/`; `cmd/` stays a thin shell.
- Errors go to stderr with a non-zero exit; results go to stdout.
- Every exported function in `internal/` has a table-driven test.
- Update docs in the same change when behavior or process changes.

## Required Checks Before Finishing

- `go build ./...` and `go test ./...` pass.
- `gofmt` clean (no diff).

## Safe Change Workflow

1. Read the affected files fully before editing.
2. Make the smallest change that solves the task.
3. Build and test, then review the diff with git before committing.
