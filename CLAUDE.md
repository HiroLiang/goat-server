# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Goat-Server is the Go backend for Goat-Chat, a Tauri desktop application with Ollama agent integration. It provides REST APIs and WebSocket endpoints for user authentication, chat management, and AI agent interaction.

## Commands

```shell
# Install CLI tools (swag, godog)
make setup

# Run locally
make run

# Build binary to bin/goat-api
make build

# Run all tests (unit + BDD)
make test

# Run unit tests only
go test ./internal/... -v

# Run BDD tests only
go test -v ./features

# Run a single unit test
go test ./internal/application/user/... -run TestAssignRoleToUser_Success -v

# Regenerate Swagger docs (required after modifying handler annotations)
make swag
```

## Architecture

The project follows Clean Architecture with four main layers:

```
internal/
├── domain/        # Entities, value objects, repository interfaces, domain errors
├── application/   # Use cases (orchestrate domain objects)
├── infrastructure/ # Concrete implementations (Postgres, Redis, JWT, Argon2)
└── interface/http/ # Gin handlers, middleware, DTOs
```

**Dependency direction:** `interface → application → domain ← infrastructure`

### Bootstrap Flow (`internal/bootstrap/`)

`main.go` → `bootstrap.CreateApp()` → `Start()`:
1. `BuildDeps()` — wires repositories and services from infrastructure implementations
2. `BuildUseCases()` — constructs use cases with dependencies
3. `NewServer()` + `RegisterRestRoutes()` — sets up Gin with middleware and routes

### Request Flow

HTTP request → Gin middleware chain → Handler → `adapter.BuildInput(c, data)` → UseCase → Domain → Repository

`adapter.BuildInput` packages the Gin context (IP, auth token, user ID) into `shared.UseCaseInput[T]`, which is the standard input type for all use cases.

### Authentication

Two-layer middleware pattern:
- `AuthMiddleware` — validates Bearer token against Redis session store, sets `authContext` in Gin context (non-blocking, continues even without valid token)
- `RequireAuthMiddleware` — aborts with 401 if `authContext` is absent

Tokens are JWT-based but validated against Redis session store (server-side sessions). Login returns the token in the `Authorization` response header.

### Error Handling Pattern

Domain errors are defined per-domain (e.g., `internal/domain/user/errors.go`). Each handler has a dedicated `*_error_translator.go` that translates domain errors to HTTP status codes using `errors.Is`. Unhandled errors are passed to the global error middleware via `c.Error(err)`.

### Testing

- **Unit tests**: Live alongside source files (e.g., `user_usecase_test.go`). Use `testify/mock` for mocking repositories.
- **BDD tests** (`features/`): Use Godog with Gherkin `.feature` files. The suite boots a real `httptest.Server` with mock dependencies loaded from `dev-doc/config/config.yaml`.

Mock implementations are in:
- `internal/infrastructure/persistence/mock/` — repository mocks
- `internal/infrastructure/auth/mock/` — auth mocks
- `internal/infrastructure/shared/mock/` — rate limiter mock

### Configuration

- `config/config.yaml` — main config (auth token expiration, DB DSN, Redis, rate limits, HMAC secret)
- `.env` — `APP_ENV` and `SERVER_PORT`
- BDD tests use `dev-doc/config/config.yaml` instead of the main config

### Infrastructure

- **Database**: PostgreSQL via `pgx` + `sqlx`. Query building uses `Masterminds/squirrel`. Each entity has a `*_record.go` (DB struct), `*_mapper.go` (record ↔ domain), and `*_repo.go`.
- **Cache**: Redis via `go-redis`. UserRole repo uses a Redis cache layer wrapping Postgres.
- **Logger**: `go.uber.org/zap` via `internal/logger/logger.go`.

### Adding a New Domain Feature

1. Define entity/VOs/errors/repository interface in `internal/domain/<name>/`
2. Add use case input/output types and `UseCase` struct in `internal/application/<name>/`
3. Implement the repository in `internal/infrastructure/persistence/postgres/<name>/`
4. Register the repository and use case in `internal/bootstrap/dependencies.go` and `usecases.go`
5. Add handler + DTO + error translator in `internal/interface/http/handler/<name>/`
6. Register routes in `internal/bootstrap/rest.go`
