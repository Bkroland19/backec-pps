# Makefile for Point Prevalence Survey API

.PHONY: help build run test clean docker-build docker-run docker-stop swagger

# Default target
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker Compose services"
	@echo "  swagger      - Generate Swagger documentation"
	@echo "  deps         - Download dependencies"

# Build the application
build:
	go build -o main .

# Run the application
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f main
	go clean

# Download dependencies
deps:
	go mod download
	go mod tidy

# Generate Swagger documentation
swagger:
	swag init

# Docker commands
docker-build:
	docker build -t point-prevalence-survey .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Development setup
dev-setup: deps swagger
	@echo "Development setup complete!"
	@echo "Run 'make run' to start the application"
	@echo "Or run 'make docker-run' to start with Docker Compose"
