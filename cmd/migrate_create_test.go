package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMigrationFile(t *testing.T) {
	// Setup temp directory
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		description   string
		migrationType migrationType
		wantErr       bool
		wantFileExt   string
	}{
		{
			name:          "create SQL migration",
			description:   "add user table",
			migrationType: sqlMigration,
			wantErr:       false,
			wantFileExt:   ".sql",
		},
		{
			name:          "create Go migration",
			description:   "seed initial data",
			migrationType: goMigration,
			wantErr:       false,
			wantFileExt:   ".go",
		},
		{
			name:          "invalid description",
			description:   "ab",
			migrationType: sqlMigration,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, err := createMigrationFile(tmpDir, tt.description, tt.migrationType)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, filePath)

			// Verify file exists
			_, err = os.Stat(filePath)
			require.NoError(t, err)

			// Verify file extension
			assert.Equal(t, tt.wantFileExt, filepath.Ext(filePath))

			// Verify file content
			content, err := os.ReadFile(filePath)
			require.NoError(t, err)
			assert.NotEmpty(t, content)

			if tt.migrationType == sqlMigration {
				assert.Contains(t, string(content), "-- +goose Up")
				assert.Contains(t, string(content), "-- +goose Down")
			} else {
				assert.Contains(t, string(content), "package migrations")
				assert.Contains(t, string(content), "goose.AddMigration")
			}
		})
	}
}
