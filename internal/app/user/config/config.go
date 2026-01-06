package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/pivaldi/go-cleanstack/internal/common/platform/config"
	"github.com/spf13/viper"
)

type Config struct {
	Platform config.Platform
}

var cfg *Config

func SetConfig(c *Config) {
	cfg = c
}

func MustGet() *Config {
	if cfg == nil {
		panic(errors.New("config is not loaded. Use config.Load() to initialize"))
	}
	return cfg
}

func Get() *Config {
	return cfg
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

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error loading default configuration : %w", err)
	}

	viper.SetConfigName("config_" + env)

	if err := viper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.Platform.AppEnv = config.AppEnv(env)

	return cfg, nil
}
