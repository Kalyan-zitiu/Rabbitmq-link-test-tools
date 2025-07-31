package probe

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/example/rabbitprobe/internal/conn"
	"github.com/example/rabbitprobe/internal/metrics"
)

// Engine performs periodic publish-confirm probes.
type Engine struct {
	manager  *conn.Manager
	exchange string
	routing  string
	interval time.Duration
	log      *logrus.Logger
}

// New creates a probe engine.
func New(mgr *conn.Manager, ex, rk string, interval time.Duration, log *logrus.Logger) *Engine {
	return &Engine{manager: mgr, exchange: ex, routing: rk, interval: interval, log: log}
}

// Start starts the probe loop.
func (e *Engine) Start(ctx context.Context) {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.ping(ctx)
		}
	}
}

func (e *Engine) ping(ctx context.Context) {
	ch, err := e.manager.Channel()
	if err != nil {
		return
	}
	ch.Confirm(false)

	start := time.Now()
	body := fmt.Sprintf("ping:%d", start.UnixNano())
	if err := ch.PublishWithContext(ctx, e.exchange, e.routing, false, false, amqp.Publishing{Body: []byte(body)}); err != nil {
		e.log.WithError(err).Error("probe publish failed")
		return
	}
	confirm := make(chan amqp.Confirmation, 1)
	ch.NotifyPublish(confirm)
	select {
	case c := <-confirm:
		if !c.Ack {
			e.log.Error("probe not acked")
			return
		}
		rtt := time.Since(start)
		metrics.RTT.Observe(float64(rtt.Milliseconds()))
		e.log.WithField("latency_ms", rtt.Milliseconds()).Info("probe success")
	case <-time.After(e.interval * 2):
		e.log.Error("probe timeout")
		_ = ch.Close()
	}
}
