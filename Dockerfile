# syntax=docker/dockerfile:1

# ==============================================================================
# Stage 1: Build Stage (Compiles static Go binary for target OS/Architecture)
# ==============================================================================
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

WORKDIR /src

# Download dependencies first (cached unless go.mod changes)
COPY go.mod ./
RUN go mod download

# Copy application source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build arguments injected by Docker Buildx / GitHub Actions
ARG TARGETOS=linux
ARG TARGETARCH=arm64

# Compile static binary with stripped debug symbols for minimal image size
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ==============================================================================
# Stage 2: Hardened Runtime Stage (Non-root, minimal attack surface)
# ==============================================================================
FROM alpine:3.20

# Install CA certificates for HTTPS and create a non-root user/group (UID 10001)
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup -H -s /sbin/nologin

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /out/server /app/server

# Runtime environment defaults (overridable via Kubernetes Deployment env)
ENV APP_VERSION="1.0.0" \
    APP_ENV="production"

# Switch to non-root user
USER 10001:10001

# Expose unprivileged HTTP port
EXPOSE 8080

# Execute binary directly as PID 1 so SIGTERM is handled gracefully
ENTRYPOINT ["/app/server"]

