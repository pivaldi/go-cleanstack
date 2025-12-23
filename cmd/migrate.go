package cmd

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
		RunE: func(_ *cobra.Command, _ []string) error {
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
		RunE: func(_ *cobra.Command, _ []string) error {
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
