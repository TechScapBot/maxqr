# Build stage - use golang:alpine with GOTOOLCHAIN=auto to download required Go version
FROM golang:alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Enable auto toolchain to download Go 1.24 if needed
ENV GOTOOLCHAIN=auto

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies (cached layer)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build with optimizations for production
# -trimpath: Remove file system paths from binary
# -ldflags: Strip debug info and set version
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-w -s -X main.Version=1.0.0" \
    -o /app/maxqr-api \
    ./cmd/server

# Final stage - minimal runtime image
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set timezone
ENV TZ=Asia/Ho_Chi_Minh

# Copy binary from builder
COPY --from=builder /app/maxqr-api /maxqr-api

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/maxqr-api"]
