# Makefile

# Default target
all: build

# Install dependencies
install:
	go mod tidy

# Build the application
build:
	go build -o bin/main main.go

# Run the application
run:
	go run main.go

# Run database migrations
migrate:
	@echo "Applying database migrations..."
	mysql -u root -p < db/schema.sql

# Clean build artifacts
clean:
	rm -rf bin/ tmp/

# Generate code from OpenAPI spec
generate:
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen -config tools/models.cfg.yml openapi/swagger.yml
	oapi-codegen -config tools/server.cfg.yml openapi/swagger.yml

# Run tests
test:
	go test ./...

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main main.go

# Docker build
docker-build:
	docker build -t h25s_20-backend .

# Docker run
docker-run:
	docker-compose up -d

# Docker stop
docker-stop:
	docker-compose down

# Run with auto-reload for development
dev:
	go run main.go

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  make install    - Install dependencies"
	@echo "  make build      - Build the application"
	@echo "  make run        - Run the application"
	@echo "  make migrate    - Run database migrations"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make generate   - Generate code from OpenAPI spec"
	@echo "  make test       - Run tests"
	@echo "  make build-prod - Build for production"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run - Run with Docker Compose"
	@echo "  make docker-stop - Stop Docker containers"
	@echo "  make dev        - Run with auto-reload"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Run linter"

.PHONY: all install build run migrate clean generate test build-prod docker-build docker-run docker-stop dev fmt lint help
