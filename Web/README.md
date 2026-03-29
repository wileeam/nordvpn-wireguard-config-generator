# NordGen Web Application

The NordGen web application is a full-stack deployment that combines a Vue.js single-page frontend with a high-performance Go backend into a single container. It exposes an API for server discovery, WireGuard key exchange, and configuration generation.

## Structure

| Directory | Description |
|-----------|-------------|
| [`web-Frontend/`](./web-Frontend/) | Vue 3 + Vite + Tailwind CSS single-page application |
| [`web-Backend/`](./web-Backend/) | Go + Fiber backend API server |

For backend-specific development and API documentation see [`web-Backend/README.md`](./web-Backend/README.md).

## Docker Deployment

The easiest way to deploy the complete web application is using Docker. The multi-stage `Dockerfile` in this directory automatically builds the Vue.js frontend and the Go backend into a single optimised container.

### Prerequisites

- Docker 20.10+ (with BuildKit support)
- Docker Compose (optional, but recommended)

### Method 1: Docker Compose (Recommended)

From this `Web/` directory, run:

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
# Build the image from the Web/ directory
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

When deploying behind a reverse proxy, ensure the `X-Forwarded-For` header is passed correctly for accurate rate limiting. The backend is already configured to trust this header via the `ProxyHeader` setting.

Example nginx configuration:
```nginx
location / {
    proxy_pass http://localhost:3000;
    proxy_set_header Host $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Real-IP $remote_addr;
}
```
