SQLC ?= $(HOME)/go/bin/sqlc
SWAG ?= $(HOME)/go/bin/swag

.PHONY: db-up db-down migrate-up migrate-down sqlc swag

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

swag: ## Generate Swagger docs
	$(SWAG) init -g cmd/main.go
