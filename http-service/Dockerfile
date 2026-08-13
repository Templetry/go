# Build stage
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /template-app ./cmd/server

# Runtime stage
FROM gcr.io/distroless/static-debian12
COPY --from=build /template-app /template-app
EXPOSE 8080
ENTRYPOINT ["/template-app"]
