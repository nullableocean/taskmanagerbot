.PHONY: up down build migrate-up migrate-down migrate-force migrate-version setup clean logs

# Main
up:
	docker compose up --build -d

down: 
	docker compose down

build:
	docker compose build

# Migration
migrate-up:
	./migrate.sh up

migrate-down:
	./migrate.sh down

migrate-force:
	@echo "Usage: make migrate-force VERSION=1"
	./migrate.sh force $(VERSION)

migrate-version:
	./migrate.sh version

# Setup and utility commands
setup:
	@echo "Setting up project..."
	@if [ ! -f .env ]; then \
		./setup.sh && \
		make up && \
		make migrate-up; \
	else \
		echo ".env already exists. Running build and migrations..." && \
		make up && \
		make migrate-up; \
	fi

logs:
	docker compose logs -f

# Help
help:
	@echo "Available commands:"
	@echo "  up                    - Start containers in detached mode"
	@echo "  down                  - Stop and remove containers"
	@echo "  build                 - Build containers"
	@echo "  migrate-up            - Run database migrations up"
	@echo "  migrate-down          - Run database migrations down"
	@echo "  migrate-force VERSION=1 - Force migration version"
	@echo "  migrate-version       - Show current migration version"
	@echo "  setup                 - Setup project (creates .env if needed)"
	@echo "  logs                  - Show container logs"
	@echo "  help                  - Show this help message"