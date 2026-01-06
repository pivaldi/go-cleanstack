package persistence

import (
	"database/sql"
	"time"
)

// UserDTO represents the database model for users.
// This is infrastructure-specific and has no dependency on domain.
type UserDTO struct {
	ID        int64
	Email     string
	Password  string
	FirstName sql.NullString
	LastName  sql.NullString
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

// userRow is the internal struct for sqlx scanning.
type userRow struct {
	ID        int64          `db:"id"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	FirstName sql.NullString `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	Role      string         `db:"role"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}

func (r *userRow) toDTO() *UserDTO {
	return &UserDTO{
		ID:        r.ID,
		Email:     r.Email,
		Password:  r.Password,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Role:      r.Role,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt,
	}
}
