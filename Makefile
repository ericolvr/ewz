# ANSI color codes
COLOR_RESET=\033[0m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m
COLOR_RED=\033[31m

MAIN_PATH=cmd/main.go
DB_CONTAINER=ewz-postgres
DB_VOLUME_NAME=ewz_postgres_data

.PHONY: help install run db-start db-stop db-clean db-migrate db-reset build

help:
	@echo ""
	@echo "  $(COLOR_YELLOW)Available targets:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)install$(COLOR_RESET)    - Install dependencies"
	@echo "  $(COLOR_GREEN)db-start$(COLOR_RESET)   - Start Postgres container"
	@echo "  $(COLOR_GREEN)db-stop$(COLOR_RESET)    - Stop Postgres container"
	@echo "  $(COLOR_GREEN)db-clean$(COLOR_RESET)   - Remove container and volume"
	@echo "  $(COLOR_GREEN)db-migrate$(COLOR_RESET) - Apply SQL schema"
	@echo "  $(COLOR_GREEN)db-reset$(COLOR_RESET)   - Recreate database"
	@echo "  $(COLOR_GREEN)run$(COLOR_RESET)        - Run development server"
	@echo "  $(COLOR_GREEN)build$(COLOR_RESET)      - Build binary"
	@echo ""

install:
	@echo "$(COLOR_YELLOW)Installing Go dependencies...$(COLOR_RESET)"
	go mod tidy
	@echo "$(COLOR_GREEN)Dependencies ready$(COLOR_RESET)"

db-start:
	@echo "$(COLOR_YELLOW)Starting Postgres container...$(COLOR_RESET)"
	@if [ ! -f .env ]; then cp .env-sample .env; fi
	docker compose --env-file .env up postgres -d
	@echo "$(COLOR_GREEN)Database started$(COLOR_RESET)"

db-stop:
	@echo "$(COLOR_YELLOW)Stopping Postgres container...$(COLOR_RESET)"
	docker compose stop postgres
	docker compose rm -f postgres
	@echo "$(COLOR_GREEN)Database stopped$(COLOR_RESET)"

db-clean:
	@echo "$(COLOR_YELLOW)Cleaning database data...$(COLOR_RESET)"
	docker compose stop postgres
	docker compose rm -f postgres
	docker volume rm $(DB_VOLUME_NAME) 2>/dev/null || true
	@echo "$(COLOR_GREEN)Database cleaned$(COLOR_RESET)"

db-migrate:
	@echo "$(COLOR_YELLOW)Applying schema...$(COLOR_RESET)"
	@if [ ! -f .env ]; then cp .env-sample .env; fi
	@until docker exec $(DB_CONTAINER) pg_isready -U $$(grep DB_USER .env | cut -d'=' -f2) -d $$(grep DB_NAME .env | cut -d'=' -f2) > /dev/null 2>&1; do \
		echo "$(COLOR_YELLOW)Waiting for database...$(COLOR_RESET)"; \
		sleep 2; \
	done
	docker exec -i $(DB_CONTAINER) psql -v ON_ERROR_STOP=1 -U $$(grep DB_USER .env | cut -d'=' -f2) -d $$(grep DB_NAME .env | cut -d'=' -f2) < scripts/init.sql
	@echo "$(COLOR_GREEN)Schema applied$(COLOR_RESET)"

db-reset: db-clean db-start db-migrate
	@echo "$(COLOR_GREEN)Database reset completed$(COLOR_RESET)"

run:
	@echo "$(COLOR_YELLOW)Starting server...$(COLOR_RESET)"
	@if [ ! -f .env ]; then \
		if [ -f .env-sample ]; then \
			cp .env-sample .env; \
			echo "$(COLOR_BLUE)Created .env from .env-sample$(COLOR_RESET)"; \
		else \
			echo "$(COLOR_RED).env-sample not found$(COLOR_RESET)"; \
			exit 1; \
		fi; \
	fi
	go run $(MAIN_PATH)

build:
	@echo "$(COLOR_YELLOW)Building application...$(COLOR_RESET)"
	go build -o bin/backend $(MAIN_PATH)
	@echo "$(COLOR_GREEN)Build completed$(COLOR_RESET)"
