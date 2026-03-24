.PHONY: help postgres-up postgres-down auth-up auth-build auth-migrate auth-migrate-create auth-migrate-down auth-sqlc build lint test clean

help:
	@echo "Available commands:"
	@echo "  make postgres-up        - Start postgres container"
	@echo "  make postgres-down     - Stop postgres container"
	@echo "  make auth-up           - Run auth service"
	@echo "  make auth-build        - Build auth service"
	@echo "  make auth-migrate      - Run auth service migrations"
	@echo "  make auth-migrate-create name=foo - Create migration for auth"
	@echo "  make auth-migrate-down - Rollback auth migrations"
	@echo "  make auth-sqlc         - Generate sqlc for auth"
	@echo "  make build             - Build all services"
	@echo "  make lint              - Run linters"
	@echo "  make test              - Run tests"
	@echo "  make clean             - Clean build artifacts"

# Database
postgres-up:
	docker compose up -d postgres
	@echo "Waiting for postgres to be ready..."
	@sleep 3

postgres-down:
	docker compose down

# Auth Service
auth-up:
	cd services/auth && go run cmd/main.go

auth-build:
	cd services/auth && go build -o bin/auth ./cmd

auth-migrate:
	cd services/auth && goose -dir=./migrations postgres "user=$(shell grep POSTGRES_USER .env | cut -d= -f2 || echo root) password=$(shell grep POSTGRES_PASSWORD .env | cut -d= -f2 || echo secret) host=localhost port=$(shell grep POSTGRES_PORT .env | cut -d= -f2 || echo 5432) dbname=$(shell grep POSTGRES_DB .env | cut -d= -f2 || echo koffee_auth)" up

auth-migrate-create:
	@echo "Usage: make auth-migrate-create name=create_users_table"
	@if [ -n "$(name)" ]; then \
		cd services/auth && goose create $(name) sql; \
	fi

auth-migrate-down:
	cd services/auth && goose -dir=./migrations postgres "user=$(shell grep POSTGRES_USER .env | cut -d= -f2 || echo root) password=$(shell grep POSTGRES_PASSWORD .env | cut -d= -f2 || echo secret) host=localhost port=$(shell grep POSTGRES_PORT .env | cut -d= -f2 || echo 5432) dbname=$(shell grep POSTGRES_DB .env | cut -d= -f2 || echo koffee_auth)" down

auth-sqlc:
	cd services/auth && sqlc generate

# Build all
build: auth-build

# Lint
lint:
	golangci-lint run

# Test
test:
	go test ./...

# Clean
clean:
	rm -rf services/*/bin
