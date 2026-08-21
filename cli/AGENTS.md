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

```sh templetry:checks
go build ./...
go test ./...
```

## Safe Change Workflow

1. Read the affected files fully before editing.
2. Make the smallest change that solves the task.
3. Build and test, then review the diff with git before committing.

## This project came from a template

Four facts you cannot infer from the code in front of you:

- **Never hand-edit `.templetry-answers.yml`.** It records what generated this project. Editing it makes the next update merge against a state that never existed.
- **Before writing a capability by hand, run `templetry pieces`.** Auth, RBAC, audit trails, API keys and whole CRUD resources may already exist as pieces for this template. Adopting one is `templetry add <name>`, and it brings its own tests.
- **`templetry update` pulls improvements from the template** through a three-way merge that keeps your edits. Use it instead of copying files from the template by hand.
- **Directives like `tpl:if` belong to the template, not here.** If you find one in this project, it is a rendering bug worth reporting — do not try to interpret it.
