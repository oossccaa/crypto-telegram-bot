# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o crypto-bot .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/crypto-bot .

# Copy config example file (users need to provide actual config.yaml via environment or volume mount)
COPY --from=builder /app/config.yaml.example ./config.yaml.example

# Expose port (if needed for health checks)
EXPOSE 8080

# Run the binary
CMD ["./crypto-bot"]