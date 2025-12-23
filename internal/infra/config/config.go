package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var cfg *Config

type Config struct {
	App      appConfig
	Server   serverConfig
	Database databaseConfig
	Log      logConfig
}

type appConfig struct {
	Env string
}

type serverConfig struct {
	Port int
}

type databaseConfig struct {
	URL string
}

type logConfig struct {
	Level string
}

func Load(configPath string) (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		return nil, errors.New("APP_ENV environment variable is not set")
	}

	viper.SetConfigName("config_default")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configPath)
	_ = viper.ReadInConfig()

	viper.SetConfigName("config_" + env)

	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.App.Env = env

	return cfg, nil
}

func GetConfig() (*Config, error) {
	if cfg == nil {
		return nil, errors.New("config is not loaded")
	}

	return cfg, nil
}
