FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.work ./
COPY pkg/ pkg/
COPY proto/ proto/
COPY services/ services/

WORKDIR /src/services/trip-core
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/trip-core ./cmd/...

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/trip-core /app/trip-core

EXPOSE 8081
ENTRYPOINT ["/app/trip-core"]
