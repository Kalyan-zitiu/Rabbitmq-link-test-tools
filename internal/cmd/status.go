package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/example/rabbitprobe/internal/metrics"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "show connection status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("connected: %.0f\n", metrics.ConnectedValue())
		},
	}
}
