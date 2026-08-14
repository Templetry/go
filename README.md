# Templetry parent: go

Go templates for [Templetry](https://github.com/Templetry). One **parent repo**, multiple **forms** — each form is a subdirectory that compiles on its own and carries its own `template.yml` ([ADR-0011](https://github.com/Templetry/wiki/blob/main/adr/0011-template-forms.md)).

| Form | What it is | Status |
|---|---|---|
| [`cli/`](cli/) | CLI starter — stdlib flags, versioned binary, table-driven tests | ✅ ready |
| [`http-service/`](http-service/) | HTTP service — stdlib mux (1.22 patterns), httptest suite, distroless Dockerfile | ✅ ready |

| [`rest-sqlite/`](rest-sqlite/) | REST API over SQLite — pure-Go driver (no CGO), migrations, repositories, httptest suite | ✅ ready |

Pieces ([ADR-0014](https://github.com/Templetry/wiki/blob/main/adr/0014-lazy-pieces.md)) live in `_pieces/` (the underscore keeps `go build ./...` clean in the template itself) and wire themselves through **sockets** — `api.Register` for routes, `store.Register` for schema — so adopting one touches no existing file:

| Form | Piece | What it adds |
|---|---|---|
| `http-service` | `version-endpoint` | `GET /version` with a configurable version string |
| `rest-sqlite` | `crud-resource` | **A whole entity**: table, repository, CRUD routes and tests, renamed to your object |

```sh
templetry pieces ./my-svc
templetry add crud-resource ./my-svc --set entity=Product
```

## Usage

```sh
templetry init go/http-service --out ./my-svc \
  --set "project_name=My Service" --set "module_path=github.com/me/my-svc"
```

Forms are **chosen**, not combined. Inside a form, the manifest's features are freely combinable.
