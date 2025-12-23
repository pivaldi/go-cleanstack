# Go-Cleanstack Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a production-ready Go project skeleton with CLI, Connect RPC API, DDD layers, sqlx persistence, and full test pyramid.

**Architecture:** Package-based DDD with domain/app/infra layers. Connect RPC generates API handlers from protobuf. Cobra provides CLI. sqlx handles database access with golang-migrate for migrations.

**Tech Stack:** Go 1.23, Cobra, Connect RPC, buf, sqlx, golang-migrate, Viper, testify, testcontainers-go, gotestsum, golangci-lint, just, Docker

---

## Phase 1: Project Foundation

### Task 1: Initialize Go Module and Directory Structure

**Files:**
- Create: `go-cleanstack/go.mod`
- Create: `go-cleanstack/cmd/main.go`
- Create: `go-cleanstack/.gitignore`

**Step 1: Initialize Go module**

```bash
cd /workspace/go-cleanstack
go mod init github.com/user/go-cleanstack
```

**Step 2: Create directory structure**

```bash
mkdir -p cmd
mkdir -p internal/domain/entity
mkdir -p internal/domain/repository
mkdir -p internal/app/service
mkdir -p internal/infra/api/proto/cleanstack/v1
mkdir -p internal/infra/api/gen
mkdir -p internal/infra/api/interceptor
mkdir -p internal/infra/cli
mkdir -p internal/infra/persistence/migrations
mkdir -p internal/infra/config
mkdir -p tests/integration
mkdir -p tests/e2e
mkdir -p tests/testutil
```

**Step 3: Create minimal main.go**

```go
// cmd/main.go
package main

func main() {
	println("go-cleanstack")
}
```

**Step 4: Create .gitignore**

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

# Go
vendor/
```

**Step 5: Verify build**

```bash
go build -o bin/app ./cmd
./bin/app
```

Expected: Prints "go-cleanstack"

**Step 6: Commit**

```bash
git add .
git commit -m "feat: initialize go module and directory structure"
```

---

### Task 2: Add Configuration Layer

**Files:**
- Create: `go-cleanstack/internal/infra/config/config.go`
- Create: `go-cleanstack/internal/infra/config/config_test.go`
- Create: `go-cleanstack/config_development.toml`
- Create: `go-cleanstack/config_staging.toml.example`
- Create: `go-cleanstack/config_production.toml.example`
- Create: `go-cleanstack/.envrc.example`

**Step 1: Write failing test for config loading**

```go
// internal/infra/config/config_test.go
package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultsToDevEnvironment(t *testing.T) {
	// Setup: create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config_development.toml")
	err := os.WriteFile(configPath, []byte(`
[server]
port = 8080

[database]
url = "postgres://localhost/test"

[log]
level = "debug"
`), 0644)
	require.NoError(t, err)

	// Unset APP_ENV to test default
	os.Unsetenv("APP_ENV")

	cfg, err := Load(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "postgres://localhost/test", cfg.Database.URL)
	assert.Equal(t, "debug", cfg.Log.Level)
}
```

**Step 2: Run test to verify it fails**

```bash
cd /workspace/go-cleanstack
go mod tidy
gotestsum -- ./internal/infra/config/...
```

Expected: FAIL - config package doesn't exist

**Step 3: Add dependencies**

```bash
go get github.com/spf13/viper
go get github.com/stretchr/testify
```

**Step 4: Write minimal config implementation**

```go
// internal/infra/config/config.go
package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	URL string
}

type LogConfig struct {
	Level string
}

