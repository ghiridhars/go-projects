# Variables
APP_NAME=basic
DOCKER_IMAGE_NAME=$(APP_NAME):latest

# Default target to build the binary and Docker image
all: build docker-build

# Build the Go application
build:
	@echo "Building the Go binary..."
	cd src && go build -o main .

# Run the Go application locally
run: build
	@echo "Running the Go application locally..."
	cd src && ./main

# Build the Docker image
docker-build:
	@echo "Building the Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) .

# Run the Docker container
docker-run:
	@echo "Running the Docker container..."
	@docker run --name postgres -e POSTGRES_USER=ghiri -e POSTGRES_PASSWORD=develop -e POSTGRES_DB=items -p 5432:5432 -d postgres:latest
	@docker run -d -p 8080:8080 --name $(APP_NAME) $(DOCKER_IMAGE_NAME)

# Stop the Docker container
docker-stop:
	@echo "Stopping the Docker container..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# Clean up binaries and Docker containers/images
clean:
	@echo "Cleaning up..."
	cd src && rm -f main
	docker rmi -f $(DOCKER_IMAGE_NAME) || true

.PHONY: all build run docker-build docker-run docker-stop clean
