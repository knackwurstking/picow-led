all: init build

include .env

BINARY_NAME := picow-led
SERVER_APP_NAME := ${BINARY_NAME}

clean:
	git clean -xfd

init:
	go mod tidy -v
	git submodule init
	git submodule update --recursive
	cd frontend && npm install

dev:
	which gow || (echo 'gow is not installed, install with: `go install github.com/mitranim/gow@latest`' && exit 1)
	gow -e=go,json -v -r run ./cmd/${BINARY_NAME} server --addr ${SERVER_ADDR} --cache .

run:
	go run ./cmd/${BINARY_NAME} server -a ${SERVER_ADDR} --cache .

test:
	go test -v ./...
	cd frontend && npm run check

build:
	cd frontend && npm run build
	rm -rf ./cmd/picow-led/frontend-build
	mkdir  ./cmd/picow-led/frontend-build
	cp -r ./frontend/.svelte-kit/output/prerendered/pages/* ./cmd/picow-led/frontend-build/
	cp -r ./frontend/.svelte-kit/output/client/* ./cmd/picow-led/frontend-build/
	go build -v --tags=frontend -o bin/${BINARY_NAME} ./cmd/${BINARY_NAME}

# NOTE: Standard systemd stuff

define SYSTEMD_SERVICE_FILE
[Unit]
Description=API for controlling my picow devices
After=network.target

[Service]
EnvironmentFile=%h/.config/${BINARY_NAME}/.env
ExecStart=${SERVER_APP_NAME} server

[Install]
WantedBy=default.target
endef

UNAME := $(shell uname)
check-linux:
ifneq ($(UNAME), Linux)
	@echo 'This won’t work here since you’re not on Linux.'
	@exit 1
endif

export SYSTEMD_SERVICE_FILE
install: check-linux
	echo "$$SYSTEMD_SERVICE_FILE" > ${HOME}/.config/systemd/user/${SERVER_APP_NAME}.service
	echo "--> Created a service file @ ${HOME}/.config/systemd/user/${SERVER_APP_NAME}.service"
	systemctl --user daemon-reload
	sudo cp ./bin/${SERVER_APP_NAME} /usr/local/bin/

start: check-linux
	systemctl --user restart ${SERVER_APP_NAME}

stop: check-linux
	systemctl --user stop ${SERVER_APP_NAME}

log: check-linux
	journalctl --user -u ${SERVER_APP_NAME} --follow --output cat