func Load(configPath string) (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName("config_" + env)
	viper.SetConfigType("toml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

**Step 5: Run test to verify it passes**

```bash
gotestsum -- ./internal/infra/config/...
```

Expected: PASS

**Step 6: Create config files**

```toml
# config_development.toml
[server]
port = 8080

[database]
url = "postgres://user:pass@localhost:5432/cleanstack?sslmode=disable"

[log]
level = "debug"
```

```toml
# config_staging.toml.example
[server]
port = 8080

[database]
url = "postgres://user:pass@staging-db:5432/cleanstack?sslmode=require"

[log]
level = "info"
```

```toml
# config_production.toml.example
[server]
port = 8080

[database]
url = "postgres://user:pass@prod-db:5432/cleanstack?sslmode=require"

[log]
level = "warn"
```

```bash
# .envrc.example
export APP_ENV=development
```

**Step 7: Commit**

```bash
git add .
git commit -m "feat: add configuration layer with Viper"
```

---

### Task 3: Add Domain Layer (Item Entity)

**Files:**
- Create: `go-cleanstack/internal/domain/entity/item.go`
- Create: `go-cleanstack/internal/domain/entity/item_test.go`
- Create: `go-cleanstack/internal/domain/repository/item.go`

**Step 1: Write failing test for Item entity**

```go
// internal/domain/entity/item_test.go
package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewItem_CreatesValidItem(t *testing.T) {
	item := NewItem("test-id", "Test Item", "A test item description")

	assert.Equal(t, "test-id", item.ID)
	assert.Equal(t, "Test Item", item.Name)
	assert.Equal(t, "A test item description", item.Description)
	assert.False(t, item.CreatedAt.IsZero())
}

func TestItem_Validate_RequiresName(t *testing.T) {
	item := &Item{
		ID:          "test-id",
		Name:        "",
		Description: "desc",
		CreatedAt:   time.Now(),
	}

	err := item.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestItem_Validate_AcceptsValidItem(t *testing.T) {
	item := NewItem("test-id", "Valid Name", "desc")

	err := item.Validate()
	assert.NoError(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
gotestsum -- ./internal/domain/entity/...
```

Expected: FAIL - Item type doesn't exist

**Step 3: Write minimal Item implementation**

```go
// internal/domain/entity/item.go
package entity

import (
	"errors"
	"time"
)

type Item struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
}

func NewItem(id, name, description string) *Item {
	return &Item{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
}

func (i *Item) Validate() error {
	if i.Name == "" {
		return errors.New("name is required")
	}
	return nil
}
```

**Step 4: Run test to verify it passes**

```bash
gotestsum -- ./internal/domain/entity/...
```

Expected: PASS

**Step 5: Create repository interface**

```go
// internal/domain/repository/item.go
package repository

import (
	"context"

	"github.com/user/go-cleanstack/internal/domain/entity"
)

type ItemRepository interface {
	Create(ctx context.Context, item *entity.Item) error
	GetByID(ctx context.Context, id string) (*entity.Item, error)
	List(ctx context.Context) ([]*entity.Item, error)
	Delete(ctx context.Context, id string) error
}
```

**Step 6: Verify build**

```bash
go build ./...
```

Expected: Build succeeds

**Step 7: Commit**

```bash
git add .
git commit -m "feat: add domain layer with Item entity and repository interface"
```

---

### Task 4: Add CLI Foundation (Cobra)

**Files:**
- Create: `go-cleanstack/internal/infra/cli/root.go`
- Create: `go-cleanstack/internal/infra/cli/version.go`
- Modify: `go-cleanstack/cmd/main.go`

**Step 1: Add Cobra dependency**

```bash
go get github.com/spf13/cobra
```

**Step 2: Create root command**

```go
// internal/infra/cli/root.go
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/go-cleanstack/internal/infra/config"
)

var cfg *config.Config

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cleanstack",
		Short: "Go-Cleanstack application",
		Long:  "A production-ready Go application with CLI and API",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cfg, err = config.Load(".")
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			return nil
		},
	}

	rootCmd.AddCommand(NewVersionCmd())

	return rootCmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func GetConfig() *config.Config {
	return cfg
}
```

**Step 3: Create version command**

```go
// internal/infra/cli/version.go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("go-cleanstack version %s (built %s)\n", Version, BuildTime)
		},
	}
}
```

**Step 4: Update main.go**

```go
// cmd/main.go
package main

import "github.com/user/go-cleanstack/internal/infra/cli"

func main() {
	cli.Execute()
}
```

**Step 5: Verify CLI works**

```bash
go run ./cmd version
```

Expected: Prints "go-cleanstack version dev (built unknown)"

**Step 6: Commit**

```bash
git add .
git commit -m "feat: add CLI foundation with Cobra"
```

---

### Task 5: Add Protobuf and Connect RPC Setup

**Files:**
- Create: `go-cleanstack/internal/infra/api/buf.yaml`
- Create: `go-cleanstack/internal/infra/api/buf.gen.yaml`
- Create: `go-cleanstack/internal/infra/api/proto/cleanstack/v1/item.proto`

**Step 1: Create buf.yaml**

```yaml
# internal/infra/api/buf.yaml
version: v2
modules:
  - path: proto
lint:
  use:
    - DEFAULT
breaking:
  use:
    - FILE
```

**Step 2: Create buf.gen.yaml**

```yaml
# internal/infra/api/buf.gen.yaml
version: v2
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/user/go-cleanstack/internal/infra/api/gen
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

**Step 3: Create item.proto**

```protobuf
// internal/infra/api/proto/cleanstack/v1/item.proto
syntax = "proto3";

package cleanstack.v1;

option go_package = "cleanstack/v1;cleanstackv1";

message Item {
  string id = 1;
  string name = 2;
  string description = 3;
  string created_at = 4;
}

message CreateItemRequest {
  string name = 1;
  string description = 2;
}

message CreateItemResponse {
  Item item = 1;
}

message GetItemRequest {
  string id = 1;
}

message GetItemResponse {
  Item item = 1;
}

message ListItemsRequest {}

message ListItemsResponse {
  repeated Item items = 1;
}

message DeleteItemRequest {
  string id = 1;
}

message DeleteItemResponse {}

service ItemService {
  rpc CreateItem(CreateItemRequest) returns (CreateItemResponse);
  rpc GetItem(GetItemRequest) returns (GetItemResponse);
  rpc ListItems(ListItemsRequest) returns (ListItemsResponse);
  rpc DeleteItem(DeleteItemRequest) returns (DeleteItemResponse);
}
```

**Step 4: Generate code with buf**

```bash
cd /workspace/go-cleanstack/internal/infra/api
buf generate
```

**Step 5: Verify generated files**

```bash
ls -la gen/cleanstack/v1/
```

Expected: See `item.pb.go` and `item_connect.go`

**Step 6: Update go.mod with connect dependencies**

```bash
cd /workspace/go-cleanstack
go mod tidy
```

**Step 7: Commit**

```bash
git add .
git commit -m "feat: add protobuf definitions and Connect RPC generation"
```

---

## Phase 2: Application and Infrastructure Layers

### Task 6: Add Application Service Layer

**Files:**
- Create: `go-cleanstack/internal/app/service/item_service.go`
- Create: `go-cleanstack/internal/app/service/item_service_test.go`

**Step 1: Write failing test for ItemService**

```go
// internal/app/service/item_service_test.go
package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-cleanstack/internal/domain/entity"
)

type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) Create(ctx context.Context, item *entity.Item) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockItemRepository) GetByID(ctx context.Context, id string) (*entity.Item, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockItemRepository) List(ctx context.Context) ([]*entity.Item, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Item), args.Error(1)
}

