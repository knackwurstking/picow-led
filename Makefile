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
	npx tsc

test:
	go test -v ./...

run:
	make test
	make generate
	go run .

build:
	make test
	make generate
	go build -v -ldflags="-w -s" -o bin/${BINARY_NAME}

generate-pwa-assets:
	npx pwa-assets-generator

# TODO: Add all "rpi-server-project" related commands here
