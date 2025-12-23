package cmd

import (
	"fmt"

	"github.com/pivaldi/go-cleanstack/internal/app/adapters"
	"github.com/pivaldi/go-cleanstack/internal/app/service"
	"github.com/pivaldi/go-cleanstack/internal/infra/api"
	"github.com/pivaldi/go-cleanstack/internal/infra/persistence"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(_ *cobra.Command, _ []string) error {
			db, err := persistence.NewDB(cfg.Database.URL)
			if err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer db.Close()

			logger.Info("connected to database")

			infraRepo := persistence.NewItemRepo(db, logger)
			itemRepo := adapters.NewItemRepositoryAdapter(infraRepo)
			itemService := service.NewItemService(itemRepo, logger)

			server := api.NewServer(cfg.Server.Port, itemService, logger)

			return server.Start()
		},
	}
}
