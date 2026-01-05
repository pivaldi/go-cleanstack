package cmd

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/pivaldi/go-cleanstack/internal/infra/config"
	"github.com/pivaldi/go-cleanstack/internal/infra/persistence/migrations"
	stringpkg "github.com/pivaldi/go-cleanstack/pkg/string"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

type migrationType int

const (
	sqlMigration migrationType = iota
	goMigration
)

const (
	defaultMigrationsDir = "internal/infra/persistence/migrations"
)

var (
	migrationsDir = defaultMigrationsDir
)

// SetMigrationsDir sets the migrations directory (for testing).
func SetMigrationsDir(dir string) {
	migrationsDir = dir
}

// NewMigrateCmd creates the migrate command with subcommands.
func NewMigrateCmd(cfg *config.Config) *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
	}

	migrateCmd.AddCommand(newMigrateUpCmd(cfg))
	migrateCmd.AddCommand(newMigrateDownCmd(cfg))
	migrateCmd.AddCommand(newMigrateStatusCmd(cfg))
	migrateCmd.AddCommand(newMigrateVersionCmd(cfg))
	migrateCmd.AddCommand(newMigrateCreateCmd())

	// Set verbose mode to show SQL
	goose.SetVerbose(true)

	return migrateCmd
}

func newMigrateUpCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run all pending migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			db, err := getDBConnection(cfg)
			if err != nil {
				return err
			}
			defer db.Close()

			// Run migrations (use embedded or filesystem)
			if err := runMigrations(ctx, db, goose.UpContext); err != nil {
				return fmt.Errorf("failed to run migrations: %w", err)
			}

			// Get new version
			newVersion, err := goose.GetDBVersionContext(ctx, db)
			if err != nil {
				return fmt.Errorf("failed to get new version: %w", err)
			}

			fmt.Printf("\nMigrations complete. Database at version %d\n", newVersion)

			return nil
		},
	}
}

func newMigrateDownCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			db, err := getDBConnection(cfg)
			if err != nil {
				return err
			}
			defer db.Close()

			// Get current version
			currentVersion, err := goose.GetDBVersionContext(ctx, db)
			if err != nil {
				return fmt.Errorf("failed to get current version: %w", err)
			}

			fmt.Printf("Current database version: %d\n", currentVersion)

			// Rollback one migration (use embedded or filesystem)
			if err := runMigrations(ctx, db, goose.DownContext); err != nil {
				return fmt.Errorf("failed to rollback migration: %w", err)
			}

			// Get new version
			newVersion, err := goose.GetDBVersionContext(ctx, db)
			if err != nil {
				return fmt.Errorf("failed to get new version: %w", err)
			}

			fmt.Printf("\nRollback complete. Database at version %d\n", newVersion)

			return nil
		},
	}
}

func newMigrateStatusCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			db, err := getDBConnection(cfg)
			if err != nil {
				return err
			}
			defer db.Close()

			// Get migration status (use embedded or filesystem)
			if err := runMigrations(ctx, db, goose.StatusContext); err != nil {
				return fmt.Errorf("failed to get migration status: %w", err)
			}

			return nil
		},
	}
}

func newMigrateVersionCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show current database version",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx := context.Background()
			db, err := getDBConnection(cfg)
			if err != nil {
				return err
			}
			defer db.Close()

			version, err := goose.GetDBVersionContext(ctx, db)
			if err != nil {
				return fmt.Errorf("failed to get database version: %w", err)
			}

			fmt.Printf("Database version: %d\n", version)

			return nil
		},
	}
}

func newMigrateCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a new migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			reader := bufio.NewReader(os.Stdin)

			// Prompt for description
			fmt.Print("Migration description: ")
			description, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read description: %w", err)
			}
			description = strings.TrimSpace(description)

			// Validate description
			if err := migrations.ValidateDescription(description); err != nil {
				return fmt.Errorf("error retrieving description: %w", err)
			}

			// Prompt for migration type
			fmt.Println("Migration type: [1] SQL  [2] Go")
			fmt.Print("> ")
			choiceStr, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read migration type: %w", err)
			}
			choiceStr = strings.TrimSpace(choiceStr)

			var migrationType migrationType
			switch choiceStr {
			case "1":
				migrationType = sqlMigration
			case "2":
				migrationType = goMigration
			default:
				return errors.New("invalid migration type, choose 1 or 2")
			}

			// Create migration file
			filePath, err := createMigrationFile(migrationsDir, description, migrationType)
			if err != nil {
				return err
			}

			fmt.Printf("\nCreated: %s\n", filePath)

			return nil
		},
	}
}

// createMigrationFile creates a new migration file.
func createMigrationFile(dir, description string, migrationType migrationType) (string, error) {
	// Validate description
	if err := migrations.ValidateDescription(description); err != nil {
		return "", fmt.Errorf("migration description is not valid: %w", err)
	}

	// Normalize description for filename using existing NormalizeFileName
	normalizedDesc, err := stringpkg.NormalizeFileName(description)
	if err != nil {
		return "", fmt.Errorf("failed to normalize description: %w", err)
	}

	// Generate timestamp
	timestamp := time.Now().Format("20060102150405")

	// Generate filename
	var filename string
	if migrationType == sqlMigration {
		filename = fmt.Sprintf("%s_%s.sql", timestamp, normalizedDesc)
	} else {
		filename = fmt.Sprintf("%s_%s.go", timestamp, normalizedDesc)
	}

	filePath := filepath.Join(dir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Render template
	var content string
	if migrationType == sqlMigration {
		content, err = renderSQLTemplate()
	} else {
		content, err = renderGoTemplate(normalizedDesc)
	}
	if err != nil {
		return "", err
	}

	//
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		return "", fmt.Errorf("failed to write migration file: %w", err)
	}

	// Get absolute path for output
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		absPath = filePath
	}

	return absPath, nil
}

// getDBConnection returns a database connection.
func getDBConnection(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("configuration not loaded")
	}

	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(context.TODO()); err != nil {
		db.Close()

		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// runMigrations executes migrations using embedded FS or filesystem depending on availability.
type migrationFunc func(context.Context, *sql.DB, string, ...goose.OptionsFunc) error

func runMigrations(ctx context.Context, db *sql.DB, fn migrationFunc) error {
	// Check if external migrations directory exists
	if _, err := os.Stat(migrationsDir); err == nil {
		// Use filesystem migrations (development mode)
		return fn(ctx, db, migrationsDir)
	}

	// Use embedded migrations (production mode)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	goose.SetBaseFS(migrations.FS)

	return fn(ctx, db, ".")
}
