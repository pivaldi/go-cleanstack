package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderSQLTemplate(t *testing.T) {
	result, err := renderSQLTemplate()
	require.NoError(t, err)

	assert.Contains(t, result, "-- +goose Up")
	assert.Contains(t, result, "-- +goose Down")
	assert.Contains(t, result, "-- +goose StatementBegin")
	assert.Contains(t, result, "-- +goose StatementEnd")
	assert.Contains(t, result, "-- TODO: Add your schema changes here")
	assert.Contains(t, result, "-- TODO: Add rollback logic here")
}

func TestRenderGoTemplate(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantContain []string
	}{
		{
			name:        "simple description",
			description: "add-user-table",
			wantContain: []string{
				"package migrations",
				"func init()",
				"goose.AddMigration(upAddUserTable, downAddUserTable)",
				"func upAddUserTable(tx *sql.Tx) error",
				"func downAddUserTable(tx *sql.Tx) error",
			},
		},
		{
			name:        "complex description",
			description: "add-user-authentication-system",
			wantContain: []string{
				"goose.AddMigration(upAddUserAuthenticationSystem, downAddUserAuthenticationSystem)",
				"func upAddUserAuthenticationSystem(tx *sql.Tx) error",
				"func downAddUserAuthenticationSystem(tx *sql.Tx) error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renderGoTemplate(tt.description)
			require.NoError(t, err)

			for _, want := range tt.wantContain {
				assert.Contains(t, result, want)
			}
		})
	}
}
