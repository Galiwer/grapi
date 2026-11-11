# DevOps Demo API (`grapi`)

A lightweight Go 1.22 REST API designed to demonstrate an end-to-end GitOps CI/CD pipeline using **GitHub Actions**, **GitHub Container Registry (GHCR)**, **k3s (ARM64)**, **Argo CD**, and **Traefik Ingress**.

## Directory Structure

```text
application/
├── cmd/
│   └── server/
│       └── main.go              # HTTP server entry point & graceful SIGTERM shutdown
├── internal/
│   └── handlers/
│       ├── handlers.go          # HTTP route handlers (/, /health, /ready, /version, /api/message)
│       └── handlers_test.go     # Unit tests for /health and /version
├── go.mod                       # Go module definition (github.com/Galiwer/grapi)
├── Dockerfile                   # Multi-stage, non-root (UID 10001) ARM64-compatible Dockerfile
├── .dockerignore                # Excludes test files and local metadata from build context
└── README.md
```

## Exposed HTTP Endpoints (Port `8080`)

| Method & Path      | Purpose                                             | Example Response                                                                                                  |
| ------------------ | --------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `GET /`            | Application overview & active version               | `{"application":"DevOps Demo API","description":"REST API for demonstrating CI/CD and GitOps","version":"1.0.0"}` |
| `GET /health`      | Kubernetes `livenessProbe` endpoint                 | `{"status":"healthy"}`                                                                                            |
| `GET /ready`       | Kubernetes `readinessProbe` endpoint                | `{"status":"ready"}`                                                                                              |
| `GET /version`     | Detailed version, environment & UTC timestamp       | `{"application":"DevOps Demo API","environment":"local","timestamp":"...","version":"1.0.0"}`                     |
| `GET /api/message` | Sample business payload endpoint for update testing | `{"message":"Hello! Your Go REST API is running successfully.","status":"success"}`                               |

## Environment Variables

| Variable      | Default | Description                             |
| ------------- | ------- | --------------------------------------- |
| `APP_VERSION` | `1.0.0` | Version reported by `/` and `/version`  |
| `APP_ENV`     | `local` | Environment name reported by `/version` |

## Running & Testing Locally

```bash
# Run unit tests
cd application
go test ./... -v

# Run the server locally on :8080
go run ./cmd/server

# Build & run with Docker (ARM64)
docker build --platform linux/arm64 -t grapi:local .
docker run --rm -p 8080:8080 grapi:local
```
