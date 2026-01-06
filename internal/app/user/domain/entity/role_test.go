package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRole_Values(t *testing.T) {
	assert.Equal(t, Role("admin"), RoleAdmin)
	assert.Equal(t, Role("user"), RoleUser)
}

func TestRole_ParseRole(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Role
		wantErr bool
	}{
		{"valid admin", "admin", RoleAdmin, false},
		{"valid user", "user", RoleUser, false},
		{"invalid role", "invalid", Role(""), true},
		{"empty string", "", Role(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRole(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRole_String(t *testing.T) {
	assert.Equal(t, "admin", RoleAdmin.String())
	assert.Equal(t, "user", RoleUser.String())
}

func TestRole_JSONMarshal(t *testing.T) {
	data, err := json.Marshal(RoleAdmin)
	require.NoError(t, err)
	assert.Equal(t, `"admin"`, string(data))
}

func TestRole_JSONUnmarshal(t *testing.T) {
	var role Role
	err := json.Unmarshal([]byte(`"user"`), &role)
	require.NoError(t, err)
	assert.Equal(t, RoleUser, role)
}
