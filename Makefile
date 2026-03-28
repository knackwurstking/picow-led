.PHONY: all init generate install-tailwind macos-install macos-update

BINARY_NAME := picow-led
BIN_DIR := ./bin
INSTALL_PATH := /usr/local/bin

SERVICE_FILE := $(HOME)/Library/LaunchAgents/com.$(BINARY_NAME).plist

APP_DATA := $(HOME)/Library/Application Support/$(BINARY_NAME)
LOG_FILE := $(APP_DATA)/$(BINARY_NAME).log

all: generate init build

generate: install-tailwind
	@templ generate
	@npx tailwindcss -i ./internal/assets/public/css/input.css -o ./internal/assets/public/css/output.css --minify

init: generate
	@go mod tidy
	@npm install

install-tailwind:
	@test -f package.json || npm init -y
	@npm install --save-dev tailwindcss postcss autoprefixer

run: generate
	@gow -e=go,js,css -r run .

build: generate
	@go build -v -o bin/$(BINARY_NAME) .

define LAUNCHCTL_PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>Label</key>
        <string>com.$(BINARY_NAME)</string>

        <key>ProgramArguments</key>
        <array>
                <string>$(INSTALL_PATH)/$(BINARY_NAME)</string>
        </array>

        <key>RunAtLoad</key>
        <true/>

        <key>KeepAlive</key>
        <true/>

        <key>StandardOutPath</key>
        <string>$(LOG_FILE)</string>

        <key>StandardErrorPath</key>
        <string>$(LOG_FILE)</string>

        <key>EnvironmentVariables</key>
        <dict>
                <key>DB_PATH</key>
                <string>$(HOME)</string>
                <key>SERVER_PATH_PREFIX</key>
                <string>/$(BINARY_NAME)</string>
        </dict>
</dict>
</plist>
endef

export LAUNCHCTL_PLIST

macos-install: all
	@echo "Installing $(BINARY_NAME) for macOS..."
	mkdir -p $(INSTALL_PATH)
	sudo cp $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	echo "$$LAUNCHCTL_PLIST" > $(SERVICE_FILE)
	@echo "$(BINARY_NAME) installed successfully"

macos-update: all
	sudo cp $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Using launchctl command for restarting the service..."
