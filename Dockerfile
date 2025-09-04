# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/server

# Release stage
FROM alpine:latest

# Install grpc_health_probe and curl for health checks
RUN apk add --no-cache ca-certificates curl && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.19/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Expose both HTTP REST API and gRPC ports
EXPOSE 8080 9090

# Set environment variables
ENV PORT=8080
ENV GRPC_PORT=9090

# Run the gRPC server
CMD ["/app/server"]
