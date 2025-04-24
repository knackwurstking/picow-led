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
	go generate -v
	cd web/js && make

generate-dist:
	make generate
	rm -rf dist
	cp -r public dist
	go run -v . generate ./dist

build:
	SERVER_PATH_PREFIX= make generate-dist
	go build -v -ldflags="-w -s" -o bin/${BINARY_NAME}

generate-pwa-assets:
	npx pwa-assets-generator
