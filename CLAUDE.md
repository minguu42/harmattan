# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Harmattan is a Go-based REST API service that provides authentication and project management functionality. It uses OpenAPI 3.0 specification with code generation via ogen.

## Key Architecture

- **API Layer**: `api/` contains the main application entry point, handlers, and usecases
- **Domain Layer**: `internal/domain/` contains core business entities (Project, User)
- **Database Layer**: `internal/database/` contains database access logic using GORM
- **Generated Code**: `internal/openapi/` contains auto-generated OpenAPI client/server code
- **Shared Libraries**: `lib/` contains reusable utilities (logging, clock, ID generation, etc.)

## Core Commands

### Code Generation
```bash
# Generate OpenAPI code from spec (run from api/ directory)
go generate ./...
```

### Build & Run
```bash
# Build the API server
go build -o bin/api ./api

# Run the API server
go run ./api/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./lib/clock
```

### Linting & Code Quality
```bash
# Run staticcheck (configured in staticcheck.conf)
go tool staticcheck ./...

# Format code
go tool goimports -w .
```

## OpenAPI Integration

The project uses ogen for OpenAPI code generation:
- API specification: `api/openapi.yaml`
- Generated code: `internal/openapi/`
- Code generation command is in `api/main.go` via `//go:generate`

## Testing Architecture

- Uses testify for assertions and test utilities
- Custom test helpers in `lib/` packages (clocktest, databasetest, idgentest)
- Database tests use testcontainers for integration testing
- HTTP tests use ikawaha/httpcheck for API testing

## Key Dependencies

- **Web Framework**: Uses generated ogen server code
- **Database**: GORM with MySQL driver
- **Authentication**: JWT via golang-jwt/jwt
- **Testing**: testify, testcontainers, httpcheck
- **Logging**: Custom applog wrapper

## Development Notes

- Timezone is set to JST (Japan Standard Time) in main.go
- Uses ULID for ID generation
- Custom middleware for access logging, recovery, and request ID context
- Database models use GORM conventions
- Authentication uses JWT with custom security handler