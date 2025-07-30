package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/example/rabbitprobe/internal/probe"
)

var (
	probeEx       string
	probeRK       string
	probeInterval time.Duration
	probeCtx      context.Context
	probeCancel   context.CancelFunc
)

func newProbeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "probe",
		Short: "start or stop probe",
	}

	start := &cobra.Command{
		Use:   "start",
		Short: "start probe",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, log, ctx, cancel := setup()
			probeCtx = ctx
			probeCancel = cancel
			eng := probe.New(mgr, probeEx, probeRK, probeInterval, log)
			go eng.Start(ctx)
			<-ctx.Done()
			return nil
		},
	}
	start.Flags().StringVar(&probeEx, "ex", "", "exchange")
	start.Flags().StringVar(&probeRK, "rk", "", "routing key")
	start.Flags().DurationVar(&probeInterval, "interval", 200*time.Millisecond, "probe interval")

	stop := &cobra.Command{
		Use:   "stop",
		Short: "stop probe",
		Run: func(cmd *cobra.Command, args []string) {
			if probeCancel != nil {
				probeCancel()
			}
		},
	}
	cmd.AddCommand(start, stop)
	return cmd
}
