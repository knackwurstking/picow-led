BINARY_NAME = picow-led
APP_DATA = $(HOME)/Library/Application Support/picow-led
DATABASE_PATH = $(APP_DATA)/picow-led.db
TEST_DATABASE_PATH = ./picow-led.test.db

all: init build

init:
	@echo "Initializing project..."
	templ generate
	go mod tidy

generate:
	@echo "Generating templ files..."
	templ generate

# Build the binary
build:
	@echo "Building binary..."
	go build -o ./bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

# Run the application
run:
	@echo "Running server without building..."
	make init
	make generate
	go run ./cmd/$(BINARY_NAME) server \
		-debug \
		-log-format text \
		-database-path $(TEST_DATABASE_PATH) \
		-path-prefix /picow-led

test:
	@echo "Running tests..."
	rm -rf $(TEST_DATABASE_PATH)
	go test -v ./...

# Clean up build files
clean:
	@echo "Cleaning up..."
	go clean
	git clean -xfd

define LAUNCHCTL_PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.picow-led</string>

	<key>ProgramArguments</key>
	<array>
		<string>/usr/local/bin/picow-led</string>
		<string>server</string>
		<string>-path-prefix</string>
		<string>/picow-led</string>
		<string>-addr</string>
		<string>:50836</string>
		<string>-log-format</string>
		<string>text</string>
		<string>-database-path</string>
		<string>$(DATABASE_PATH)</string>
	</array>

	<key>RunAtLoad</key>
	<true/>

	<key>KeepAlive</key>
	<true/>

	<key>StandardOutPath</key>
	<string>$(APP_DATA)/picow-led.log</string>

	<key>StandardErrorPath</key>
	<string>$(APP_DATA)/picow-led.log</string>
</dict>
</plist>
endef

export LAUNCHCTL_PLIST

macos-install:
	@echo "Installing picow-led for macOS..."
	mkdir -p /usr/local/bin
	sudo cp ./bin/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	mkdir -p $(APP_DATA)
	@echo "$$LAUNCHCTL_PLIST" > ~/Library/LaunchAgents/com.picow-led.plist
	@echo "picow-led installed successfully"

macos-start-service:
	@echo "Starting picow-led service..."
	launchctl load -w ~/Library/LaunchAgents/com.picow-led.plist
	launchctl start com.picow-led

macos-stop-service:
	@echo "Stopping picow-led service..."
	launchctl stop com.picow-led
	launchctl unload -w ~/Library/LaunchAgents/com.picow-led.plist

macos-restart-service:
	@echo "Restarting picow-led service..."
	make macos-stop-service
	make macos-start-service

macos-print-service:
	@echo "picow-led service information:"
	launchctl print gui/$$(id -u)/com.picow-led || echo "Service not loaded or running"

macos-watch-service:
	@echo "picow-led watch server logs @ \"$(APP_DATA)/picow-led.log\":"
	@if [ -f "$(APP_DATA)/picow-led.log" ]; then \
		echo "Watching logs... Press Ctrl+C to stop"; \
		tail -f "$(APP_DATA)/picow-led.log"; \
	else \
		echo "Log file not found. Make sure the service is running or has been started."; \
		echo "Log file path: $(APP_DATA)/picow-led.log"; \
	fi

macos-update: all
	make macos-stop-service
	make macos-install
	make macos-start-service
	make macos-watch-service
