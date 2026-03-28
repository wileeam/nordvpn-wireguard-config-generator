# Multi-stage Dockerfile for NordGen Web Application
# Stage 1: Build the Vue.js frontend
FROM node:23-alpine AS frontend-builder

WORKDIR /frontend

# Copy frontend package files
COPY web-Frontend/package.json ./

# Install dependencies
RUN npm install

# Copy frontend source code
COPY web-Frontend/ ./

# Build the frontend
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.25.6-alpine3.23 AS backend-builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid 10001 appuser

# Copy Go modules files
COPY web-Backend/go.mod web-Backend/go.sum ./

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy backend source code
COPY web-Backend/ ./

# Copy the built frontend from the previous stage
COPY --from=frontend-builder /frontend/dist ./public

# Build the Go binary
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOAMD64=v3 go build -buildvcs=false -ldflags="-s -w" -trimpath -o server main.go

# Stage 3: Create the final minimal image
FROM scratch

WORKDIR /app

# Copy necessary files from builder
COPY --from=backend-builder /etc/passwd /etc/passwd
COPY --from=backend-builder /etc/group /etc/group
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary and static assets
COPY --from=backend-builder /build/server .
COPY --from=backend-builder /build/public ./public

# Use non-root user
USER appuser:appuser

# Expose the application port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/server", "--help"] || exit 1

# Run the application
ENTRYPOINT ["./server"]
