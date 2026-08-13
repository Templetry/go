# TemplateApp

Go HTTP service generated with [Templetry](https://github.com/Templetry): stdlib mux (Go 1.22 method+path patterns), health endpoint, httptest suite, optional distroless Dockerfile.

```sh
go run ./cmd/server          # listens on :8080 (PORT overrides)
go test ./...
docker build -t template-app .   # docker feature
```

Routes: `GET /healthz` · `GET /api/hello/{name}`.
