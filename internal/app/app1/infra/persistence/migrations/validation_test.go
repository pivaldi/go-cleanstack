package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid description",
			description: "add user table",
			wantErr:     false,
		},
		{
			name:        "valid with underscores",
			description: "add_user_authentication",
			wantErr:     false,
		},
		{
			name:        "empty description",
			description: "",
			wantErr:     true,
			errContains: "description is required",
		},
		{
			name:        "too short",
			description: "ab",
			wantErr:     true,
			errContains: "description must be between 3 and 100 characters",
		},
		{
			name:        "too long",
			description: "this is a very long description that exceeds the maximum allowed length of one hundred characters for migration descriptions",
			wantErr:     true,
			errContains: "description must be between 3 and 100 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.description)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "dash separated to camel case",
			input: "add-user-table",
			want:  "AddUserTable",
		},
		{
			name:  "underscore separated to camel case",
			input: "add_user_table",
			want:  "AddUserTable",
		},
		{
			name:  "spaces to camel case",
			input: "add user table",
			want:  "AddUserTable",
		},
		{
			name:  "single word",
			input: "users",
			want:  "Users",
		},
		{
			name:  "mixed separators",
			input: "add-user_authentication-table",
			want:  "AddUserAuthenticationTable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
