.PHONY: all test lint build up down

all: test build

up:
	docker compose -f deploy/docker-compose.yml --env-file deploy/.env.example up -d

down:
	docker compose -f deploy/docker-compose.yml down

test:
	go test ./pkg/... ./proto/... ./services/...

build:
	go build -o bin/gateway-auth ./services/gateway-auth/cmd/...
	go build -o bin/trip-core ./services/trip-core/cmd/...
	go build -o bin/geo-places ./services/geo-places/cmd/...
