# GoCleanstack

## THIS IS A WORK IN PROGRESS !
## DO NOT USE IN PRODUCTION !

A production-ready **Go multi-module workspace** demonstrating Clean Architecture principles with Domain-Driven Design and Hexagonal Architecture, governed by a unified CLI project manager.

## Go Workspace Architecture

This project uses **Go workspaces** (`go.work`) to manage multiple independent modules:

| Module | Path | Purpose |
|--------|------|---------|
| Root | `/` | CLI orchestrator, aggregates sub-applications |
| Common | `/internal/common` | Shared platform utilities (logging, config, errors) |
| App1 | `/internal/app/app1` | First application with its own domain and infrastructure |

**Benefits:**
- Independent versioning and testing per module
- Clear boundaries between applications
- Shared utilities without code duplication
- Easy to add new applications (app2, app3, etc.)

## Features

- **Clean Architecture**: Strict separation of concerns with domain, application, and infrastructure layers
- **Hexagonal Architecture**: Ports and adapters pattern for dependency inversion
- **Connect RPC**: Modern RPC framework built on protobuf and HTTP/2
- **CLI Interface**: Cobra-based command-line interface
- **Database**: PostgreSQL with sqlx and Goose v3 migrations
- **Configuration**: Viper-based configuration with environment-specific TOML files
- **Testing**: Full test pyramid (unit, integration, e2e) with testcontainers
- **Docker**: Production-ready Docker setup with docker-compose
- **Dev Tools**: justfile, golangci-lint, pre-commit hooks

## Architecture

This project follows Hexagonal Architecture principles with clear boundaries, organized as a Go workspace:

```
/                           - Root module (CLI orchestrator)
├── main.go                 - Entry point, aggregates sub-applications
├── cmd/                    - Root CLI commands (migrate)

/internal/common/           - Common module (shared utilities)
├── platform/
│   ├── logging/           - Logger abstraction (zap-based)
│   ├── config/            - Generic configuration loader
│   ├── apperr/            - Application error types
│   ├── clierr/            - CLI error handling
│   └── reqid/             - Request ID utilities
└── transport/
    └── connectx/          - Connect RPC interceptors and error handling

/internal/app/app1/         - App1 module (first application)
├── domain/                - Pure business logic (no external dependencies)
│   ├── entity/            - Domain entities with business rules
│   └── ports/             - Interfaces defining what domain needs
├── service/               - Use cases and business workflows
├── adapters/              - Bridges between domain and infrastructure
├── config/                - App-specific configuration
├── cmd/                   - App CLI commands (serve, version)
└── infra/                 - External concerns and frameworks
    ├── api/               - Connect RPC handlers, protobuf, HTTP server
    └── persistence/       - Database access with DTOs and migrations
```

### Dependency Rules

- **Domain** has no dependencies on any other layer
- **Application** depends only on domain (via ports)
- **Infrastructure** provides implementations of domain ports via adapters
- Infrastructure never imports domain directly (uses DTOs instead)

## Project Structure

