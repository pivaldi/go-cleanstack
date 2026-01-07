//go:build integration || e2e

package testutil

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/pivaldi/go-cleanstack/internal/app/user/infra/persistence/migrations"
)

// SetupTestDB connects to the database and runs migrations.
func SetupTestDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations using embedded files
	ctx := context.Background()
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db.DB, "."); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// CleanupTestDB truncates all tables to reset state between tests.
func CleanupTestDB(db *sqlx.DB) {
	_, _ = db.Exec("TRUNCATE TABLE users CASCADE")
}
