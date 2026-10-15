FROM golang:1.27.2-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/todo_app /app/cmd/main.go

FROM alpine:3.24 AS runner
RUN adduser -D -u 10001 app && mkdir -p /app/out/logs && chown -R app:app /app
WORKDIR /app
COPY --from=builder --chown=app:app /app/bin/todo_app /app
USER app
ENV PROJECT_ROOT=/app
CMD ["/app/todo_app"]
