# TemplateApp

REST API over SQLite generated with [Templetry](https://github.com/Templetry): pure-Go driver (no CGO, no build toolchain), idempotent migrations at startup, repository layer, httptest suite against an in-memory database.

```sh
go run ./cmd/server        # :8080, db template-app.db (PORT and DATABASE_DSN override)
go test ./...
docker build -t template-app .   # docker feature
```

## Endpoints

`GET /healthz` · `POST /api/notes` · `GET /api/notes` · `GET|PUT|DELETE /api/notes/{id}`

## Adding resources

This form ships a **piece per object**: adopt it once per entity and get the table, repository, CRUD routes and tests — wired through the `store.Register` and `api.Register` sockets, so no existing file is touched.

```sh
templetry add crud-resource --set entity=Product
```
