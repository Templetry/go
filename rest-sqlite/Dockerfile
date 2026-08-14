# Build stage
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# CGO stays off: the SQLite driver is pure Go.
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /template-app ./cmd/server

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=build /template-app /template-app
ENV DATABASE_DSN=/data/template-app.db
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["/template-app"]
