package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/pivaldi/go-cleanstack/internal/app/cmd"
	appConfig "github.com/pivaldi/go-cleanstack/internal/app/config"
	"github.com/pivaldi/go-cleanstack/internal/infra/config"
	"github.com/pivaldi/go-cleanstack/internal/platform/clierr"
	"github.com/pivaldi/go-cleanstack/internal/platform/logger/zap"
	"github.com/pivaldi/go-cleanstack/internal/platform/logging"
	"github.com/spf13/cobra"
)

func GetRootCmd() *cobra.Command {
	var (
		configPath string = "."
		logLevel   string
	)

	rootCmd := &cobra.Command{
		Use:   "app",
		Short: "Application",
		Long:  "Application",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			cfg := appConfig.GetConfig()
			if cfg == nil {
				var err error
				cfg, err = config.Load(configPath)
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				// CLI flag overrides config
				if logLevel != "" {
					cfg.Log.Level = logLevel
				}

				appConfig.SetConfig(cfg)
			}

			setup()

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "",
		"log level (debug, info, warn, error) - overrides config file")
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "",
		"The path where live the configuration files config_default.toml and config_"+os.Getenv("APP_ENV")+".toml")

	rootCmd.AddCommand(cmd.NewVersionCmd())
	rootCmd.AddCommand(cmd.NewServeCmd())
	// app.cmd.AddCommand(NewMigrateCmd())

	return rootCmd
}

func setup() {
	cfg := appConfig.GetConfig()
	logger, err := zap.NewDevelopment(cfg.Log.Level)
	if err != nil {
		panic(err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		panic(errors.New("APP_ENV environment variable is not set"))
	}

	logger.Info("application starting",
		logging.String("env", env),
		logging.String("log_level", cfg.Log.Level),
	)

	logging.SetLogger(logger)
}

func main() {
	if err := GetRootCmd().Execute(); err != nil {
		clierr.ExitOnError(err, true)
	}

}
