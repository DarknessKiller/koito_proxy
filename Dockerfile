# Build stage
FROM golang:1.27-alpine AS builder
RUN apk add --no-cache tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o koito_proxy ./cmd/app

# App image
FROM scratch
WORKDIR /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/koito_proxy koito_proxy
VOLUME ["/app/data"]

ENV PROXY_DB=/app/data/koito_proxy.db
ENV PROXY_PORT=4112
ENV TZ=Asia/Kuala_Lumpur

EXPOSE ${PROXY_PORT}
ENTRYPOINT ["./koito_proxy"]
