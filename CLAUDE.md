# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a production-ready Go application skeleton demonstrating **Clean Architecture** with **Hexagonal Architecture** (Ports & Adapters pattern) and **Domain-Driven Design**. The project uses Connect RPC (built on protobuf and HTTP/2) for the API layer.

## Essential Commands

### Development
```bash
just dev              # Run development server (requires APP_ENV=development)
go run . serve    # Alternative way to run server
```

### Testing
```bash
just test             # Unit tests only
just test-int         # Integration tests (requires Docker)
just test-e2e         # End-to-end tests (requires Docker)
just test-all         # All tests including integration and e2e
just test-cover       # Generate coverage report (creates coverage.html)

# Run specific tests
go test ./internal/app/service/...
go test -run TestItemService_CreateItem ./internal/app/service/...
```

### Code Quality
```bash
just lint             # Run golangci-lint (requires Go 1.24+)
just lint-fix         # Auto-fix linting issues
```

### Database
```bash
just migrate-up       # Run migrations
just migrate-down     # Rollback last migration
go run . migrate up     # Alternative way
```

### Code Generation
```bash
just generate-api     # Regenerate Connect RPC code from protobuf
```

### Build & Docker
```bash
just build            # Build binary to bin/cleanstack
just up               # Start docker-compose for current APP_ENV (reads from .envrc)
just down             # Stop docker-compose
just logs             # View logs from all running services
```

**Port allocation** (set automatically by APP_ENV):
- development: App runs locally on 4224, DB in Docker on 5435
- staging: Full stack in Docker - App 4225, DB 5436
- production: Full stack in Docker - App 4226, DB 5437

## Multi-Environment Workflow

This project uses `APP_ENV` to drive configuration across all tools (config files, docker-compose, justfile).

### Architecture
- **Development**: Hybrid - App runs locally (`just dev`), database in Docker
  - Allows hot reload and debugging
  - Database URL: `localhost:5435`
- **Staging/Production**: Full stack in Docker
  - Both app and database run in containers
  - Database URL: Docker network using container name

### Setup for Multiple Environments
1. Clone/checkout repo to separate directories for each environment:
```bash
# Example directory structure
~/projects/cleanstack-development/
~/projects/cleanstack-staging/
~/projects/cleanstack-production/
```

2. In each directory, run `./configure` and select the appropriate environment
3. This creates `.envrc` with `export APP_ENV=<environment>` and `config_<environment>.toml`

### Usage
```bash
# Development workflow
source .envrc          # Load APP_ENV=development
just up                # Start database in Docker on port 5435
just dev               # Run app locally on port 4224

# Staging/Production workflow
source .envrc          # Load APP_ENV=staging or production
just up                # Start both app + database in Docker
```

### How APP_ENV Works
- **`.envrc`**: Sets APP_ENV environment variable
- **`config.go`**: Loads `config_default.toml` + `config_${APP_ENV}.toml`
- **`docker-compose.yml`**: Uses `${APP_ENV}` for service names, ports, database names
- **`justfile`**: Validates APP_ENV and sets ports before starting docker-compose
- **`./configure`**: Creates environment-specific config files from `.example` templates

### Database Access
```bash
# Each environment on different port
psql -h localhost -p 5435 -U user -d cleanstack_development  # Development
psql -h localhost -p 5436 -U user -d cleanstack_staging      # Staging
psql -h localhost -p 5437 -U user -d cleanstack_production   # Production
```

### Running Multiple Environments Simultaneously
Since each environment is in its own directory with its own `APP_ENV`, you can run all three at once without conflicts:
```bash
# Terminal 1: Development
cd ~/projects/cleanstack-development && source .envrc && just up && just dev

# Terminal 2: Staging
cd ~/projects/cleanstack-staging && source .envrc && just up

# Terminal 3: Production
cd ~/projects/cleanstack-production && source .envrc && just up
```

Each environment runs on its own ports (4224/5435, 4225/5436, 4226/5437) without conflicts.

## Architecture: Hexagonal (Ports & Adapters)

This codebase enforces strict architectural boundaries. Understanding the dependency flow is critical:

```
┌──────────────────────────────────────────────────────────┐
│ Infrastructure Layer (., internal/infra/)             │
│ - API handlers, HTTP server, database, config, CLI      │
│ - Depends on: Application layer (via service imports)   │
│ - Implements: Domain ports (via adapters)               │
└────────────────────┬─────────────────────────────────────┘
                     │
                     ↓
┌──────────────────────────────────────────────────────────┐
│ Application Layer (internal/app/)                        │
│ - Services (use cases, business workflows)               │
│ - Adapters (bridge domain ↔ infrastructure)             │
│ - Depends on: Domain layer (ports only, never entities) │
└────────────────────┬─────────────────────────────────────┘
                     │
                     ↓
┌──────────────────────────────────────────────────────────┐
│ Domain Layer (internal/domain/)                          │
│ - Entities (pure business logic)                         │
│ - Ports (interfaces defining infrastructure needs)       │
│ - Depends on: NOTHING (zero external dependencies)       │
└──────────────────────────────────────────────────────────┘
```

### Critical Architectural Rules

1. **Domain is pure**: No imports of app, infra, or external packages (only standard library)
2. **Infrastructure never imports domain entities**: Use DTOs instead
3. **Adapters bridge the gap**: Convert between domain entities and infrastructure DTOs
4. **Logging is infrastructure**: Domain has NO logger, app/infra layers receive logger via dependency injection

### The DTO/Entity Separation Pattern

This codebase uses a strict separation to maintain clean boundaries:

**Infrastructure DTOs** (`internal/infra/persistence/models.go`):
```go
type ItemDTO struct {
    ID          string
    Name        string
    Description string
    CreatedAt   time.Time
}
```

**Domain Entities** (`internal/domain/entity/item.go`):
```go
type Item struct {
    ID          string
    Name        string
    Description string
    CreatedAt   time.Time
}
```

**Adapters** (`internal/app/adapters/item_repo_adapter.go`):
- Convert entity → DTO when calling infrastructure
- Convert DTO → entity when returning to domain
- Bridge the `ports.ItemRepository` interface to `persistence.ItemRepo`

### Dependency Wiring Pattern

The application wires dependencies in `cmd/serve.go`:

```go
// Infrastructure creates DB and concrete repo (uses DTOs)
db := persistence.NewDB(cfg.Database.URL)
infraRepo := persistence.NewItemRepo(db, logger)

// Adapter wraps infra repo to implement domain port
itemRepo := adapters.NewItemRepositoryAdapter(infraRepo)

// Service uses port interface (no knowledge of DTOs)
itemService := service.NewItemService(itemRepo, logger)

// API handler uses service
server := api.NewServer(cfg.Server.Port, itemService, logger)
```

## Configuration System

The application uses **environment-based configuration** with Viper:

### Configuration Files
- `config_default.toml` - Base configuration (always loaded)
- `config_<env>.toml` - Environment-specific overrides (merged with default)

### Configuration Loading
1. Set `APP_ENV` environment variable (REQUIRED: no default)
2. Application loads `config_default.toml` first
3. Then merges `config_<env>.toml` (e.g., `config_development.toml`, `config_production.toml`)
4. CLI flags override config file values

Example:
```bash
APP_ENV=development go run ./cmd serve
APP_ENV=production go run ./cmd serve --log-level debug
```

### Logger Configuration
- Logger is initialized in `cmd/root.go` via `PersistentPreRunE`
- Log level precedence: CLI flag `--log-level` → config file → default
- Development environment → console logger (human-readable)
- Production environment → JSON logger (structured)
- Logger passed explicitly to all app/infra constructors (no globals)

## Testing Strategy

### Test Organization
- **Unit tests**: Co-located with code (`*_test.go`), use mocks, no external dependencies
- **Integration tests**: In `tests/integration/`, use testcontainers for real PostgreSQL
- **E2E tests**: In `tests/e2e/`, test full stack with real dependencies

### Test File Requirements
All test files should:
- Use `zap.NewNop()` for logger in constructors
- Set `APP_ENV` environment variable in config tests: `os.Setenv("APP_ENV", "development")`
- Clean up with `t.Cleanup(func() { os.Unsetenv("APP_ENV") })`

### Common Test Patterns
```go
// Service layer unit test (uses mock repo)
mockRepo := new(MockItemRepository)
logger := zap.NewNop()
svc := service.NewItemService(mockRepo, logger)

// Integration test (uses testcontainers)
logger := zap.NewNop()
repo := persistence.NewItemRepo(db, logger)
```

