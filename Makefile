# Makefile for picow-led project

BINARY_NAME = ./cmd/picow-led

all: init build

init:
	@echo "Initializing project..."
	go mod tidy

# Build the binary
build:
	@echo "Building binary..."
	go build -o ./bin/picow-led $(BINARY_NAME)

# Run the application
run:
	@echo "Running server without building..."
	go run $(BINARY_NAME) -debug -log-format=text

# Clean up build files
clean:
	@echo "Cleaning up..."
	go clean
	rm -rf bin
