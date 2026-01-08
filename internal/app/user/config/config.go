package config

import (
	"errors"
	"os"

	"github.com/pivaldi/go-cleanstack/internal/common/platform/config"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logger/zap"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logging"
)

type Config struct {
	Platform config.Platform
}

func (c *Config) SetAppEnv(appEnv config.AppEnv) {
	if c == nil {
		return
	}

	c.Platform.SetAppEnv(appEnv)
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

func Setup(cfg *Config) {
	if cfg == nil {
		panic(errors.New("config is not loaded. Use config.Load() to initialize"))
	}

	SetConfig(cfg)

	logger, err := zap.NewDevelopment(cfg.Platform.Log.Level)
	if err != nil {
		panic(err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		panic(errors.New("APP_ENV environment variable is not set"))
	}

	logger.Info("application starting",
		logging.String("env", env),
		logging.String("log_level", cfg.Platform.Log.Level),
	)

	// TODO: Implement setup logic in the app infrastructure not in common
	logging.SetLogger(logger)
}
