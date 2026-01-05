package config

import "github.com/pivaldi/go-cleanstack/internal/infra/config"

var cfg *config.Config

func SetConfig(c *config.Config) {
	cfg = c
}

func GetConfig() *config.Config {
	return cfg
}
