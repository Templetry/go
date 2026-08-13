# TemplateApp

Go CLI generated with [Templetry](https://github.com/Templetry): stdlib flags, versioned binary, table-driven tests.

```sh
go build -o template-app ./cmd/template-app
./template-app --name you
go test ./...
```

Release builds stamp the version: `go build -ldflags "-X main.version=v1.0.0" ./cmd/template-app`.
