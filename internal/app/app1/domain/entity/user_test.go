package entity

import (
	"testing"
	"time"

	"github.com/pivaldi/presence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	user := NewUser("test@example.com", "password123", RoleUser)

	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "password123", user.Password)
	assert.Equal(t, RoleUser, user.Role)
	assert.False(t, user.FirstName.IsSet())
	assert.False(t, user.LastName.IsSet())
}

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    *User
		wantErr error
	}{
		{
			name:    "valid user",
			user:    NewUser("test@example.com", "password123", RoleUser),
			wantErr: nil,
		},
		{
			name:    "valid admin",
			user:    NewUser("admin@example.com", "adminpass123", RoleAdmin),
			wantErr: nil,
		},
		{
			name:    "empty email",
			user:    NewUser("", "password123", RoleUser),
			wantErr: ErrEmailRequired,
		},
		{
			name:    "invalid email format",
			user:    NewUser("notanemail", "password123", RoleUser),
			wantErr: ErrEmailInvalid,
		},
		{
			name:    "empty password",
			user:    NewUser("test@example.com", "", RoleUser),
			wantErr: ErrPasswordRequired,
		},
		{
			name:    "password too short",
			user:    NewUser("test@example.com", "short", RoleUser),
			wantErr: ErrPasswordTooShort,
		},
		{
			name:    "invalid role",
			user:    NewUser("test@example.com", "password123", Role("invalid")),
			wantErr: ErrRoleInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUser_SetFirstName(t *testing.T) {
	user := NewUser("test@example.com", "password123", RoleUser)
	assert.False(t, user.FirstName.IsSet())

	user.SetFirstName("John")

	assert.True(t, user.FirstName.IsSet())
	val, err := user.FirstName.Value()
	require.NoError(t, err)
	assert.Equal(t, "John", val)
}

func TestUser_SetLastName(t *testing.T) {
	user := NewUser("test@example.com", "password123", RoleUser)
	assert.False(t, user.LastName.IsSet())

	user.SetLastName("Doe")

	assert.True(t, user.LastName.IsSet())
	val, err := user.LastName.Value()
	require.NoError(t, err)
	assert.Equal(t, "Doe", val)
}

func TestUser_MarkDeleted(t *testing.T) {
	user := NewUser("test@example.com", "password123", RoleUser)
	assert.False(t, user.IsDeleted())

	user.MarkDeleted()

	assert.True(t, user.IsDeleted())
	assert.True(t, user.DeletedAt.IsSet())
}

func TestUser_IsDeleted(t *testing.T) {
	t.Run("not deleted", func(t *testing.T) {
		user := NewUser("test@example.com", "password123", RoleUser)
		assert.False(t, user.IsDeleted())
	})

	t.Run("deleted", func(t *testing.T) {
		user := NewUser("test@example.com", "password123", RoleUser)
		user.MarkDeleted()
		assert.True(t, user.IsDeleted())
	})

	t.Run("explicitly null", func(t *testing.T) {
		user := NewUser("test@example.com", "password123", RoleUser)
		user.DeletedAt = presence.Null[time.Time]()
		assert.False(t, user.IsDeleted())
	})
}
