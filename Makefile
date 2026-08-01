SQLC ?= $(HOME)/go/bin/sqlc
SWAG ?= $(HOME)/go/bin/swag
GOLANGCI_LINT ?= go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

.PHONY: up db-up db-down db-kill sqlc swag fmt-check lint vet build ci

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

fmt-check: ## Check Go formatting
	test -z "$$(gofmt -l .)"

lint: ## Run golangci-lint
	GOCACHE=/tmp/go-build $(GOLANGCI_LINT) run ./...

vet: ## Run go vet
	GOCACHE=/tmp/go-build go vet ./...

build: ## Build all packages
	GOCACHE=/tmp/go-build go build ./...

ci: fmt-check lint vet build ## Run local CI checks
