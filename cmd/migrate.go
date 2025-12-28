package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	stringpkg "github.com/pivaldi/go-cleanstack/pkg/string"
	"github.com/spf13/cobra"
)

var (
	migrationsDir = "migrations"
)

type migrationLogger struct{}

func (l *migrationLogger) Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (l *migrationLogger) Verbose() bool {
	return true
}

func newMigrator() (*migrate.Migrate, error) {
	m, err := migrate.New(
		"file://"+migrationsDir,
		cfg.Database.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	m.Log = &migrationLogger{}

	return m, nil
}

func SetMigrationsDir(dir string) {
	migrationsDir = dir
}

func NewMigrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
	}

	migrateCmd.AddCommand(newMigrateUpCmd())
	migrateCmd.AddCommand(newMigrateDownCmd())
	migrateCmd.AddCommand(newMigrateCreateCmd())

	return migrateCmd
}

func newMigrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run all pending migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			m, err := newMigrator()

			if err != nil {
				return err
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
			m, err := newMigrator()

			if err != nil {
				return err
			}
			defer m.Close()

			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("failed to rollback migration: %w", err)
			}

			fmt.Println("Migration rolled back successfully")

			return nil
		},
	}
}

func newMigrateCreateCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			if name == "" {
				return errors.New("migration name is required")
			}

			ut := time.Now().UnixNano()

			err := createMigrationFile(name, ut)
			if err != nil {
				return fmt.Errorf("failed to create migration file: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the migration")

	return cmd
}

// createMigrationFile implements the logic to create a migration file with the given name and timestamp
// Return an error if the file creation fails
func createMigrationFile(name string, timestamp int64) error {
	format := "%s_%s.%s.sql"
	lastCreated := ""
	version := strconv.FormatInt(timestamp, 10)

	for _, action := range []string{"up", "down"} {
		fullName := fmt.Sprintf(format, version, name, action)
		path, err := stringpkg.NormalizeFileName(filepath.Clean(filepath.Join(migrationsDir, fullName)))
		if err != nil {
			return fmt.Errorf("failed to normalize migration file name: %w", err)
		}

		file, err := os.Create(path)
		if err != nil {
			if lastCreated != "" {
				if errr := os.Remove(lastCreated); errr != nil {
					return fmt.Errorf("failed to remove partial migration %s : %w", path, err)
				}
			}

			return fmt.Errorf("failed to create %s migration file: %w", path, err)
		}

		err = file.Close()
		if err != nil {
			return fmt.Errorf("failed to close migration file: %w", err)
		}

		lastCreated = path
		fmt.Printf("Migration created successfully: %s\n", path)
	}

	return nil
}
