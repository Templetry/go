# Templetry parent: go

Go templates for [Templetry](https://github.com/Templetry). One **parent repo**, multiple **forms** — each form is a subdirectory that compiles on its own and carries its own `template.yml` ([ADR-0011](https://github.com/Templetry/wiki/blob/main/adr/0011-template-forms.md)).

| Form | What it is | Status |
|---|---|---|
| [`cli/`](cli/) | CLI starter — stdlib flags, versioned binary, table-driven tests | ✅ ready |
| [`http-service/`](http-service/) | HTTP service — stdlib mux (1.22 patterns), httptest suite, distroless Dockerfile | ✅ ready |

Pieces ([ADR-0014](https://github.com/Templetry/wiki/blob/main/adr/0014-lazy-pieces.md)): `http-service` ships `version-endpoint`, which adds `GET /version` **without touching a single existing file** — the form exposes an `api.Register` socket and the piece plugs in from its own `init`. That is the documented pattern for ecosystems whose wiring lives in code rather than in JSON.

```sh
templetry pieces ./my-svc
templetry add version-endpoint ./my-svc --set version_value=1.2.3
```

## Usage

```sh
templetry init go/http-service --out ./my-svc \
  --set "project_name=My Service" --set "module_path=github.com/me/my-svc"
```

Forms are **chosen**, not combined. Inside a form, the manifest's features are freely combinable.
