.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down sqlc

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go

run: ## Run the application
	@echo "Running..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.txt -covermode=atomic ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.txt

docker-up: ## Start Docker Compose services
	@echo "Starting Docker Compose services..."
	@docker-compose up -d

docker-down: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	@docker-compose down

docker-logs: ## Show Docker Compose logs
	@docker-compose logs -f

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build

sqlc: ## Generate SQLC code
	@echo "Generating SQLC code..."
	@sqlc generate

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@migrate create -ext sql -dir migrations -seq $(NAME)

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/resell?sslmode=disable" up

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	@migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/resell?sslmode=disable" down 1

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	@migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/resell?sslmode=disable" force $(VERSION)

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

dev: docker-up ## Start development environment
	@echo "Development environment started"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "Keycloak: http://localhost:8180"
	@echo ""
	@echo "Run 'make migrate-up' to apply database migrations"
	@echo "Run 'make run' to start the API server"
