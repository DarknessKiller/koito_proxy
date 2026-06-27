# Frontend build stage
FROM node:26-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Go build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/internal/admin/dist ./internal/admin/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o koito_proxy ./cmd/app

# App image
FROM scratch
WORKDIR /app
COPY --from=builder /app/koito_proxy koito_proxy
VOLUME ["/app/data"]

ENV PROXY_DB=/app/data/koito_proxy.db
ENV PROXY_PORT=4112

EXPOSE 4112
ENTRYPOINT ["./koito_proxy"]
