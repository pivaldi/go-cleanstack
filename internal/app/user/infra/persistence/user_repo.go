package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepo is the infrastructure implementation.
// It works only with DTOs and has NO dependency on domain.
type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, dto *UserDTO) (*UserDTO, error) {
	query := `
		INSERT INTO users (email, password, first_name, last_name, role, created_at)
		VALUES ($1, crypt($2, gen_salt('bf')), $3, $4, $5, NOW())
		RETURNING id, email, password, first_name, last_name, role, created_at, updated_at, deleted_at
	`

	var row userRow
	err := r.db.GetContext(ctx, &row, query,
		dto.Email,
		dto.Password,
		dto.FirstName,
		dto.LastName,
		dto.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	return row.toDTO(), nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*UserDTO, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var row userRow
	err := r.db.GetContext(ctx, &row, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return row.toDTO(), nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*UserDTO, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var row userRow
	err := r.db.GetContext(ctx, &row, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return row.toDTO(), nil
}

func (r *UserRepo) List(ctx context.Context, offset, limit int) ([]*UserDTO, int64, error) {
	// Get total count
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated results
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var rows []userRow
	if err := r.db.SelectContext(ctx, &rows, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	dtos := make([]*UserDTO, len(rows))
	for i, row := range rows {
		dtos[i] = row.toDTO()
	}

	return dtos, total, nil
}

func (r *UserRepo) Update(ctx context.Context, dto *UserDTO) (*UserDTO, error) {
	query := `
		UPDATE users SET
			email = COALESCE(NULLIF($2, ''), email),
			password = CASE WHEN $3 = '' THEN password ELSE crypt($3, gen_salt('bf')) END,
			first_name = $4,
			last_name = $5,
			role = COALESCE(NULLIF($6, ''), role),
			updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, email, password, first_name, last_name, role, created_at, updated_at, deleted_at
	`

	var row userRow
	err := r.db.GetContext(ctx, &row, query,
		dto.ID,
		dto.Email,
		dto.Password,
		dto.FirstName,
		dto.LastName,
		dto.Role,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}

	return row.toDTO(), nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to execute soft delete: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
