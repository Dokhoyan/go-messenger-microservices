.PHONY: help build run test clean docker-build docker-run

build:
	@echo "Building application..."
	@go build -o bin/api-gateway ./cmd/server

run:
	@echo "Running application..."
	@go run ./cmd/server

clean:
	@echo "Cleaning..."
	@rm -rf bin/

docker-build:
	@echo "Building Docker image..."
	@docker-compose build

docker-run:
	@echo "Running with Docker Compose..."
	@docker-compose up -d

