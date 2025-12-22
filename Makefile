# 1. Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# 2. Variable helper (defaults to disable SSL for local dev if not specified)
# We append ?sslmode=disable if it's not already in the URL to prevent "SSL off" errors
MIGRATE_CMD=migrate -path db/migrations -database "$(DATABASE_URL)?sslmode=require"

# --- Commands ---

# Run the server (Loading .env automatically)
run:
	go run cmd/api/main.go

# Create a new migration file
# Usage: make create-migration name=some_name
create-migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

# Apply all up migrations
migrate-up:
	@echo "Running Up Migrations..."
	$(MIGRATE_CMD) up

# Revert the last migration
migrate-down:
	@echo "Running Down Migration..."
	$(MIGRATE_CMD) down 1

force:
	$(MIGRATE_CMD) force $(version)

# Clean up dependencies
tidy:
	go mod tidy