.PHONY: init generate install-tailwind

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
	@go build -v -o bin/picow-led .
