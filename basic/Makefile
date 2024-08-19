.PHONY: build run docker-build docker-run docker-run-dev clean

# Build the Go application binary
build:
	@echo "Building the Go application..."
	docker-compose run --rm app go build -o main .

# Run the Go application (dev mode)
run:
	@echo "Running the Go application..."
	docker-compose run --rm app go run main.go

# Build the Docker image
docker-build:
	@echo "Building the Docker image..."
	docker-compose build

# Run the Docker containers
docker-run:
	@echo "Starting the Docker containers..."
	docker-compose up -d

# Run in dev mode with live reload
docker-run-dev:
	@echo "Running Docker containers in dev mode with live reload..."
	docker-compose up --build

# Clean the built files and containers
clean:
	@echo "Cleaning up..."
	docker-compose down
	rm -f main