```
.
├── go.work                       # Go workspace definition
├── go.mod                        # Root module
├── main.go                       # Entry point (CLI orchestrator)
├── cmd/                          # Root CLI commands
│   └── migrate.go               # Database migration commands
│
├── internal/
│   ├── common/                   # Common module (go.mod)
│   │   ├── platform/            # Cross-cutting utilities
│   │   │   ├── logging/         # Logger abstraction
│   │   │   ├── config/          # Generic config loader
│   │   │   ├── apperr/          # Application errors
│   │   │   ├── clierr/          # CLI error handling
│   │   │   └── reqid/           # Request ID utilities
│   │   └── transport/           # Transport utilities
│   │       └── connectx/        # Connect RPC interceptors
│   │
│   └── app/
│       └── app1/                 # App1 module (go.mod)
│           ├── main.go          # App entry point (standalone use)
│           ├── cmd/             # App CLI commands
│           │   ├── root.go      # Root command with logger init
│           │   ├── serve.go     # HTTP server command
│           │   └── version.go   # Version command
│           ├── config/          # App-specific configuration
│           │   └── config.go
│           ├── domain/          # Core business logic
│           │   ├── entity/      # Domain entities
│           │   │   └── item.go
│           │   └── ports/       # Port interfaces
│           │       └── repository.go
│           ├── service/         # Use cases
│           │   └── item_service.go
│           ├── adapters/        # Domain-to-infra adapters
│           │   └── item_repo_adapter.go
│           └── infra/           # Infrastructure layer
│               ├── api/         # Connect RPC API
│               │   ├── proto/   # Protobuf definitions
│               │   ├── gen/     # Generated code
│               │   ├── handler/ # Request handlers
│               │   └── server.go
│               └── persistence/ # Database access
│                   ├── db.go
│                   ├── models.go
│                   ├── item_repo.go
│                   └── migrations/
│
├── tests/
│   ├── integration/              # Integration tests
│   ├── e2e/                      # End-to-end tests
│   └── testutil/                 # Test utilities
│
├── config_default.toml           # Default configuration
├── docker-compose.yml            # Docker orchestration
├── Dockerfile                    # Production container
├── justfile                      # Task runner
└── .golangci.yml                 # Linter configuration
```

## Prerequisites

- Go 1.25 or higher (optional if the app is built from Docker)
- PostgreSQL 16 or higher
- Docker and docker-compose (for testing and **local development**)
- Node.js/npm (for **local development**)
- Python 3 (for pre-commit hooks and **local development**)