func (m *MockItemRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestItemService_CreateItem_Success(t *testing.T) {
	mockRepo := new(MockItemRepository)
	svc := NewItemService(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Item")).Return(nil)

	item, err := svc.CreateItem(context.Background(), "Test Item", "Test Description")

	assert.NoError(t, err)
	assert.Equal(t, "Test Item", item.Name)
	assert.Equal(t, "Test Description", item.Description)
	mockRepo.AssertExpectations(t)
}

func TestItemService_CreateItem_ValidationError(t *testing.T) {
	mockRepo := new(MockItemRepository)
	svc := NewItemService(mockRepo)

	item, err := svc.CreateItem(context.Background(), "", "Description")

	assert.Error(t, err)
	assert.Nil(t, item)
	mockRepo.AssertNotCalled(t, "Create")
}
```

**Step 2: Run test to verify it fails**

```bash
gotestsum -- ./internal/app/service/...
```

Expected: FAIL - ItemService doesn't exist

**Step 3: Write minimal ItemService implementation**

```go
// internal/app/service/item_service.go
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/go-cleanstack/internal/domain/entity"
	"github.com/user/go-cleanstack/internal/domain/repository"
)

type ItemService struct {
	repo repository.ItemRepository
}

func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(ctx context.Context, name, description string) (*entity.Item, error) {
	item := entity.NewItem(uuid.New().String(), name, description)

	if err := item.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemService) GetItem(ctx context.Context, id string) (*entity.Item, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ItemService) ListItems(ctx context.Context) ([]*entity.Item, error) {
	return s.repo.List(ctx)
}

func (s *ItemService) DeleteItem(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
```

**Step 4: Add uuid dependency and run tests**

```bash
go get github.com/google/uuid
gotestsum -- ./internal/app/service/...
```

Expected: PASS

**Step 5: Commit**

```bash
git add .
git commit -m "feat: add application service layer with ItemService"
```

---

### Task 7: Add Connect RPC Handler

**Files:**
- Create: `go-cleanstack/internal/infra/api/handler/item_handler.go`
- Create: `go-cleanstack/internal/infra/api/server.go`

**Step 1: Create Connect RPC handler**

```go
// internal/infra/api/handler/item_handler.go
package handler

import (
	"context"

	"connectrpc.com/connect"
	"github.com/user/go-cleanstack/internal/app/service"
	cleanstackv1 "github.com/user/go-cleanstack/internal/infra/api/gen/cleanstack/v1"
	"github.com/user/go-cleanstack/internal/infra/api/gen/cleanstack/v1/cleanstackv1connect"
)

type ItemHandler struct {
	service *service.ItemService
}

func NewItemHandler(svc *service.ItemService) *ItemHandler {
	return &ItemHandler{service: svc}
}

var _ cleanstackv1connect.ItemServiceHandler = (*ItemHandler)(nil)

func (h *ItemHandler) CreateItem(ctx context.Context, req *connect.Request[cleanstackv1.CreateItemRequest]) (*connect.Response[cleanstackv1.CreateItemResponse], error) {
	item, err := h.service.CreateItem(ctx, req.Msg.Name, req.Msg.Description)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&cleanstackv1.CreateItemResponse{
		Item: &cleanstackv1.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}), nil
}

func (h *ItemHandler) GetItem(ctx context.Context, req *connect.Request[cleanstackv1.GetItemRequest]) (*connect.Response[cleanstackv1.GetItemResponse], error) {
	item, err := h.service.GetItem(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&cleanstackv1.GetItemResponse{
		Item: &cleanstackv1.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}), nil
}

func (h *ItemHandler) ListItems(ctx context.Context, req *connect.Request[cleanstackv1.ListItemsRequest]) (*connect.Response[cleanstackv1.ListItemsResponse], error) {
	items, err := h.service.ListItems(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoItems := make([]*cleanstackv1.Item, len(items))
	for i, item := range items {
		protoItems[i] = &cleanstackv1.Item{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return connect.NewResponse(&cleanstackv1.ListItemsResponse{
		Items: protoItems,
	}), nil
}

func (h *ItemHandler) DeleteItem(ctx context.Context, req *connect.Request[cleanstackv1.DeleteItemRequest]) (*connect.Response[cleanstackv1.DeleteItemResponse], error) {
	if err := h.service.DeleteItem(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&cleanstackv1.DeleteItemResponse{}), nil
}
```

**Step 2: Create server setup**

```go
// internal/infra/api/server.go
package api

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"github.com/user/go-cleanstack/internal/app/service"
	"github.com/user/go-cleanstack/internal/infra/api/gen/cleanstack/v1/cleanstackv1connect"
	"github.com/user/go-cleanstack/internal/infra/api/handler"
)

type Server struct {
	port        int
	itemService *service.ItemService
}

func NewServer(port int, itemService *service.ItemService) *Server {
	return &Server{
		port:        port,
		itemService: itemService,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	itemHandler := handler.NewItemHandler(s.itemService)
	path, h := cleanstackv1connect.NewItemServiceHandler(itemHandler)
	mux.Handle(path, h)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Server starting on %s\n", addr)

	return http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{}))
}
```

**Step 3: Add connect dependencies**

```bash
go get connectrpc.com/connect
go get golang.org/x/net/http2
go get golang.org/x/net/http2/h2c
go mod tidy
```

**Step 4: Verify build**

```bash
go build ./...
```

Expected: Build succeeds

**Step 5: Commit**

```bash
git add .
git commit -m "feat: add Connect RPC handler and server"
```

---

### Task 8: Add Persistence Layer with sqlx

**Files:**
- Create: `go-cleanstack/internal/infra/persistence/item_repo.go`
- Create: `go-cleanstack/internal/infra/persistence/db.go`
- Create: `go-cleanstack/internal/infra/persistence/migrations/000001_create_items_table.up.sql`
- Create: `go-cleanstack/internal/infra/persistence/migrations/000001_create_items_table.down.sql`

**Step 1: Create database connection helper**

```go
// internal/infra/persistence/db.go
package persistence

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
```

**Step 2: Create migration files**

```sql
-- internal/infra/persistence/migrations/000001_create_items_table.up.sql
CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

```sql
-- internal/infra/persistence/migrations/000001_create_items_table.down.sql
DROP TABLE IF EXISTS items;
```

**Step 3: Create repository implementation**

```go
// internal/infra/persistence/item_repo.go
package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/user/go-cleanstack/internal/domain/entity"
	"github.com/user/go-cleanstack/internal/domain/repository"
)

type itemRow struct {
	ID          string       `db:"id"`
	Name        string       `db:"name"`
	Description string       `db:"description"`
	CreatedAt   sql.NullTime `db:"created_at"`
}

type ItemRepo struct {
	db *sqlx.DB
}

func NewItemRepo(db *sqlx.DB) *ItemRepo {
	return &ItemRepo{db: db}
}

var _ repository.ItemRepository = (*ItemRepo)(nil)

func (r *ItemRepo) Create(ctx context.Context, item *entity.Item) error {
	query := `INSERT INTO items (id, name, description, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, item.ID, item.Name, item.Description, item.CreatedAt)
	return err
}

func (r *ItemRepo) GetByID(ctx context.Context, id string) (*entity.Item, error) {
	var row itemRow
	query := `SELECT id, name, description, created_at FROM items WHERE id = $1`
	err := r.db.GetContext(ctx, &row, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("item not found")
	}
	if err != nil {
		return nil, err
	}

	return &entity.Item{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt.Time,
	}, nil
}

func (r *ItemRepo) List(ctx context.Context) ([]*entity.Item, error) {
	var rows []itemRow
	query := `SELECT id, name, description, created_at FROM items ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}

	items := make([]*entity.Item, len(rows))
	for i, row := range rows {
		items[i] = &entity.Item{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
			CreatedAt:   row.CreatedAt.Time,
		}
	}

	return items, nil
}

func (r *ItemRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
```

**Step 4: Add sqlx and postgres driver dependencies**

```bash
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
go mod tidy
```

**Step 5: Verify build**

```bash
go build ./...
```

Expected: Build succeeds

**Step 6: Commit**

```bash
git add .
git commit -m "feat: add persistence layer with sqlx repository"
```

---

### Task 9: Add Serve Command with Full Wiring

**Files:**
- Create: `go-cleanstack/internal/infra/cli/serve.go`
- Modify: `go-cleanstack/internal/infra/cli/root.go`

**Step 1: Create serve command**

```go
// internal/infra/cli/serve.go
package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/user/go-cleanstack/internal/app/service"
	"github.com/user/go-cleanstack/internal/infra/api"
	"github.com/user/go-cleanstack/internal/infra/persistence"
)

func NewServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := GetConfig()

			// Connect to database
			db, err := persistence.NewDB(cfg.Database.URL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			log.Println("Connected to database")

			// Wire up dependencies
			itemRepo := persistence.NewItemRepo(db)
			itemService := service.NewItemService(itemRepo)

			// Start server
			server := api.NewServer(cfg.Server.Port, itemService)
			return server.Start()
		},
	}
}
```

**Step 2: Update root.go to add serve command**

```go
// internal/infra/cli/root.go
// Add to NewRootCmd():
rootCmd.AddCommand(NewServeCmd())
```

**Step 3: Verify build**

```bash
go build ./...
```

Expected: Build succeeds

**Step 4: Commit**

```bash
git add .
git commit -m "feat: add serve command with full dependency wiring"
```

---

### Task 10: Add Migration Commands

**Files:**
- Create: `go-cleanstack/internal/infra/cli/migrate.go`
- Modify: `go-cleanstack/internal/infra/cli/root.go`

**Step 1: Create migrate commands**

```go
// internal/infra/cli/migrate.go
package cli

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func NewMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
	}

	migrateCmd.AddCommand(newMigrateUpCmd())
	migrateCmd.AddCommand(newMigrateDownCmd())

	return migrateCmd
}

func newMigrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run all pending migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := GetConfig()

			m, err := migrate.New(
				"file://internal/infra/persistence/migrations",
				cfg.Database.URL,
			)
			if err != nil {
				return fmt.Errorf("failed to create migrator: %w", err)
			}
			defer m.Close()

			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to run migrations: %w", err)
			}

			fmt.Println("Migrations completed successfully")
			return nil
		},
	}
}

func newMigrateDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := GetConfig()

			m, err := migrate.New(
				"file://internal/infra/persistence/migrations",
				cfg.Database.URL,
			)
			if err != nil {
				return fmt.Errorf("failed to create migrator: %w", err)
			}
			defer m.Close()

			if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to rollback migration: %w", err)
			}

			fmt.Println("Migration rolled back successfully")
			return nil
		},
	}
}
```

**Step 2: Update root.go to add migrate command**

```go
// internal/infra/cli/root.go
// Add to NewRootCmd():
rootCmd.AddCommand(NewMigrateCmd())
```

**Step 3: Add golang-migrate dependency**

```bash
go get github.com/golang-migrate/migrate/v4
go mod tidy
```

**Step 4: Verify build**

```bash
go build ./...
```

Expected: Build succeeds

**Step 5: Commit**

```bash
git add .
git commit -m "feat: add database migration commands"
```

---

## Phase 3: Tooling and DevOps

### Task 11: Add justfile

**Files:**
- Create: `go-cleanstack/justfile`

**Step 1: Create justfile**

```just
# go-cleanstack justfile

# Development
dev:
    go run ./cmd serve

# Code generation
generate:
    cd internal/infra/api && buf generate

# Testing
test:
    gotestsum -- ./...

test-int:
    gotestsum -- -tags=integration ./tests/integration/...

test-e2e:
    gotestsum -- -tags=e2e ./tests/e2e/...

test-all:
    gotestsum -- -tags=integration,e2e ./...

test-cover:
    gotestsum -- -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

# Database
migrate-up:
    go run ./cmd migrate up

migrate-down:
    go run ./cmd migrate down

# Linting
lint:
    golangci-lint run

lint-fix:
    golangci-lint run --fix

# Build
build:
    go build -o bin/cleanstack ./cmd

# Docker
up:
    docker-compose up -d

down:
    docker-compose down

logs:
    docker-compose logs -f

# Clean
clean:
    rm -rf bin/ coverage.out coverage.html
```

**Step 2: Verify justfile works**

```bash
just --list
```

Expected: Shows all available recipes

**Step 3: Commit**

```bash
git add justfile
git commit -m "feat: add justfile for task automation"
```

---

### Task 12: Add golangci-lint Configuration

**Files:**
- Create: `go-cleanstack/.golangci.yml`

**Step 1: Create .golangci.yml**

```yaml
# .golangci.yml
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
  gocritic:
    enabled-tags:
      - diagnostic
      - style
    disabled-checks:
      - ifElseChain

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
    - path: gen/
      linters:
        - all
```

**Step 2: Run linter**

```bash
cd /workspace/go-cleanstack
golangci-lint run
```

Expected: Either PASS or shows fixable issues

**Step 3: Fix any issues if present**

**Step 4: Commit**

```bash
git add .golangci.yml
git commit -m "feat: add golangci-lint configuration"
```

---

### Task 13: Add Pre-commit Hooks

**Files:**
- Create: `go-cleanstack/.pre-commit-config.yaml`

**Step 1: Create .pre-commit-config.yaml**

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/golangci/golangci-lint
    rev: v1.61.0
    hooks:
      - id: golangci-lint
        args: [--fix]

  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files

  - repo: local
    hooks:
      - id: buf-lint
        name: buf lint
        entry: bash -c 'cd internal/infra/api && buf lint'
        language: system
        files: '\.proto$'
        pass_filenames: false

      - id: buf-generate
        name: buf generate
        entry: bash -c 'cd internal/infra/api && buf generate'
        language: system
        files: '\.proto$'
        pass_filenames: false
```

**Step 2: Commit**

```bash
git add .pre-commit-config.yaml
git commit -m "feat: add pre-commit hooks configuration"
```

---

### Task 14: Add Docker Configuration

**Files:**
- Create: `go-cleanstack/Dockerfile`
- Create: `go-cleanstack/docker-compose.yml`

**Step 1: Create Dockerfile**

```dockerfile
# Dockerfile

# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/cleanstack ./cmd

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/cleanstack /usr/local/bin/cleanstack
COPY --from=builder /app/internal/infra/persistence/migrations /migrations
COPY --from=builder /app/config_development.toml /config_development.toml

EXPOSE 8080

ENTRYPOINT ["cleanstack"]
CMD ["serve"]
```

**Step 2: Create docker-compose.yml**

```yaml
# docker-compose.yml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
    depends_on:
      db:
        condition: service_healthy
    volumes:
      - ./config_development.toml:/config_development.toml:ro

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: cleanstack
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d cleanstack"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
```

**Step 3: Verify docker-compose config**

```bash
docker-compose config
```

Expected: Valid YAML output

**Step 4: Commit**

```bash
git add Dockerfile docker-compose.yml
git commit -m "feat: add Docker and docker-compose configuration"
```

---

## Phase 4: Testing

### Task 15: Add Test Utilities

**Files:**
- Create: `go-cleanstack/tests/testutil/db.go`
- Create: `go-cleanstack/tests/testutil/containers.go`

**Step 1: Create test container utilities**

```go
// tests/testutil/containers.go
//go:build integration || e2e

package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Container testcontainers.Container
	URI       string
}

func StartPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	return &PostgresContainer{
		Container: container,
		URI:       uri,
	}, nil
}

