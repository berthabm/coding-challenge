# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api-go ./cmd/server

# Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /out/api-go ./api-go

# Non-root
RUN adduser -D -H -u 10001 app \
  && chown -R app:app /app
USER app

ENV APP_HOST=0.0.0.0
ENV APP_PORT=8080

EXPOSE 8080

# El binario debe escuchar APP_HOST / APP_PORT (config por env)
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/health

CMD ["./api-go"]
