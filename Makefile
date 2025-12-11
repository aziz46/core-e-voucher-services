.PHONY: help build up down logs test clean restart healthcheck migrate

help:
	@echo "Available commands:"
	@echo "  make build              - Build Docker images for all services"
	@echo "  make up                 - Start all services with docker-compose"
	@echo "  make down               - Stop all services"
	@echo "  make logs               - View logs from all services"
	@echo "  make test               - Run all tests"
	@echo "  make test-unit          - Run unit tests only"
	@echo "  make test-integration   - Run integration tests"
	@echo "  make clean              - Clean up containers, volumes, and build artifacts"
	@echo "  make restart            - Restart all services"
	@echo "  make healthcheck        - Check health of all services"
	@echo "  make migrate            - Run database migrations"
	@echo "  make seed               - Seed database with test data"
	@echo "  make lint               - Run linter (golangci-lint)"
	@echo "  make fmt                - Format code"
	@echo "  make vet                - Run go vet"

build:
	@echo "Building Docker images..."
	docker-compose build

up:
	@echo "Starting services..."
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@make healthcheck

down:
	@echo "Stopping services..."
	docker-compose down

logs:
	docker-compose logs -f

logs-credit:
	docker-compose logs -f credit-service

logs-billing:
	docker-compose logs -f billing-service

logs-ppob:
	docker-compose logs -f ppob-core

test:
	@echo "Running all tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Test coverage:"
	@go tool cover -func=coverage.out | tail -1

test-unit:
	@echo "Running unit tests..."
	go test -v -short ./...

test-integration:
	@echo "Running integration tests..."
	go test -v -run Integration ./...

clean:
	@echo "Cleaning up..."
	docker-compose down -v
	rm -f coverage.out
	go clean ./...

restart: down up

healthcheck:
	@echo "Checking service health..."
	@curl -f http://localhost:8080/health > /dev/null 2>&1 && echo "✓ PPOB Core: OK" || echo "✗ PPOB Core: FAILED"
	@curl -f http://localhost:8081/health > /dev/null 2>&1 && echo "✓ Credit Service: OK" || echo "✗ Credit Service: FAILED"
	@curl -f http://localhost:8082/health > /dev/null 2>&1 && echo "✓ Billing Service: OK" || echo "✗ Billing Service: FAILED"

migrate:
	@echo "Running migrations..."
	docker-compose exec postgres psql -U postgres -d e_voucher -f /docker-entrypoint-initdb.d/001_init.sql
	docker-compose exec postgres psql -U postgres -d e_voucher -f /docker-entrypoint-initdb.d/002_seed.sql

seed:
	@echo "Seeding database..."
	docker-compose exec postgres psql -U postgres -d e_voucher -f /docker-entrypoint-initdb.d/002_seed.sql

lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