func (p *PostgresContainer) Terminate(ctx context.Context) error {
	return p.Container.Terminate(ctx)
}
```

**Step 2: Create db test utilities**

```go
// tests/testutil/db.go
//go:build integration || e2e

package testutil

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func SetupTestDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Run migrations
	m, err := migrate.New(
		"file://../../internal/infra/persistence/migrations",
		databaseURL,
	)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return db, nil
}

func CleanupTestDB(db *sqlx.DB) {
	db.Exec("TRUNCATE TABLE items CASCADE")
}
```

**Step 3: Add testcontainers dependency**

```bash
go get github.com/testcontainers/testcontainers-go
go mod tidy
```

**Step 4: Commit**

```bash
git add tests/
git commit -m "feat: add test utilities with testcontainers"
```

---

### Task 16: Add Integration Tests

**Files:**
- Create: `go-cleanstack/tests/integration/item_repo_test.go`

**Step 1: Create integration test for ItemRepo**

```go
// tests/integration/item_repo_test.go
//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/go-cleanstack/internal/domain/entity"
	"github.com/user/go-cleanstack/internal/infra/persistence"
	"github.com/user/go-cleanstack/tests/testutil"
)

func TestItemRepo_CRUD(t *testing.T) {
	ctx := context.Background()

	// Start postgres container
	pgContainer, err := testutil.StartPostgresContainer(ctx)
	require.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	// Setup test DB
	db, err := testutil.SetupTestDB(pgContainer.URI)
	require.NoError(t, err)
	defer db.Close()

	repo := persistence.NewItemRepo(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		item := entity.NewItem("test-id-1", "Test Item", "Test Description")
		err := repo.Create(ctx, item)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, "test-id-1")
		require.NoError(t, err)
		assert.Equal(t, item.ID, retrieved.ID)
		assert.Equal(t, item.Name, retrieved.Name)
		assert.Equal(t, item.Description, retrieved.Description)
	})

	t.Run("List", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		item1 := entity.NewItem("test-id-2", "Item 1", "Desc 1")
		item2 := entity.NewItem("test-id-3", "Item 2", "Desc 2")
		require.NoError(t, repo.Create(ctx, item1))
		require.NoError(t, repo.Create(ctx, item2))

		items, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		item := entity.NewItem("test-id-4", "To Delete", "Will be deleted")
		require.NoError(t, repo.Create(ctx, item))

		err := repo.Delete(ctx, "test-id-4")
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, "test-id-4")
		assert.Error(t, err)
	})
}
```

**Step 2: Run integration tests (requires Docker)**

```bash
gotestsum -- -tags=integration ./tests/integration/...
```

Expected: PASS (if Docker is available)

**Step 3: Commit**

```bash
git add tests/integration/
git commit -m "feat: add integration tests for ItemRepo"
```

---

### Task 17: Add E2E Tests

**Files:**
- Create: `go-cleanstack/tests/e2e/api_test.go`

**Step 1: Create E2E test for API**

```go
// tests/e2e/api_test.go
//go:build e2e

