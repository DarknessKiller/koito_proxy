# Build stage
FROM golang:1.26.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o koito_proxy ./cmd/app

# App image
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache tzdata
COPY --from=builder /app/koito_proxy .
COPY migrations ./migrations
RUN mkdir -p /app/data
VOLUME ["/app/data"]
EXPOSE 4112
ENTRYPOINT ["./koito_proxy"]
