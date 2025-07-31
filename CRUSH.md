# CRUSH.md - RMM23 Go Project Guide

## Build Commands
- Build: `make build`
- Test: `make test` (runs `go test ./...`)
- Lint: `make lint` (runs `golangci-lint run ./... --fix`)
- Vet: `make vet`
- Run: `make run`
- Race: `make race`

## Go Test Single File
`go test -v ./src/mod_pkg_name -run TestFunctionName`

## Code Style
- **Imports**: stdlib first, then third-party, then project imports
- **Naming**: CamelCase for exports, lowercase for internals
- **Modules**: Each module has const.go, errors.go, func.go, init.go, method.go, type.go, var.go
- **Error handling**: Use `mod_errors` package, return errors explicitly
- **Formatting**: Use `golangci-lint` with config from .golangci.json
- **Tags**: Struct tags aligned and sorted (tagalign enabled)

## Project Structure
- Main entry: `src/main.go`
- Module packages: `src/mod_*` (mod_db, mod_crypto, mod_git, etc.)
- Logging package: `src/l`
- Build target: `bin/rmm23`