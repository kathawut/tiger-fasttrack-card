# Builder stage
# Use a Go version that matches your go.mod toolchain directive (go1.24.4)
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
# It's recommended to use a .dockerignore file to exclude unnecessary files
COPY . .

# Build a small, static, and optimized binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -o /app/tiger-fasttrack-card main.go

# Production stage
# Use a minimal and secure base image from Chainguard. It's similar to distroless
# but includes essential tools like curl for health checks.
FROM gcr.io/distroless/static-debian11:nonroot

# The Chainguard image already includes ca-certificates and runs as a non-root user by default.
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/tiger-fasttrack-card /app/tiger-fasttrack-card

# Copy any static assets if needed
# COPY --from=builder /app/static ./static

# Expose port
EXPOSE 8080

USER nonroot:nonroot

# Add health check
# Use curl (available in the chainguard image) instead of wget
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD curl --fail --silent --show-error http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/tiger-fasttrack-card"]
