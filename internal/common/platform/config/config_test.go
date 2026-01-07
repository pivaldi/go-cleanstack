package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("APP_ENV", "development")
	t.Cleanup(func() { os.Unsetenv("APP_ENV") })

	cfg := &Config{}
	t.Run("Load default config", func(t *testing.T) {
		err := Load("/", cfg)
		require.NoError(t, err)
	})
	t.Run("Check config values", func(t *testing.T) {
		assert.NotEqual(t, 0, cfg.Server.Port)
		assert.Contains(t, cfg.Database.URL, "://")
		assert.Equal(t, "debug", cfg.Log.Level)
		assert.Equal(t, "development", string(cfg.AppEnv))
	})
}

func TestLoad_DefaultsToDevEnvironment(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config_development.toml")
	err := os.WriteFile(configPath, []byte(`
[platform.server]
port = 8080

[platform.database]
url = "postgres://localhost/test"

[platform.log]
level = "warn"
`), 0644)
	require.NoError(t, err)

	os.Setenv("APP_ENV", "development")
	t.Cleanup(func() { os.Unsetenv("APP_ENV") })

	cfg := &Config{}
	t.Run("Load dev config", func(t *testing.T) {
		err = Load(tmpDir, cfg)
		require.NoError(t, err)
	})

	t.Run("Check dev config values", func(t *testing.T) {
		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "postgres://localhost/test", cfg.Database.URL)
		assert.Equal(t, "warn", cfg.Log.Level)
		assert.Equal(t, "development", string(cfg.AppEnv))
	})
}