package e2e

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/go-cleanstack/internal/app/service"
	cleanstackv1 "github.com/user/go-cleanstack/internal/infra/api/gen/cleanstack/v1"
	"github.com/user/go-cleanstack/internal/infra/api/gen/cleanstack/v1/cleanstackv1connect"
	"github.com/user/go-cleanstack/internal/infra/api/handler"
	"github.com/user/go-cleanstack/internal/infra/persistence"
	"github.com/user/go-cleanstack/tests/testutil"
)

func TestItemAPI_E2E(t *testing.T) {
	ctx := context.Background()

	// Start postgres container
	pgContainer, err := testutil.StartPostgresContainer(ctx)
	require.NoError(t, err)
	defer pgContainer.Terminate(ctx)

	// Setup test DB
	db, err := testutil.SetupTestDB(pgContainer.URI)
	require.NoError(t, err)
	defer db.Close()

	// Wire up dependencies
	itemRepo := persistence.NewItemRepo(db)
	itemService := service.NewItemService(itemRepo)
	itemHandler := handler.NewItemHandler(itemService)

	// Create test server
	mux := http.NewServeMux()
	path, h := cleanstackv1connect.NewItemServiceHandler(itemHandler)
	mux.Handle(path, h)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Create client
	client := cleanstackv1connect.NewItemServiceClient(
		http.DefaultClient,
		server.URL,
	)

	t.Run("CreateItem", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		resp, err := client.CreateItem(ctx, connect.NewRequest(&cleanstackv1.CreateItemRequest{
			Name:        "E2E Test Item",
			Description: "Created via E2E test",
		}))

		require.NoError(t, err)
		assert.NotEmpty(t, resp.Msg.Item.Id)
		assert.Equal(t, "E2E Test Item", resp.Msg.Item.Name)
	})

	t.Run("ListItems", func(t *testing.T) {
		testutil.CleanupTestDB(db)

		// Create items first
		_, err := client.CreateItem(ctx, connect.NewRequest(&cleanstackv1.CreateItemRequest{
			Name:        "Item 1",
			Description: "First item",
		}))
		require.NoError(t, err)

		_, err = client.CreateItem(ctx, connect.NewRequest(&cleanstackv1.CreateItemRequest{
			Name:        "Item 2",
			Description: "Second item",
		}))
		require.NoError(t, err)

		// List items
		resp, err := client.ListItems(ctx, connect.NewRequest(&cleanstackv1.ListItemsRequest{}))
		require.NoError(t, err)
		assert.Len(t, resp.Msg.Items, 2)
	})
}
```

**Step 2: Run E2E tests (requires Docker)**

```bash
gotestsum -- -tags=e2e ./tests/e2e/...
```

Expected: PASS (if Docker is available)

**Step 3: Commit**

```bash
git add tests/e2e/
git commit -m "feat: add E2E tests for API"
```

---

### Task 18: Add README

**Files:**
- Create: `go-cleanstack/README.md`

**Step 1: Create README.md**

```markdown
# Go-Cleanstack

