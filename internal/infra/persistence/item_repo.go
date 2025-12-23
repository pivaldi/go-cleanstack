package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type itemRow struct {
	ID          string       `db:"id"`
	Name        string       `db:"name"`
	Description string       `db:"description"`
	CreatedAt   sql.NullTime `db:"created_at"`
}

// ItemRepo is the infrastructure implementation
// It works only with DTOs and has NO dependency on domain
type ItemRepo struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewItemRepo(db *sqlx.DB, logger *zap.Logger) *ItemRepo {
	return &ItemRepo{
		db:     db,
		logger: logger,
	}
}

func (r *ItemRepo) Create(ctx context.Context, item *ItemDTO) error {
	query := `INSERT INTO items (id, name, description, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, item.ID, item.Name, item.Description, item.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

func (r *ItemRepo) GetByID(ctx context.Context, id string) (*ItemDTO, error) {
	var row itemRow
	query := `SELECT id, name, description, created_at FROM items WHERE id = $1`
	err := r.db.GetContext(ctx, &row, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return &ItemDTO{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt.Time,
	}, nil
}

func (r *ItemRepo) List(ctx context.Context) ([]*ItemDTO, error) {
	var rows []itemRow
	query := `SELECT id, name, description, created_at FROM items ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	items := make([]*ItemDTO, len(rows))
	for i, row := range rows {
		items[i] = &ItemDTO{
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
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}
