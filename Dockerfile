# Build stage
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN GOOS=linux go build -o koito_proxy ./cmd/proxy

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies and timezone data
RUN apk add --no-cache tzdata

# Copy binary from builder
COPY --from=builder /app/koito_proxy .

# Copy migrations
COPY migrations ./migrations

# Create data directory for database
RUN mkdir -p /app/data

# Volume for database persistence
VOLUME ["/app/data"]

# Expose Port 4112
EXPOSE 4112

# Copy entrypoint and make executable
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Run the application via entrypoint (entrypoint applies TZ at runtime)
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["./koito_proxy"]
