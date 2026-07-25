DATABASE_URL ?= postgres://rental:rental@localhost:5433/rental?sslmode=disable
PORT ?= 8081
SQLC ?= $(HOME)/go/bin/sqlc
GOOSE ?= go run github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: db-up db-down migrate-up migrate-down sqlc

up:
	docker compose up -d --wait postgres
	go build cmd/main.go
	./main
	
db-up: ## Start Postgres
	docker compose up -d postgres

db-down: ## Stop Postgres and remove its volume
	docker compose down -v

migrate-up: ## Apply migrations with goose
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" up

migrate-down: ## Roll back the last migration with goose
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" down

sqlc: ## Generate Go code from SQL
	$(SQLC) generate