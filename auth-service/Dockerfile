# Use official Golang image as builder with version 1.22 or higher
FROM golang:1.22 AS builder

# Set working directory
WORKDIR /auth-service

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source code
COPY . .

# Build the Go binary
RUN GOOS=linux GOARCH=amd64 go build -o /auth-service/auth-service

# Use a minimal runtime image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /auth-service/auth-service /root/

# Make sure the binary is executable
RUN chmod +x /root/auth-service

# Expose Fiber's default port
EXPOSE 3000

# Run the binary
CMD ["/root/auth-service"]
