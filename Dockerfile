# Stage 1: Build
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/auth-service ./cmd/http/main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/auth-service .
# Copy .env file if needed (usually passed via docker-compose env_file)
# COPY .env .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./auth-service"]
