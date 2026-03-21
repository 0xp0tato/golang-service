.PHONY: help run-data-generator run-producer run-consumer run-all build build-data-generator clean deps-data-generator deps-all

help:
	@echo "Available targets:"
	@echo "  run-data-generator     	- Run the data-generator service"
	@echo "  run-all       			- Run both services"
	@echo "  build         			- Build all services"
	@echo "  build-data-generator	- Build data-generator service only"
	@echo "  clean         			- Remove build artifacts"

run-data-generator:
	@cd services/data-generator && go run main.go

run-producer:
	@cd services/producer && go run main.go

run-consumer:
	@cd services/consumer && go run main.go

run-all:
	@echo "Running all services..."
	@cd services/data-generator && go run main.go &
	@cd services/producer && go run main.go &
	@cd services/consumer && go run main.go &
	@wait
	@echo "All services started"

build: build-data-generator
	@echo "All services built successfully!"

build-data-generator:
	@echo "Building data-generator service..."
	@mkdir -p bin
	@cd services/data-generator && go build -o ../../bin/data-generator main.go

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin
	@cd services/data-generator && go clean
	@echo "Clean complete!"

# Dependency management helpers
deps-data-generator:
	@echo "Updating data-generator service dependencies..."
	@cd services/data-generator && go mod tidy && go mod download

deps-all: deps-data-generator
	@echo "All dependencies updated!"