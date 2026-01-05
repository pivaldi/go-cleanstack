package cmd

import (
	"fmt"

	appConfig "github.com/pivaldi/go-cleanstack/internal/app/config"
	"github.com/spf13/cobra"
)

var (
	BuildTime = "unknown"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("go-cleanstack version %s (built %s)\n", appConfig.GetConfig().AppEnv, BuildTime)
		},
	}
}
