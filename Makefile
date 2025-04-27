all: init build

BINARY_NAME := "picow-led"

clean:
	git clean -xfd

init:
	npm install
	go mod tidy -v

generate:
	go mod tidy -v
	templ generate 

watch-templ:
	templ generate --watch

dev:
	make generate
	which gow || (echo 'gow is not installed, install with: `go install github.com/mitranim/gow@latest`' && exit 1)
	gow -v -r run . server -a :8887

run:
	make generate
	go run . server -a :8887

test:
	go test -v ./...

build:
	make test
	make generate
	go build -v -ldflags="-w -s" -o bin/${BINARY_NAME}

generate-pwa-assets:
	npx pwa-assets-generator

# TODO: Add all "rpi-server-project" related commands here
