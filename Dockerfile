# CrossForge Dockerfile
# Multi-stage build for minimal production image

# =============================================================================
# Build Stage
# =============================================================================
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO is needed for SQLite with go-sqlite3, use modernc.org/sqlite for pure Go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o crossforge \
    ./cmd/crossforge

# =============================================================================
# Production Stage
# =============================================================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite

# Create non-root user
RUN addgroup -g 1000 crossforge && \
    adduser -u 1000 -G crossforge -h /app -D crossforge

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/crossforge /app/crossforge

# Copy static assets and templates
COPY --from=builder /build/web /app/web
COPY --from=builder /build/internal/frontend/templates /app/templates
COPY --from=builder /build/internal/storage/migrations /app/migrations

# Create data directory
RUN mkdir -p /app/data && chown -R crossforge:crossforge /app

# Switch to non-root user
USER crossforge

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set environment defaults
ENV SERVER_HOST=0.0.0.0 \
    SERVER_PORT=8080 \
    DATABASE_PATH=/app/data/crossforge.db \
    LLM_ENDPOINT=http://host.docker.internal:11434

# Run the application
ENTRYPOINT ["/app/crossforge"]
CMD ["serve"]
