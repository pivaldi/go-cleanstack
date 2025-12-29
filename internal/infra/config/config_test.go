package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultsToDevEnvironment(t *testing.T) {
	// Setup: create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config_development.toml")
	err := os.WriteFile(configPath, []byte(`
[server]
port = 8080

[database]
url = "postgres://localhost/test"

[log]
level = "debug"
`), 0644)
	require.NoError(t, err)

	// Set APP_ENV to development
	os.Setenv("APP_ENV", "development")
	t.Cleanup(func() { os.Unsetenv("APP_ENV") })

	cfg, err := Load(tmpDir)
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "postgres://localhost/test", cfg.Database.URL)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "development", string(cfg.AppEnv))
}
