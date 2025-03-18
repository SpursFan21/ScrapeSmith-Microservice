# Use official Golang image as builder
FROM golang:1.21 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o auth-service

# Use a minimal runtime image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/auth-service .

# Expose Fiber's default port
EXPOSE 3000

# Run the binary
CMD ["./auth-service"]
