include .env
MIGRATION_PATH= ./migrations

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

# Create migration
migration:
	@migrate create -seq -ext sql -dir ${MIGRATION_PATH} $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path=${MIGRATION_PATH} -database=${DB_MIGRATOR_ADDR} up

migrate-down:
	@migrate -path=${MIGRATION_PATH} -database=${DB_MIGRATOR_ADDR} down $(filter-out $@,$(MAKECMDGOALS))

migrate-force:
	@migrate -path=${MIGRATION_PATH} -database=${DB_MIGRATOR_ADDR} force $(filter-out $@,$(MAKECMDGOALS))

.PHONY: all build run test clean watch docker-run docker-down itest migration migrate-up migrate-down migrate-force
