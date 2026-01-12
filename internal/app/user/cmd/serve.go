package cmd

import (
	"fmt"

	"github.com/pivaldi/go-cleanstack/internal/app/user/adapters"
	"github.com/pivaldi/go-cleanstack/internal/app/user/api"
	appConfig "github.com/pivaldi/go-cleanstack/internal/app/user/config"
	"github.com/pivaldi/go-cleanstack/internal/app/user/infra/persistence"
	"github.com/pivaldi/go-cleanstack/internal/app/user/service"
	"github.com/pivaldi/go-cleanstack/internal/common/platform/logger/zap"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := appConfig.Get()

			db, err := persistence.NewDB(cfg.Platform.Database.URL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			logger, err := zap.NewLogger(string(cfg.Platform.AppEnv), cfg.Platform.Log.Level)
			if err != nil {
				return fmt.Errorf("failed to create logger: %w", err)
			}

			logger.Info("connected to database")

			infraRepo := persistence.NewUserRepo(db)
			userRepo := adapters.NewUserRepositoryAdapter(infraRepo)
			userService := service.NewUserService(userRepo, logger)

			server := api.NewServer(cfg.Platform.Server.Port, userService, logger)

			return server.Start()
		},
	}
}
