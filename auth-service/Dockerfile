# Stage 1: Builder
FROM golang:1.22-alpine AS builder
WORKDIR /src

# Install any required tools (e.g., git) for fetching modules
RUN apk add --no-cache git

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary for Linux (disable CGO for a static binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-service .

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /src/auth-service .

# Ensure the binary has execution permissions
RUN chmod +x ./auth-service

# Expose the port your service listens on (Fiber's default is 3000)
EXPOSE 3000

# Start the service
CMD ["/app/auth-service"]
