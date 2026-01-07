package entity

import (
	"errors"
	"net/mail"
	"time"

	"github.com/pivaldi/presence"
)

const minPasswordLength = 8

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailInvalid     = errors.New("email format is invalid")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrRoleInvalid      = errors.New("role is invalid")
)

type User struct {
	ID        int64                  `db:"id"`
	Email     string                 `db:"email"`
	Password  string                 `db:"password"` // hashed by PostgreSQL with pgcrypto
	FirstName presence.Of[string]    `db:"first_name"`
	LastName  presence.Of[string]    `db:"last_name"`
	Role      Role                   `db:"role"`
	CreatedAt time.Time              `db:"created_at"`
	UpdatedAt presence.Of[time.Time] `db:"updated_at"`
	DeletedAt presence.Of[time.Time] `db:"deleted_at"` // soft delete
}

// NewUser creates a new User with required fields.
// ID and CreatedAt are set by the database.
func NewUser(email, password string, role Role) *User {
	return &User{
		Email:    email,
		Password: password,
		Role:     role,
	}
}

// Validate checks all required fields and formats.
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrEmailRequired
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return ErrEmailInvalid
	}

	if u.Password == "" {
		return ErrPasswordRequired
	}

	if len(u.Password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	if !u.Role.IsValid() {
		return ErrRoleInvalid
	}

	return nil
}

// SetFirstName sets the first name.
func (u *User) SetFirstName(name string) {
	u.FirstName = presence.FromValue(name)
}

// SetLastName sets the last name.
func (u *User) SetLastName(name string) {
	u.LastName = presence.FromValue(name)
}

// MarkDeleted marks the user as soft-deleted.
func (u *User) MarkDeleted() {
	u.DeletedAt = presence.FromValue(time.Now())
}

// IsDeleted returns true if the user is soft-deleted.
func (u *User) IsDeleted() bool {
	return u.DeletedAt.IsSet() && !u.DeletedAt.IsNull()
}
