package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/pivaldi/go-cleanstack/internal/app/user/infra/persistence/migrations"
)

const sqlMigrationTemplate = `-- +goose Up
-- +goose StatementBegin
-- TODO: Add your schema changes here

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- TODO: Add rollback logic here

-- +goose StatementEnd
`

const goMigrationTemplate = `package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(up{{.CamelCase}}, down{{.CamelCase}})
}

func up{{.CamelCase}}(tx *sql.Tx) error {
	// TODO: Implement migration logic
	return nil
}

func down{{.CamelCase}}(tx *sql.Tx) error {
	// TODO: Implement rollback logic
	return nil
}
`

// renderSQLTemplate renders the SQL migration template.
func renderSQLTemplate() (string, error) {
	return sqlMigrationTemplate, nil
}

// renderGoTemplate renders the Go migration template with the given description.
func renderGoTemplate(description string) (string, error) {
	tmpl, err := template.New("migration").Parse(goMigrationTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := map[string]string{
		"CamelCase": migrations.ToCamelCase(description),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