A production-ready Go project skeleton with CLI, Connect RPC API, DDD architecture, and full test pyramid.

## Features

- **CLI**: Cobra-based command line interface
- **API**: Connect RPC with protobuf definitions
- **Architecture**: Simplified DDD with layered architecture
- **Database**: sqlx with golang-migrate migrations
- **Testing**: Full test pyramid (unit, integration, e2e)
- **Linting**: golangci-lint with pre-commit hooks
- **Docker**: Multi-stage build with docker-compose

## Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- [buf](https://buf.build/docs/cli/installation)
- [just](https://github.com/casey/just)
- [gotestsum](https://github.com/gotestyourself/gotestsum)
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [direnv](https://direnv.net/)

### Setup

```bash
# Clone repository
git clone https://github.com/user/go-cleanstack.git
cd go-cleanstack

# Setup environment
cp .envrc.example .envrc
direnv allow

# Start database
just up

# Run migrations
just migrate-up

# Start server
just dev
```

### Development Commands

```bash
just dev          # Run development server
just generate     # Generate protobuf code
just test         # Run unit tests
just test-int     # Run integration tests
just test-e2e     # Run E2E tests
just test-all     # Run all tests
just lint         # Run linter
just build        # Build binary
```

## Project Structure

```
go-cleanstack/
├── cmd/                          # Application entrypoint
├── internal/
│   ├── domain/                   # Business logic (no external deps)
│   │   ├── entity/               # Domain entities
│   │   └── repository/           # Repository interfaces
│   ├── app/                      # Application layer
│   │   └── service/              # Use cases / services
│   └── infra/                    # Infrastructure layer
│       ├── api/                  # Connect RPC API
│       ├── cli/                  # Cobra commands
│       ├── persistence/          # Database repositories
│       └── config/               # Configuration
├── tests/
│   ├── integration/              # Integration tests
│   ├── e2e/                      # End-to-end tests
│   └── testutil/                 # Test utilities
├── config_development.toml       # Dev configuration
├── docker-compose.yml            # Local development stack
└── justfile                      # Task automation
```

## Configuration

Configuration is loaded based on `APP_ENV` environment variable:

- `development` (default): `config_development.toml`
- `staging`: `config_staging.toml`
- `production`: `config_production.toml`

## License

MIT
```

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add README with setup instructions"
```

---

## Summary

**Total Tasks:** 18

**Phase 1 (Foundation):** Tasks 1-5
- Go module, directories, config, domain, CLI, protobuf

**Phase 2 (Layers):** Tasks 6-10
- Application service, Connect handler, persistence, serve command, migrations

**Phase 3 (Tooling):** Tasks 11-14
- justfile, golangci-lint, pre-commit, Docker

**Phase 4 (Testing):** Tasks 15-18
- Test utilities, integration tests, E2E tests, README

---

Plan complete and saved to `docs/plans/2025-12-21-go-cleanstack-implementation.md`.

**Two execution options:**

1. **Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

2. **Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

Which approach?
