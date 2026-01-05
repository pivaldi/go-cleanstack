package cmd

import (
	"fmt"

	"github.com/pivaldi/go-cleanstack/internal/app/adapters"
	appConfig "github.com/pivaldi/go-cleanstack/internal/app/config"
	"github.com/pivaldi/go-cleanstack/internal/app/service"
	"github.com/pivaldi/go-cleanstack/internal/infra/api"
	"github.com/pivaldi/go-cleanstack/internal/infra/persistence"
	"github.com/pivaldi/go-cleanstack/internal/platform/logging"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := appConfig.GetConfig()

			db, err := persistence.NewDB(cfg.Database.URL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			logging.GetLogger().Info("connected to database")

			infraRepo := persistence.NewItemRepo(db)
			itemRepo := adapters.NewItemRepositoryAdapter(infraRepo)
			itemService := service.NewItemService(itemRepo)

			server := api.NewServer(cfg.Server.Port, itemService)

			return server.Start()
		},
	}
}
