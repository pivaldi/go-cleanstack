package cmd

import (
	"fmt"
	"os"

	appConfig "github.com/pivaldi/go-cleanstack/internal/app/user/config"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/config"
	"github.com/spf13/cobra"
)

func GetRootCmd() *cobra.Command {
	var (
		configDir string = ""
		logLevel  string
	)

	rootCmd := &cobra.Command{
		Use:   "user",
		Short: "User Application",
		Long:  "User Application",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			cfg := appConfig.Get()
			if cfg == nil {
				var err error
				cfg = &appConfig.Config{}

				err = config.Load(configDir, cfg)
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				// CLI flag overrides config
				if logLevel != "" {
					cfg.Platform.Log.Level = logLevel
				}

			}

			appConfig.Setup(cfg)

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "",
		"log level (debug, info, warn, error) - overrides config file")
	rootCmd.PersistentFlags().StringVar(&configDir, "config-dir", "",
		"The directory where lives the configuration files config_default.toml and config_"+os.Getenv("APP_ENV")+".toml")

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewServeCmd())
	// app.cmd.AddCommand(NewMigrateCmd())

	return rootCmd
}
