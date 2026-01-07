package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pivaldi/presence"
)

var ErrUserNotFound = errors.New("user not found")

// UserRepo is the infrastructure implementation.
type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *User) (*User, error) {
	query := `
		INSERT INTO users (email, password, first_name, last_name, role, created_at)
		VALUES ($1, crypt($2, gen_salt('bf')), $3, $4, $5, NOW())
		RETURNING id, email, password, first_name, last_name, role, created_at, updated_at, deleted_at
	`

	var result User
	err := r.db.GetContext(ctx, &result, query,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	return &result, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, role, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) List(ctx context.Context, offset, limit int) ([]*User, int64, error) {
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

	var users []User
	if err := r.db.SelectContext(ctx, &users, query, limit, offset); err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	result := make([]*User, len(users))
	for i := range users {
		result[i] = &users[i]
	}

	return result, total, nil
}

func (r *UserRepo) Update(ctx context.Context, user *User) (*User, error) {
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

	var result User
	err := r.db.GetContext(ctx, &result, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query: %w", err)
	}

	return &result, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, presence.FromValue(time.Now()))
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
