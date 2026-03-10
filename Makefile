.PHONY: init

init:
	@go mod tidy

generate:
	@templ generate

run: generate
	@go run .
