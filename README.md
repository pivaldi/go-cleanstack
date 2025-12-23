# GoCleanstack

## THIS IS A WORK IN PROGRESS !
## DO NOT USE IN PRODUCTION !

A production-ready Go application skeleton demonstrating Clean Architecture principles with Domain-Driven Design (DDD) and Hexagonal Architecture (Ports & Adapters).

## Features

- **Clean Architecture**: Strict separation of concerns with domain, application, and infrastructure layers
- **Hexagonal Architecture**: Ports and adapters pattern for dependency inversion
- **Connect RPC**: Modern RPC framework built on protobuf and HTTP/2
- **CLI Interface**: Cobra-based command-line interface
- **Database**: PostgreSQL with sqlx and golang-migrate
- **Configuration**: Viper-based configuration with environment-specific TOML files
- **Testing**: Full test pyramid (unit, integration, e2e) with testcontainers
- **Docker**: Production-ready Docker setup with docker-compose
- **Dev Tools**: justfile, golangci-lint, pre-commit hooks

## Architecture

This project follows Hexagonal Architecture principles with clear boundaries:

```
domain/     - Pure business logic and entities (no external dependencies)
  ├── entity/    - Domain entities with business rules
  └── ports/     - Interfaces defining what domain needs

app/        - Application orchestration layer
  ├── service/   - Use cases and business workflows
  └── adapters/  - Bridges between domain and infrastructure

infra/      - External concerns and frameworks
  ├── api/       - Connect RPC handlers, protobuf, and HTTP server
  ├── config/    - Configuration loading (Viper)
  └── persistence/ - Database access with DTOs and migrations

platform/   - Cross-cutting infrastructure utilities
  ├── logging/   - Zap logger factory
  ├── apperr/    - Application error types
  ├── clierr/    - CLI error handling
  └── reqid/     - Request ID utilities

transport/  - Transport layer utilities
  └── connectx/  - Connect RPC interceptors and error handling
```

### Dependency Rules

- **Domain** has no dependencies on any other layer
- **Application** depends only on domain (via ports)
- **Infrastructure** provides implementations of domain ports via adapters
- Infrastructure never imports domain directly (uses DTOs instead)

## Project Structure

```
.
├── main.go                      # Application entry point
├── cmd/                         # CLI commands (Cobra)
│   ├── root.go                 # Root command with logger initialization
│   ├── serve.go                # HTTP server command
│   ├── migrate.go              # Database migration commands
│   └── version.go              # Version command
├── internal/
│   ├── domain/                  # Core business logic
│   │   ├── entity/             # Domain entities
│   │   │   ├── item.go
│   │   │   └── item_test.go
│   │   └── ports/              # Port interfaces
│   │       └── repository.go
│   ├── app/                     # Application layer
│   │   ├── service/            # Use cases
│   │   │   ├── item_service.go
│   │   │   └── item_service_test.go
│   │   └── adapters/           # Domain-to-infra adapters
│   │       └── item_repo_adapter.go
│   ├── infra/                   # Infrastructure layer
│   │   ├── api/                # Connect RPC API
│   │   │   ├── proto/cleanstack/v1/
│   │   │   │   └── item.proto
│   │   │   ├── gen/            # Generated protobuf code
│   │   │   │   └── cleanstack/v1/
│   │   │   ├── handler/
│   │   │   │   └── item_handler.go
│   │   │   ├── interceptor/    # HTTP interceptors
│   │   │   └── server.go
│   │   ├── config/             # Configuration
│   │   │   ├── config.go
│   │   │   └── config_test.go
│   │   └── persistence/        # Database access
│   │       ├── db.go
│   │       ├── models.go       # DTOs
│   │       ├── item_repo.go
│   │       └── migrations/
│   │           ├── 000001_create_items_table.up.sql
│   │           └── 000001_create_items_table.down.sql
│   ├── platform/                # Cross-cutting utilities
│   │   ├── logging/            # Zap logger factory
│   │   │   └── zap.go
│   │   ├── apperr/             # Application errors
│   │   ├── clierr/             # CLI error handling
│   │   └── reqid/              # Request ID utilities
│   └── transport/               # Transport layer
│       └── connectx/           # Connect RPC utilities
│           ├── interceptors.go
│           └── connect_errors.go
├── tests/
│   ├── integration/             # Integration tests
│   │   └── item_repo_test.go
│   ├── e2e/                     # End-to-end tests
│   │   └── api_test.go
│   └── testutil/                # Test utilities
│       ├── containers.go
│       └── db.go
├── config_default.toml          # Default configuration
├── docker-compose.yml           # Docker orchestration
├── Dockerfile                   # Production container
├── justfile                     # Task runner
└── .golangci.yml                # Linter configuration
```

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 16
- Docker and docker-compose (for testing and local development)
- [buf](https://buf.build) (for protobuf generation)
- [just](https://github.com/casey/just) (task runner)
- [golangci-lint](https://golangci-lint.run) (optional, for linting)

## Getting Started

### 1. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 2. Configuration

The application uses environment-based configuration with two files:

1. `config_default.toml` - Base configuration (always loaded)
2. `config_<env>.toml` - Environment-specific overrides (merged with default)

Create environment-specific configuration:

```bash
cp config_default.toml config_development.toml
# Edit config_development.toml with development settings

cp config_default.toml config_production.toml
# Edit config_production.toml with production settings
```

**Important:** Set `APP_ENV` environment variable to choose configuration (REQUIRED - no default):

```bash
export APP_ENV=development  # or production
```

### 3. Start Database

Using Docker:

```bash
just up
```

Or manually start PostgreSQL and update the database URL in your config file.

### 4. Run Migrations

```bash
just migrate-up
```

### 5. Run the Server

```bash
export APP_ENV=development
just dev
```

Or run directly:

```bash
APP_ENV=development go run main.go serve
```

The API will be available at `http://localhost:8080`.

## Development Commands

This project uses [just](https://github.com/casey/just) as a task runner. Available commands:

### Development
```bash
just dev          # Run development server (requires APP_ENV env var)
just generate-api # Generate code from protobuf definitions
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
- `internal/domain/entity/item_test.go` - Domain entity tests
- `internal/app/service/item_service_test.go` - Service layer tests
- `internal/infra/config/config_test.go` - Configuration tests

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

Protocol Buffer definitions are in `internal/infra/api/proto/cleanstack/v1/`. To regenerate code:

```bash
just generate-api
```

This uses [buf](https://buf.build) to generate:
- Go structs from protobuf messages
- Connect RPC service interfaces and handlers
- Generated code is output to `internal/infra/api/gen/`

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
2. Update `go.mod` with your module name
3. Update protobuf package names in `internal/infra/api/proto/`
4. Implement your domain entities and business logic
5. Add corresponding service methods and API endpoints

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
