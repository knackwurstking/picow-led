# Makefile for picow-led project

BINARY_NAME = ./cmd/picow-led
TEST_DATABASE_PATH = ./services/picow-led.test.db

all: init build

init:
	@echo "Initializing project..."
	go mod tidy

generate:
	@echo "Generating templ files..."
	templ generate

# Build the binary
build:
	@echo "Building binary..."
	go build -o ./bin/picow-led $(BINARY_NAME)

# Run the application
run:
	@echo "Running server without building..."
	make init
	make generate
	go run $(BINARY_NAME) server -debug -log-format text -database-path ./picow-led.db -path-prefix /picow-led

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
