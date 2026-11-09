FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.work ./
COPY pkg/ pkg/
COPY proto/ proto/
COPY services/gateway-auth/ services/gateway-auth/

WORKDIR /src/services/gateway-auth
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/gateway-auth ./cmd/...

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/gateway-auth /app/gateway-auth

EXPOSE 8080
ENTRYPOINT ["/app/gateway-auth"]
