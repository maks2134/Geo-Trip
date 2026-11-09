FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.work ./
COPY pkg/ pkg/
COPY proto/ proto/
COPY services/geo-places/ services/geo-places/

WORKDIR /src/services/geo-places
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/geo-places ./cmd/...

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/geo-places /app/geo-places

EXPOSE 8082
ENTRYPOINT ["/app/geo-places"]
