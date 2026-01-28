package config

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pivaldi/go-cleanstack/pkg/file"
	"github.com/spf13/viper"
)

//go:embed config_default.toml
var defaultConfig []byte

type AppEnv string

type Platform struct {
	AppEnv   AppEnv
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
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

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	URL string
}

type LogConfig struct {
	Level string
}

func Load[T configI](configDir string, dest T) error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return errors.New("APP_ENV environment variable is not set")
	}

	viper.SetConfigType("toml")
	viper.ReadConfig(bytes.NewBuffer(defaultConfig))

	baseFile := "config_" + env
	fileName := baseFile + ".toml"
	confPath := filepath.Join(configDir, fileName)
	if file.Exists(confPath) {
		viper.AddConfigPath(configDir)
		viper.SetConfigName(baseFile)

		if err := viper.MergeInConfig(); err != nil {
			return fmt.Errorf("failed to merge config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&dest); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	dest.SetAppEnv(AppEnv(env))

	return nil
}
