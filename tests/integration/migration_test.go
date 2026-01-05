//go:build integration

package integration

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMigrations(t *testing.T) {
	ctx := context.Background()

	// Start PostgreSQL container with robust wait strategy
	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
				wait.ForListeningPort("5432/tcp"),
			).WithDeadline(60 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer postgres.Terminate(ctx)

	// Get connection string
	host, err := postgres.Host(ctx)
	require.NoError(t, err)
	port, err := postgres.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable"

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Get migrations path
	migrationsPath := filepath.Join("..", "..", "internal", "app", "app1", "infra", "persistence", "migrations")

	// Test: Apply migrations
	t.Run("apply migrations", func(t *testing.T) {
		err := goose.UpContext(ctx, db, migrationsPath)
		assert.NoError(t, err)

		// Verify items table exists
		var exists bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'items'
			)
		`).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "items table should exist after migration")
	})

	// Test: Rollback migrations
	t.Run("rollback migrations", func(t *testing.T) {
		err := goose.DownContext(ctx, db, migrationsPath)
		assert.NoError(t, err)

		// Verify items table does not exist
		var exists bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'items'
			)
		`).Scan(&exists)
		require.NoError(t, err)
		assert.False(t, exists, "items table should not exist after rollback")
	})

	// Test: Reapply migrations
	t.Run("reapply migrations", func(t *testing.T) {
		err := goose.UpContext(ctx, db, migrationsPath)
		assert.NoError(t, err)

		// Verify items table exists again
		var exists bool
		err = db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = 'items'
			)
		`).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists, "items table should exist after reapplying migration")
	})

	// Test: Get version
	t.Run("get version", func(t *testing.T) {
		version, err := goose.GetDBVersionContext(ctx, db)
		assert.NoError(t, err)
		assert.Greater(t, version, int64(0), "version should be greater than 0")
	})
}