**If the installation of Python 3 is system-wide, you must install [pre-commit](https://pre-commit.com/#install)
yourself.**

## Getting Started

### 1. Install Dependencies

Simply run this command :

```bash
./configure
```

This command can be launch as many times as you want.

### 2. Configuration

The application uses environment-based configuration with two files:

1. `config_default.toml` - Base configuration (always loaded)
2. `config_<env>.toml` - Environment-specific overrides (merged with default)

**Important:** Set `APP_ENV` environment variable to choose configuration (REQUIRED - no default):

```bash
export APP_ENV=development  # or production
```

Or use `[direnv](https://direnv.net/)`.

### 3. Start Database

Using Docker:

`just up` or `docker-compose up -d` or manually start PostgreSQL and update the database URL in your config file.

### 4. Run Migrations

`just migrate-up` or `./bin/cleanstack migrate up` or `go run . migrate up`

### 5. Run the Server

`just dev` or `go run . serve`

The API will be available at `http://localhost:4224`.

## Development Commands

This project uses [just](https://github.com/casey/just) as a task runner.
Available commands follow.

### Development
```bash
just generate-api # Generate code from protobuf definitions
just dev          # Run development server (requires APP_ENV env var)
```

### Testing
```bash
just test         # Run unit tests
just test-int     # Run integration tests (requires Docker)
just test-e2e     # Run end-to-end tests (requires Docker)
just test-all     # Run all tests
just test-cover   # Generate test coverage report
```

### Database
```bash
just migrate-up   # Run database migrations
just migrate-down # Rollback database migrations
```

### Code Quality
```bash
just lint         # Run linter
just lint-fix     # Run linter with auto-fix
```

### Build
```bash
just build        # Build binary to bin/cleanstack
```

### Docker
```bash
just up           # Start services with docker-compose
just down         # Stop services
just logs         # View service logs
```

### Cleanup
```bash
just clean        # Remove build artifacts and coverage files
```

## Testing

This project implements a comprehensive test pyramid:

### Unit Tests

Run unit tests that don't require external dependencies:

```bash
just test
# or
go test ./...
```

Unit tests are co-located with the code they test:
- `internal/app/app1/domain/entity/item_test.go` - Domain entity tests
- `internal/app/app1/service/item_service_test.go` - Service layer tests
- `internal/common/platform/config/config_test.go` - Configuration tests

### Integration Tests

Run integration tests that use testcontainers for real PostgreSQL:

```bash
just test-int
# or
go test -tags=integration ./tests/integration/...
```

Integration tests verify:
- Database operations with real PostgreSQL
- Repository implementations
- Data persistence and retrieval

### E2E Tests

Run end-to-end tests that test the full application stack:

```bash
just test-e2e
# or
go test -tags=e2e ./tests/e2e/...
```

E2E tests verify:
- Complete API flows
- Request/response handling
- Full dependency wiring

### Test Coverage

Generate a test coverage report:

```bash
just test-cover
open coverage.html
```

## API Usage

The application exposes a Connect RPC API. You can interact with it using any HTTP client or Connect-compatible client.

### Example: Create an Item

```bash
curl -X POST http://localhost:8080/cleanstack.v1.ItemService/CreateItem \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Example Item",
    "description": "An example item"
  }'
```

### Example: List Items

```bash
curl -X POST http://localhost:8080/cleanstack.v1.ItemService/ListItems \
  -H "Content-Type: application/json" \
  -d '{}'
```

## Docker Deployment

### Build and Run

```bash
# Build the image
docker build -t go-cleanstack .

# Run with docker-compose
docker-compose up
```

The docker-compose setup includes:
- Application container (port 8080)
- PostgreSQL database (port 5432)
- Health checks and automatic migrations

### Environment Variables

Configure the application using environment variables in docker-compose.yml:

```yaml
environment:
  - APP_ENV=production
```

## Code Generation

Protocol Buffer definitions are in `internal/app/app1/infra/api/proto/cleanstack/v1/`. To regenerate code:

```bash
just generate-api
```

This uses [buf](https://buf.build) to generate:
- Go structs from protobuf messages
- Connect RPC service interfaces and handlers
- Generated code is output to `internal/app/app1/infra/api/gen/`

## Linting

This project uses golangci-lint with a comprehensive configuration:

```bash
just lint        # Check for issues
just lint-fix    # Auto-fix issues
```

Pre-commit hooks are configured to run linting automatically.

## Contributing

This is a template/skeleton project demonstrating architecture patterns. To use it as a starting point:

1. Fork or clone the repository
2. Update module names in `go.mod`, `go.work`, and all sub-module `go.mod` files
3. Update `replace` directives in all `go.mod` files to match your module path
4. Update protobuf package names in `internal/app/app1/infra/api/proto/`
5. Implement your domain entities and business logic
6. Add corresponding service methods and API endpoints
7. To add a new application, create a new module under `internal/app/` (see CLAUDE.md for detailed steps)

## Architecture Highlights

### Hexagonal Architecture Benefits

1. **Testability**: Domain logic can be tested without databases or HTTP
2. **Flexibility**: Easy to swap infrastructure implementations
3. **Independence**: Domain layer is completely isolated from frameworks
4. **Maintainability**: Clear boundaries reduce cognitive load

### The Adapter Pattern

The `ItemRepositoryAdapter` demonstrates the adapter pattern:

```go
// Domain defines what it needs (port)
type ItemRepository interface {
    Create(ctx context.Context, item *entity.Item) error
}

// Infrastructure provides DTOs (no domain import)
type ItemRepo struct {}
func (r *ItemRepo) Create(ctx context.Context, dto *ItemDTO) error

// Adapter bridges domain and infrastructure
type ItemRepositoryAdapter struct {
    infraRepo *persistence.ItemRepo
}
func (a *ItemRepositoryAdapter) Create(ctx context.Context, item *entity.Item) error {
    dto := &persistence.ItemDTO{...} // Convert entity to DTO
    return a.infraRepo.Create(ctx, dto)
}
```

This ensures infrastructure never imports domain, maintaining clean boundaries.

## Resources

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Connect RPC](https://connectrpc.com/)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