## Linting Configuration

The project uses **golangci-lint v2.7+** with strict rules:

### Important Linter Rules
- **wrapcheck**: ALL errors from external packages must be wrapped with context
  ```go
  // ❌ Bad
  return err

  // ✅ Good
  return fmt.Errorf("failed to create item: %w", err)
  ```
- **nlreturn**: Blank line required before return statements after if/for blocks
- **lll**: Max line length 124 characters (split long function signatures)
- **revive**: Unused parameters must be renamed to `_`
- **gosec**: HTTP servers must have timeouts configured
- **mnd**: Magic numbers must be extracted to constants

### Go Version Requirement
- **Go 1.24+** required (specified in go.mod)
- golangci-lint must be built with Go 1.24+ to match project requirements

## Protocol Buffers & Code Generation

### Proto Files Location
`internal/infra/api/proto/cleanstack/v1/item.proto`

### Generating Code
```bash
just generate-api
# or
cd internal/infra/api && buf generate
```

This generates:
- Go structs from protobuf messages
- Connect RPC service interfaces in `internal/infra/api/gen/`

### After Regenerating
Update handlers in `internal/infra/api/handler/` to match new service interfaces.

## Common Development Patterns

### Adding a New Domain Entity
1. Create entity in `internal/domain/entity/` (pure logic, no dependencies)
2. Add validation methods to entity (e.g., `Validate()`)
3. Define repository port in `internal/domain/ports/`
4. Create DTO in `internal/infra/persistence/models.go`
5. Implement repository in `internal/infra/persistence/` (uses DTOs)
6. Create adapter in `internal/app/adapters/` (converts entity ↔ DTO)
7. Add service in `internal/app/service/` (uses port interface)
8. Create migration in `internal/infra/persistence/migrations/`
9. Add API handler in `internal/infra/api/handler/`
10. Update protobuf and regenerate code

### Adding a Logger to a Component
Loggers are passed via dependency injection:

```go
// Add logger field to struct
type MyService struct {
    logger *zap.Logger
}

// Add logger parameter to constructor
func NewMyService(logger *zap.Logger) *MyService {
    return &MyService{logger: logger}
}

// Use structured logging
s.logger.Info("operation started",
    zap.String("field", value),
    zap.Error(err),
)
```

**Never** create loggers inside components - always receive via constructor.

## File Naming Conventions

- Domain entities: `item.go`, `user.go` (singular noun)
- Repository implementations: `item_repo.go`, `user_repo.go`
- Adapters: `item_repo_adapter.go` (bridges domain ↔ infra)
- Services: `item_service.go`, `user_service.go`
- Handlers: `item_handler.go`, `user_handler.go`
- Tests: `*_test.go` (co-located with code)
- DTOs: Defined in `models.go` within infrastructure packages

## Critical Implementation Details

### Error Handling Philosophy
- Infrastructure errors must be wrapped with context at the boundary
- Service layer wraps repository errors with business context
- Adapters wrap infrastructure errors when crossing boundaries
- Use `fmt.Errorf("context: %w", err)` to maintain error chain

### Database Connections
- Connection pool configured in `internal/infra/persistence/db.go`
- Max open connections: 25 (constant `maxOpenConns`)
- Max idle connections: 5 (constant `maxIdleConns`)
- Migrations managed by golang-migrate

### HTTP Server Configuration
The server in `internal/infra/api/server.go` requires proper timeout configuration for security:
```go
httpServer := &http.Server{
    Addr:         addr,
    Handler:      handler,
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  30 * time.Second,
}
```

## Gotchas & Important Notes

1. **APP_ENV is required**: Application will not start without `APP_ENV` environment variable set
2. **Config files must exist**: Both `config_default.toml` and `config_<env>.toml` must be present
3. **DTO vs Entity**: Infrastructure NEVER imports domain entities - always use DTOs and adapters
4. **No global logger**: Logger must be passed to constructors, never created as global variable
5. **Integration tests need Docker**: Integration and E2E tests use testcontainers
6. **Protobuf changes require regeneration**: After modifying `.proto` files, run `just generate-api`
7. **CLI commands in cmd/ package**: CLI structure changed from `package main` to `package cmd` with commands in separate files
