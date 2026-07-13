# ---------- Build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Needed for some go modules that fetch over https during `go mod download`
RUN apk add --no-cache git ca-certificates

# Cache dependencies separately from source code
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static binary, no cgo, stripped symbols for a smaller image
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server .

# ---------- Runtime stage ----------
FROM alpine:3.20

# Needed for outbound HTTPS calls and for Postgres SSL if you enable sslmode later
RUN apk add --no-cache ca-certificates

# Non-root user — don't run the server as root inside the container
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/frontend ./frontend

RUN mkdir -p /app/uploads && chown -R app:app /app

USER app

ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["./server"]