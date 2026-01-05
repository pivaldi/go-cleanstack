package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/pivaldi/go-cleanstack/internal/app/app1/config"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logger/zap"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logging"
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
			cfg := config.Get()
			if cfg == nil {
				var err error
				cfg, err = config.Load(configPath)
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				// CLI flag overrides config
				if logLevel != "" {
					cfg.Platform.Log.Level = logLevel
				}

				config.SetConfig(cfg)
			}

			setup()

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "",
		"log level (debug, info, warn, error) - overrides config file")
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "",
		"The path where live the configuration files config_default.toml and config_"+os.Getenv("APP_ENV")+".toml")

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewServeCmd())
	// app.cmd.AddCommand(NewMigrateCmd())

	return rootCmd
}

func setup() {
	cfg := config.Get()
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

	logging.SetLogger(logger)
}
