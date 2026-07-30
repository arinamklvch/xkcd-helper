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

db-down: ## Stop Postgres
	docker compose down

db-kill: ## Removes volume
	docker compose down -v
	
sqlc: ## Generate Go code from SQL
	$(SQLC) generate