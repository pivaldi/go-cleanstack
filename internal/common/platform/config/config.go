package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type AppEnv string

type Platform struct {
	AppEnv   AppEnv
	Server   serverConfig
	Database databaseConfig
	Log      logConfig
}

func (p *Platform) SetAppEnv(appEnv AppEnv) {
	if p == nil {
		return
	}

	p.AppEnv = appEnv
}

type Config struct {
	Platform
}

func (c *Config) SetAppEnv(appEnv AppEnv) {
	if c == nil {
		return
	}

	c.Platform.AppEnv = appEnv
}

type configI interface {
	SetAppEnv(AppEnv)
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

func Load[T configI](configPath string) (T, error) {
	cfg := *new(T)
	env := os.Getenv("APP_ENV")
	if env == "" {
		return cfg, errors.New("APP_ENV environment variable is not set")
	}

	viper.SetConfigName("config_default")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configPath)
	_ = viper.ReadInConfig()

	viper.SetConfigName("config_" + env)

	if err := viper.MergeInConfig(); err != nil {
		return cfg, fmt.Errorf("failed to merge config file: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.SetAppEnv(AppEnv(env))

	return cfg, nil
}
