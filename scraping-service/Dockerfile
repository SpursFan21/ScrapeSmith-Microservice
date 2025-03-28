# Use official Golang image as builder
FROM golang:1.22-alpine AS builder
WORKDIR /src

# Install any required tools (e.g., git) for fetching modules
RUN apk add --no-cache git

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary for Linux (disable CGO for a static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scraping-service .

# Use a minimal runtime image
FROM alpine:latest
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /src/scraping-service .

# Ensure the binary has execution permissions
RUN chmod +x ./scraping-service

# Expose port
EXPOSE 3003

# Run the binary
CMD ["./scraping-service"]
