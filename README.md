# Templetry parent: go

Go templates for [Templetry](https://github.com/Templetry). One **parent repo**, multiple **forms** — each form is a subdirectory that compiles on its own and carries its own `template.yml` ([ADR-0011](https://github.com/Templetry/wiki/blob/main/adr/0011-template-forms.md)).

| Form | What it is | Status |
|---|---|---|
| [`cli/`](cli/) | CLI starter — stdlib flags, versioned binary, table-driven tests | ✅ ready |
| [`http-service/`](http-service/) | HTTP service — stdlib mux (1.22 patterns), httptest suite, distroless Dockerfile | ✅ ready |

## Usage

```sh
templetry init go/http-service --out ./my-svc \
  --set "project_name=My Service" --set "module_path=github.com/me/my-svc"
```

Forms are **chosen**, not combined. Inside a form, the manifest's features are freely combinable.
