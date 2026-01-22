.PHONY: all build test clean run-server run-collector

all: build

build:
	@echo "Building binaries..."
	@mkdir -p bin
	@go build -o bin/collector cmd/collector/main.go
	@go build -o bin/server cmd/server/main.go
	@echo "Build complete!"

test:
	@echo "Running tests..."
	@go test ./... -v

test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

run-server:
	@echo "Starting server..."
	@go run cmd/server/main.go

run-collector:
	@echo "Starting collector..."
	@go run cmd/collector/main.go

clean:
	@echo "Cleaning..."
	@rm -rf bin
	@rm -f coverage.out coverage.html
	@rm -f mqtt_catalog.db
	@echo "Clean complete!"

install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

fmt:
	@echo "Formatting code..."
	@go fmt ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

docker-build:
	@echo "Building Docker image..."
	@docker build -t mqtt-catalog:latest .

help:
	@echo "Available targets:"
	@echo "	make build			- Build both binaries"
	@echo "	make test			- Run all tests"
	@echo "	make test-coverage 	- Run tests with coverage report"
	@echo "	make run-server		- Run the server"
	@echo "	make run-collector	- Run the collector"
	@echo "	make clean			- Remove build artifacts"
	@echo "	make install-deps	- Install Go dependencies"
	@echo "	make fmt			- Format code"
	@echo "	make lint 			- Run linter"
	@echo " make docker-build	- Build Docker image"