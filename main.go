package main

import (
	"errors"
	"os"

	"github.com/pivaldi/go-cleanstack/cmd"

	"github.com/pivaldi/go-cleanstack/internal/infra/config"
	"github.com/pivaldi/go-cleanstack/internal/platform/clierr"
)

var (
	appEnvName = "APP_ENV"
	env        = os.Getenv(appEnvName)
)

func prerequisitesTest() error {
	if env == "" {
		return errors.New(appEnvName + " environment variable is not set")
	}

	return nil
}

func main() {
	if err := prerequisitesTest(); err != nil {
		clierr.ExitOnError(err, true)
	}

	if err := cmd.NewRootCmd().Execute(); err != nil {
		cfg, errr := config.GetConfig()
		if errr != nil {
			clierr.ExitOnError(errr, true)
		}

		debug := cfg.Log.Level == "debug"
		clierr.ExitOnError(err, debug)
	}
}
