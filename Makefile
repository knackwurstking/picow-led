# Makefile for picow-led project

BINARY_NAME = ./cmd/picow-led
DATABASE_PATH = ./picow-led.db
TEST_DATABASE_PATH = ./services/picow-led.test.db

all: init build

init:
	@echo "Initializing project..."
	go mod tidy

# Build the binary
build:
	@echo "Building binary..."
	go build -o ./bin/picow-led $(BINARY_NAME)

# Run the application
dev-run:
	@echo "Running server without building..."
	go run $(BINARY_NAME) server -debug -log-format=text -database-path=$(DATABASE_PATH)

test:
	@echo "Running tests..."
	rm -rf $(TEST_DATABASE_PATH)
	go test -v ./...

# Clean up build files
clean:
	@echo "Cleaning up..."
	go clean
	rm -rf bin
	rm -rf $(TEST_DATABASE_PATH)
