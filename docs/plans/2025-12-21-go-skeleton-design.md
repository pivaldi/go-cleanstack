# Go Skeleton - Design Document

## Overview

A GitHub template repository that bootstraps production-ready Go projects with:
- CLI interface (Cobra) and HTTP API (Connect RPC) in one binary
- Simplified DDD with package-based layered architecture
- Database-agnostic persistence layer using sqlx
- Full test pyramid (unit, integration, e2e)
- Code quality via golangci-lint with pre-commit hooks
- Configuration via Viper with direnv + environment-specific TOML files
- Docker + docker-compose for local development

**Key Principle:** Provide enough structure to guide best practices, but stay minimal enough that users delete little when starting a real project.

---

## Technology Stack

| Component | Choice |
|-----------|--------|
| CLI | Cobra |
| API | Connect RPC (buf generated) |
| Architecture | Package-based DDD layers |
| Database | sqlx (database-agnostic) |
| Migrations | golang-migrate |
| Config | Viper + direnv + TOML per environment |
| Tests | gotestsum, testify, testcontainers |
| Linting | golangci-lint + pre-commit hooks |
| Task runner | just |
| Docker | Multi-stage build + docker-compose |

---

## Directory Structure

```
go-skeleton/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   └── repository/
│   ├── app/
│   │   └── service/
│   └── infra/
│       ├── api/
│       │   ├── proto/
│       │   │   └── myapp/v1/
│       │   ├── gen/
│       │   ├── buf.yaml
│       │   ├── buf.gen.yaml
│       │   ├── server.go
│       │   └── interceptor/
│       ├── cli/
│       ├── persistence/
│       │   └── migrations/
│       └── config/
├── tests/
│   ├── integration/
│   ├── e2e/
│   └── testutil/
├── .envrc.example
├── config_development.toml
├── config_staging.toml.example
├── config_production.toml.example
├── .golangci.yml
├── .pre-commit-config.yaml
├── Dockerfile
├── docker-compose.yml
├── justfile
├── .gitignore
└── README.md
```

---

## Layered Architecture

### Layer Rules

