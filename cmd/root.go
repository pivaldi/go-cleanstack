package cmd

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/pivaldi/go-cleanstack/internal/infra/config"
	"github.com/pivaldi/go-cleanstack/internal/platform/logging"
	"github.com/spf13/cobra"
)

var (
	cfg      *config.Config
	logger   *zap.Logger
	logLevel string
)

func NewRootCmd() *cobra.Command {
	var (
		configPath string
	)

	rootCmd := &cobra.Command{
		Use:   "cleanstack",
		Short: "GoCleanstack application",
		Long:  "A production-ready Go application with CLI and API",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			var err error
			cfg, err = config.Load(".")
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// CLI flag overrides config
			effectiveLogLevel := cfg.Log.Level
			if logLevel != "" {
				effectiveLogLevel = logLevel
			}

			logger, err = logging.NewLogger(cfg.App.Env, effectiveLogLevel)
			if err != nil {
				return fmt.Errorf("failed to initialize logger: %w", err)
			}

			logger.Info("application starting",
				zap.String("env", cfg.App.Env),
				zap.String("log_level", effectiveLogLevel),
			)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "",
		"log level (debug, info, warn, error) - overrides config file")
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "",
		"The path where live the configuration files config_default.toml and config_"+os.Getenv("APP_ENV")+".toml")

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewServeCmd())
	rootCmd.AddCommand(NewMigrateCmd())

	return rootCmd
}
