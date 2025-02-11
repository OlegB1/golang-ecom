include .env

MIGRATIONS_PATH = ./cmd/migrations

build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -v ./...

run:
	@go run cmd/main.go

migrate-create:
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(word 2, $(MAKECMDGOALS))

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ARRD) up

migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ARRD) down
