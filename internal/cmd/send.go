package cmd

import (
	"github.com/spf13/cobra"

	"github.com/example/rabbitprobe/internal/sender"
)

var (
	sendEx    string
	sendRK    string
	sendSize  int
	sendCount int
	sendRate  int
)

func newSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "send random payloads",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, log, ctx, cancel := setup()
			defer cancel()
			s := sender.New(mgr, log)
			p := sender.Params{Exchange: sendEx, Routing: sendRK, Size: sendSize, Count: sendCount, Rate: sendRate}
			return s.Run(ctx, p)
		},
	}

	cmd.Flags().StringVar(&sendEx, "ex", "", "exchange")
	cmd.Flags().StringVar(&sendRK, "rk", "", "routing key")
	cmd.Flags().IntVar(&sendSize, "size", 1024, "message size bytes")
	cmd.Flags().IntVar(&sendCount, "count", 1, "message count")
	cmd.Flags().IntVar(&sendRate, "rate", 0, "rate per second")
	return cmd
}