- **domain/** - Depends on nothing. Pure Go, no imports from other layers.
- **app/** - Depends only on domain/. Contains use cases and service orchestration.
- **infra/** - Depends on domain/ and app/. Implements interfaces defined in domain.

### Domain Layer (`internal/domain/`)

```go
// entity/item.go
type Item struct {
    ID          string
    Name        string
    Description string
    CreatedAt   time.Time
}

// repository/item.go
type ItemRepository interface {
    Create(ctx context.Context, item *entity.Item) error
    GetByID(ctx context.Context, id string) (*entity.Item, error)
    List(ctx context.Context) ([]*entity.Item, error)
    Delete(ctx context.Context, id string) error
}
```

### Application Layer (`internal/app/service/`)

```go
// item_service.go - implements the Connect RPC interface
type ItemService struct {
    repo repository.ItemRepository
}

func (s *ItemService) CreateItem(ctx, req) (*resp, error) {
    // Orchestrates domain logic, calls repo
}
```

### Infrastructure Layer (`internal/infra/`)

- `api/` - Connect RPC server, generated handlers, interceptors
- `cli/` - Cobra commands
- `persistence/` - sqlx repository implementations + migrations
- `config/` - Viper configuration loading

---

## CLI Structure

```
app serve              # Start HTTP server
app migrate up         # Run database migrations
app migrate down       # Rollback migrations
app item list          # List items (example CRUD via CLI)
app item create        # Create item
app version            # Print version info
```

**File structure:**

```
internal/infra/cli/
├── root.go            # Root command, loads config via Viper
├── serve.go           # Start HTTP server
├── migrate.go         # Migration commands
├── item.go            # Item CRUD commands (example)
└── version.go         # Version command
```

---

## API Layer (Connect RPC)

Protobuf definitions generate type-safe handlers via buf.

**buf.gen.yaml:**

```yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/user/go-skeleton/internal/infra/api/gen
plugins:
  - remote: buf.build/protocolbuffers/go
    out: gen
    opt: paths=source_relative
  - remote: buf.build/connectrpc/go
    out: gen
    opt: paths=source_relative
inputs:
  - directory: proto
```

The `ItemService` in `internal/app/service/` implements the generated Connect interface.

---

## Configuration Strategy

**Environment selection via direnv:**

```bash
# .envrc.example
export APP_ENV=development
```

**Environment-specific TOML files:**

| File | Committed | Purpose |
|------|-----------|---------|
| config_development.toml | Yes | Safe dev defaults |
| config_staging.toml | No | Gitignored, sensitive |
| config_production.toml | No | Gitignored, sensitive |
| config_staging.toml.example | Yes | Template |
| config_production.toml.example | Yes | Template |

**config_development.toml:**

```toml
[server]
port = 8080

[database]
url = "postgres://user:pass@localhost:5432/app?sslmode=disable"

[log]
level = "debug"
```

**Loading logic:**

```go
env := os.Getenv("APP_ENV")
if env == "" {
    env = "development"
}
viper.SetConfigName("config_" + env)
viper.SetConfigType("toml")
viper.AddConfigPath(".")
```

---

## Test Strategy

### Unit Tests (in-package `_test.go` files)

- `internal/domain/entity/` - Entity validation, business rules
- `internal/app/service/` - Service logic with mocked repositories
- Use standard `testing` + `testify` for assertions/mocks

### Integration Tests (`tests/integration/`)

- Repository implementations against real database
- Use testcontainers-go to spin up database
- Tests run with `go test -tags=integration`

### End-to-End Tests (`tests/e2e/`)

- Full API tests via HTTP client against running server
- Use testcontainers for both app + database
- Tests run with `go test -tags=e2e`

### Test Helpers

- `tests/testutil/` - Shared fixtures, database setup, test containers config

---

## Linting & Pre-commit

**.golangci.yml:**

```yaml
run:
  timeout: 5m

linters:
  enable:
    - errcheck
    - govet
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - gofmt
    - goimports
    - misspell
    - revive
    - gocritic

linters-settings:
  revive:
    rules:
      - name: exported
        disabled: true
```

**.pre-commit-config.yaml:**

```yaml
repos:
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.61.0
    hooks:
      - id: golangci-lint
  - repo: https://github.com/pre-commit/pre-commit-hooks
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
  - repo: local
    hooks:
      - id: buf-lint
        name: buf lint
        entry: buf lint
        language: system
        files: '\.proto$'
```

---

## Docker

**Dockerfile:**

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/app ./cmd

# Runtime stage
FROM alpine:latest
COPY --from=builder /app/bin/app /usr/local/bin/app
COPY internal/infra/persistence/migrations /migrations
ENTRYPOINT ["app"]
CMD ["serve"]
```

**docker-compose.yml:**

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: app
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

---

## Task Runner (justfile)

```just
# Development
dev:
    go run ./cmd serve

# Code generation
generate:
    buf generate --path internal/infra/api

# Testing
test:
    gotestsum -- ./...

test-int:
    gotestsum -- -tags=integration ./tests/integration/...

test-e2e:
    gotestsum -- -tags=e2e ./tests/e2e/...

test-all:
    gotestsum -- -tags=integration,e2e ./...

# Database
migrate-up:
    go run ./cmd migrate up

migrate-down:
    go run ./cmd migrate down

# Linting
lint:
    golangci-lint run

# Docker
up:
    docker-compose up -d

down:
    docker-compose down
```

---

## Wiring (cmd/main.go)

```
config → db connection → repo → service → Connect handler → HTTP server
                                       ↘ Cobra commands
```

---

## .gitignore

```
# Binaries
/bin/
*.exe

# Environment
.envrc
config_staging.toml
config_production.toml

# Generated
internal/infra/api/gen/

# IDE
.idea/
.vscode/

# Test
coverage.out
```

---

## Next Steps

1. Create isolated git worktree for implementation
2. Create detailed implementation plan
3. Implement in incremental PRs
