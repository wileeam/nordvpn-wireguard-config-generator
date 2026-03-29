# NordGen Backend

A minimalist, high-performance backend service for generating NordVPN WireGuard configurations. Built on the **Go** programming language and the **Fiber** framework, this service handles server data caching, credential exchange, and configuration generation with extreme efficiency.

## Overview

This application serves as the API layer for the NordGen project. It interfaces directly with NordVPN's infrastructure to retrieve server lists and exchange authentication tokens for WireGuard private keys. It provides endpoints to generate configuration files in text, file, or QR code formats.

## Prerequisites

- Go 1.25+ (Recommended)

## Installation

Clone the repository and download the dependencies:

```bash
git clone https://github.com/mustafachyi/NordVPN-WireGuard-Config-Generator
cd NordVPN-WireGuard-Config-Generator
go mod download
```

## Development

Start the server directly using the Go toolchain:

```bash
go run main.go
```

The server listens on port `3000` by default.

## Production

To build and run the optimized production binary:

```bash
# Build the binary with size optimizations (strip debug symbols)
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o server main.go

# Run the binary
./server
```

## Docker Deployment

The easiest way to deploy the complete web application (frontend + backend) is using Docker. The multi-stage build process automatically builds the Vue.js frontend and Go backend into a single optimized container.

### Prerequisites

- Docker 20.10+ (with BuildKit support)
- Docker Compose (optional, but recommended)

### Method 1: Docker Compose (Recommended)

From the `Web` directory, simply run:

```bash
# Build and start the application
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the application
docker-compose down
```

The application will be available at `http://localhost:3000`.

### Method 2: Docker Build & Run

```bash
# Build the image from the Web directory
docker build -t nordgen-web:latest .

# Run the container
docker run -d \
  --name nordgen-web \
  -p 3000:3000 \
  --restart unless-stopped \
  nordgen-web:latest

# View logs
docker logs -f nordgen-web

# Stop and remove the container
docker stop nordgen-web && docker rm nordgen-web
```

### Advanced Configuration

**Custom Port Mapping:**
```bash
docker run -d -p 8080:3000 nordgen-web:latest
```

**Resource Limits:**
```bash
docker run -d \
  -p 3000:3000 \
  --memory="512m" \
  --cpus="1.0" \
  nordgen-web:latest
```

**Behind a Reverse Proxy (nginx/traefik):**

When deploying behind a reverse proxy, ensure the `X-Forwarded-For` header is properly set for accurate rate limiting. The application is configured to trust this header via the `ProxyHeader` setting.

Example nginx configuration:
```nginx
location / {
    proxy_pass http://localhost:3000;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Real-IP $remote_addr;
}
```

## Architecture

- **Core Store**: Maintains a thread-safe in-memory cache of NordVPN servers, refreshed every 5 minutes. It also handles static asset serving with pre-compressed Brotli support and ETag caching.
- **Validation**: Strict input validation ensures all data sent to upstream APIs or used in configuration generation is sanitized.
- **Performance**: Utilizes **Fiber's** zero-allocation routing and Go's native concurrency model to handle high throughput with minimal resource usage.

## Static Assets

The server looks for a `./public` directory to serve static frontend files. If an `index.html` is present, it is served for the root path and any unknown routes (SPA fallback), with the server data injected directly into the HTML to prevent an initial round-trip fetch.

## API Documentation

For detailed endpoint specifications, request/response formats, and validation rules, please refer to the [API.md](./API.md).