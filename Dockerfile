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
# Use Alpine for a lightweight image with curl available for health checks
FROM alpine:latest

# Install ca-certificates and curl for HTTPS calls and health checks
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/tiger-fasttrack-card /app/tiger-fasttrack-card

# Change ownership to appuser
RUN chown -R appuser:appuser /app
USER appuser

# Expose port
EXPOSE 8080

# Add health check using curl
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl --fail --silent http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/tiger-fasttrack-card"]
