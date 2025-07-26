.PHONY: run dev build test clean docker-up docker-down db-up db-down db-migrate

# Development
dev:
	go run main.go

run:
	go build -o bin/tiger-fasttrack-card main.go && ./bin/tiger-fasttrack-card

build:
	go build -o bin/tiger-fasttrack-card main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

# Database
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

db-up: docker-up
	@echo "PostgreSQL is starting up..."
	@sleep 5
	@echo "PostgreSQL is ready!"

db-down: docker-down

# Dependencies
deps:
	go mod tidy
	go mod download

# Format and lint
fmt:
	go fmt ./...

lint:
	golangci-lint run

# Database migrations (when you add models)
# db-migrate:
# 	go run main.go -migrate

# Help
help:
	@echo "Available commands:"
	@echo "  dev         - Run the application in development mode"
	@echo "  run         - Build and run the application"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-up   - Start Docker containers"
	@echo "  docker-down - Stop Docker containers"
	@echo "  db-up       - Start PostgreSQL database"
	@echo "  db-down     - Stop PostgreSQL database"
	@echo "  deps        - Download dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linter"
