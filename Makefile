.PHONY: build test clean run migrate dev down logs help

# Variables
COMPOSE_FILE = docker-compose.yml
SERVICE_NAME = dashbeam
DB_NAME = dashbeam_db

# Help target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build all services
build: ## Build all service binaries
	@echo "Building all services..."
	go build -o bin/auth-service ./services/auth/cmd
	go build -o bin/quiz-service ./services/quiz/cmd
	go build -o bin/analytics-service ./services/analytics/cmd
	go build -o bin/reporting-service ./services/reporting/cmd
	go build -o bin/migrator ./cmd/migrator

# Docker compose targets
dev: ## Start development environment with docker-compose
	@echo "Starting development environment..."
	docker-compose -f $(COMPOSE_FILE) up -d
	@echo "Waiting for postgres to be ready..."
	@sleep 5
	@$(MAKE) migrate
	@$(MAKE) run

up:
	docker-compose -f $(COMPOSE_FILE) up -d

down:
	docker-compose -f $(COMPOSE_FILE) down

logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

# Database migration targets
migrate:
	@echo "Running database migrations..."
	go run ./cmd/migrator up

migrate-down:
	@echo "Rolling back database migrations..."
	go run ./cmd/migrator down

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	go run ./cmd/migrator create $(NAME)

# Database management
db-reset: ## Reset database (drop and recreate)
	@echo "Resetting database..."
	docker-compose -f $(COMPOSE_FILE) exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	docker-compose -f $(COMPOSE_FILE) exec postgres psql -U postgres -c "CREATE DATABASE $(DB_NAME);"
	@$(MAKE) migrate

db-shell: ## Connect to database shell
	docker-compose -f $(COMPOSE_FILE) exec postgres psql -U postgres -d $(DB_NAME)

run:
	@echo "Starting main server..."
	go run ./cmd/server

run-auth:
	go run ./services/auth/cmd

run-quiz:
	go run ./services/quiz/cmd

run-analytics:
	go run ./services/analytics/cmd

run-reporting:
	go run ./services/reporting/cmd


lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

tidy:
	go work sync
	go mod tidy -e ./...

clean:
	rm -rf bin/
	go clean ./...

clean-docker:
	docker-compose -f $(COMPOSE_FILE) down -v --remove-orphans
	docker system prune -f

dev-setup: ## Setup development environment from scratch
	@$(MAKE) clean-docker
	@$(MAKE) dev

dev-restart: ## Restart development environment
	@$(MAKE) down
	@$(MAKE) dev

# Docker build
docker-build: ## Build docker image
	docker build -t $(SERVICE_NAME):latest .
