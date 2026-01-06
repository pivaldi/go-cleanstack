package main

import (
	"errors"
	"fmt"
	"os"

	userCmd "github.com/pivaldi/go-cleanstack/internal/app/user/cmd"
	userConfig "github.com/pivaldi/go-cleanstack/internal/app/user/config"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/clierr"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/config"
	"github.com/spf13/cobra"
)

var (
	appEnvName = "APP_ENV"
	env        = os.Getenv(appEnvName)
	cfg        *config.Config
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

	rootCmd := newRootCmd()
	rootCmd.AddCommand(userCmd.GetRootCmd())

	if err := rootCmd.Execute(); err != nil {
		clierr.ExitOnError(err, true)
	}
}

func newRootCmd() *cobra.Command {
	var (
		configPath = "."
		logLevel   string
	)

	rootCmd := &cobra.Command{
		Use:   "cleanstack",
		Short: "GoCleanstack application",
		Long:  "A production-ready Go application with CLI and API",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			var err error
			cfg, err = config.Load[*config.Config](configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// CLI flag overrides config
			if logLevel != "" {
				cfg.Log.Level = logLevel
			}

			theUserConfig, ok := any(cfg).(userConfig.Config)
			if !ok {
				return fmt.Errorf("can not cast %T to %T", cfg, userConfig.Config{})
			}

			userConfig.SetConfig(&theUserConfig)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "",
		"log level (debug, info, warn, error) - overrides config file")
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "",
		"The path where live the configuration files config_default.toml and config_"+os.Getenv("APP_ENV")+".toml")

	return rootCmd
}
