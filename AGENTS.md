# AGENTS.md

## Build/Lint/Test Commands
- `make init` - Initialize project dependencies
- `make generate` - Generate templ files
- `make build` - Build the binary
- `make test` - Run all tests
- `go test -v ./path/to/package` - Run a single test package
- `go test -v ./... -run TestFunctionName` - Run a specific test function

## Code Style Guidelines
- Use Go standard formatting (`gofmt` or `goimports`)
- Follow Go naming conventions (PascalCase for exported names, camelCase for unexported)
- Use descriptive variable and function names
- Keep functions small and focused on a single responsibility
- Use clear error messages with context
- Prefer explicit error handling over panics
- Group imports by standard library, external dependencies, and internal packages
- Use `//nolint` comments for intentional exceptions to linters
- Follow the "don't repeat yourself" principle with reusable components

## Testing
- Write unit tests for all business logic
- Use table-driven tests where appropriate
- Test error conditions and edge cases
- Use `require` for assertions that should stop execution if failed
- Use `assert` for assertions that should continue execution if failed

## Project Structure
- All Go code in `*.go` files
- Templates in `*.templ` files with generated `*_templ.go` files
- Web assets in `/assets/`
- UI components in `/handlers/components/`
- Services in `/services/`
- Models in `/models/`