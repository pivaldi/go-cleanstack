//go:build integration || e2e

package testutil

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func SetupTestDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Run migrations using Goose
	ctx := context.Background()
	migrationsPath := "../../internal/app/app1/infra/persistence/migrations"

	if err := goose.UpContext(ctx, db.DB, migrationsPath); err != nil {
		return nil, err
	}

	return db, nil
}

func CleanupTestDB(db *sqlx.DB) {
	db.Exec("TRUNCATE TABLE items CASCADE")
}
