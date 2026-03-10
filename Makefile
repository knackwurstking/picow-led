.PHONY: init generate install-tailwind

init:
	@go mod tidy
	@npm install

install-tailwind:
	@test -f package.json || npm init -y
	@npm install --save-dev tailwindcss postcss autoprefixer

generate: install-tailwind
	@templ generate
	@npx tailwindcss -i ./internal/assets/public/css/input.css -o ./internal/assets/public/css/output.css --minify

run: generate
	@go run .
