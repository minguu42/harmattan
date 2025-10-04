# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Harmattan is a Go-based REST API service that provides authentication and project management functionality. It uses OpenAPI 3.0 specification with code generation via ogen.

## Key Architecture

- **cmd/api/**: Main application entry point
- **doc/**: API documentation including OpenAPI specification
- **infra/**: Infrastructure configuration (MySQL, etc.)
- **internal/**: Private application and library code
  - **api/**: API layer components
    - **handler/**: HTTP request handlers
    - **middleware/**: HTTP middleware
    - **openapi/**: Auto-generated OpenAPI client/server code
    - **usecase/**: Business logic layer
  - **auth/**: Authentication logic
  - **database/**: Database access layer using GORM
    - **databasetest/**: Test utilities for database testing
  - **domain/**: Core business entities (Project, Task, Step, Tag, User)
  - **factory/**: Dependency injection and initialization
  - **lib/**: Shared utility libraries
    - **alog/**: Application logging (slog wrapper)
    - **clock/**: Time utilities with testable interface
    - **env/**: Environment variable loading
    - **idgen/**: ID generation (ULID-based)
    - **ptr/**: Pointer utilities

## Core Commands

```bash
# Generate OpenAPI code
go generate ./...

# Build the API server
go build -o /dev/null ./cmd/api

# Run all tests with shuffling
go test -shuffle=on ./...

# Run tests with coverage
go test -cover ./...

# Format code
go tool goimports -w .

# Run linters
go vet ./...
go tool staticcheck ./...
```

## OpenAPI Integration

The project uses ogen for OpenAPI code generation:
- **API specification**: `doc/openapi.yaml`
- **Generated code**: `internal/api/openapi/`
- **Code generation directive**: In `internal/api/handler/handler.go`
- **Configuration**: `.ogen.yaml` specifies generator options
- Generate with: `go generate ./...`

## Testing Architecture

- **Testing framework**: testify for assertions and test utilities
- **Integration testing**: testcontainers for MySQL integration tests
- **HTTP testing**: ikawaha/httpcheck for API endpoint testing
- **Test utilities**:
  - `internal/lib/clock`: Testable time interface
  - `internal/lib/idgen`: Testable ID generation
  - `internal/database/databasetest`: Database test helpers
- Run tests with: `go test -shuffle=on ./...`

## Key Dependencies

- **Go version**: 1.25.0
- **Web framework**: ogen-generated server code
- **Database**: GORM v1.31+ with MySQL driver
- **Authentication**: JWT via golang-jwt/jwt v5
- **Logging**: slog (standard library) with custom wrapper (alog)
- **Testing**: testify, testcontainers, httpcheck
- **ID generation**: ULID via oklog/ulid with custom wrapper (idgen)

## Development Notes

- **Timezone**: Set to JST (Japan Standard Time) in `cmd/api/main.go`
- **ID generation**: Uses ULID for all entity IDs
- **Logging**: Structured JSON logging via `alog` package (slog wrapper)
  - Configurable via `LOG_LEVEL` (debug/info/silent) and `LOG_INDENT` env vars
- **Middleware**: Custom middleware for access logging, recovery, and request ID context
- **Database**: GORM conventions with MySQL
- **Authentication**: JWT-based with custom security handler
- **Error handling**: Centralized error types in `internal/api/usecase/error*.go`
  - System errors: ValidationError, AuthorizationError, etc.
  - User-facing errors: DomainValidationError, ProjectNotFoundError, etc.

## Code Style

- No unnecessary code comments (self-documenting code preferred)
- Use generics where appropriate (e.g., `ptr.Ref[T]`, `ternary[T]`)
- Private helper functions for type conversions
- Consistent naming: `convert*` for type conversions, `validate*` for validation
- Domain-specific utilities over generic ones
