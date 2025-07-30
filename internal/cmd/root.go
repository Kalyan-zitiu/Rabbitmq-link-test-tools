package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/example/rabbitprobe/internal/conn"
	"github.com/example/rabbitprobe/internal/logger"
	"github.com/example/rabbitprobe/internal/metrics"
)

var (
	addrs       string
	vhost       string
	logFile     string
	metricsPort int
)

// rootCmd is the base command.
var rootCmd = &cobra.Command{
	Use:   "rabbitprobe",
	Short: "RabbitMQ link probe tool",
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&addrs, "addrs", "amqp://guest:guest@localhost:5672/", "amqp addresses comma separated")
	rootCmd.PersistentFlags().StringVar(&vhost, "vhost", "/", "vhost")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "log file path")
	rootCmd.PersistentFlags().IntVar(&metricsPort, "metrics-port", 2112, "metrics port")

	rootCmd.AddCommand(newProbeCmd())
	rootCmd.AddCommand(newSendCmd())
	rootCmd.AddCommand(newStatusCmd())
}

func setup() (*conn.Manager, *logrus.Logger, context.Context, context.CancelFunc) {
	log := logger.New(logFile, logrus.InfoLevel)
	mgr := conn.New(strings.Split(addrs, ","), vhost, log)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		addr := fmt.Sprintf(":%d", metricsPort)
		if err := metrics.Start(addr); err != nil {
			log.WithError(err).Error("metrics failed")
		}
	}()

	go mgr.Connect(ctx)
	return mgr, log, ctx, cancel
}
