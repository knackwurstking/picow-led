all: init build

BINARY_NAME := picow-led

SERVER_ADDR := :50835
SERVER_PATH_PREFIX := 

define SYSTEMD_SERVICE_FILE
[Unit]
Description=API for controlling my picow devices
After=network.target

[Service]
Environment="SERVER_ADDR=${SERVER_ADDR}"
Environment="SERVER_PATH_PREFIX=${SERVER_PATH_PREFIX}"
ExecStart=${BINARY_NAME} server

[Install]
WantedBy=default.target
endef

clean:
	git clean -xfd

init:
	go mod tidy -v
	git submodule init
	git submodule update --recursive
	cd frontend && npm ci 

dev:
	which gow || (echo 'gow is not installed, install with: `go install github.com/mitranim/gow@latest`' && exit 1)
	SERVER_PATH_PREFIX=${SERVER_PATH_PREFIX} \
		gow -e=go,json -v -r run ./cmd/${BINARY_NAME} server --addr ${SERVER_ADDR} --cache .

run:
	SERVER_PATH_PREFIX=${SERVER_PATH_PREFIX} \
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

UNAME := $(shell uname)
check-linux:
ifneq ($(UNAME), Linux)
	@echo 'This won’t work here since you’re not on Linux.'
	@exit 1
endif

export SYSTEMD_SERVICE_FILE
install: check-linux
	echo "$$SYSTEMD_SERVICE_FILE" > ${HOME}/.config/systemd/user/${BINARY_NAME}.service
	systemctl --user daemon-reload
	echo "--> Created a service file @ ${HOME}/.config/systemd/user/${BINARY_NAME}.service"
	sudo cp ./bin/${BINARY_NAME} /usr/local/bin/

start: check-linux
	systemctl --user restart ${BINARY_NAME}

stop: check-linux
	systemctl --user stop ${BINARY_NAME}

log: check-linux
	journalctl --user -u ${BINARY_NAME} --follow --output cat
