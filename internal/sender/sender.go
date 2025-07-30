package sender

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"github.com/example/rabbitprobe/internal/conn"
	"github.com/example/rabbitprobe/internal/metrics"
)

// Sender publishes random messages.
type Sender struct {
	manager *conn.Manager
	ex      string
	rk      string
	size    int
	count   int
	rate    rate.Limiter
	log     *logrus.Logger
}

// Params defines send parameters.
type Params struct {
	Exchange string
	Routing  string
	Size     int
	Count    int
	Rate     int // per second
}

// New returns a sender.
func New(mgr *conn.Manager, log *logrus.Logger) *Sender {
	return &Sender{manager: mgr, log: log}
}

// Run performs sending.
func (s *Sender) Run(ctx context.Context, p Params) error {
	s.ex = p.Exchange
	s.rk = p.Routing
	s.size = p.Size
	s.count = p.Count
	if p.Rate > 0 {
		s.rate = *rate.NewLimiter(rate.Limit(p.Rate), p.Rate)
	}

	ch, err := s.manager.Channel()
	if err != nil {
		return err
	}
	ch.Confirm(false)
	confirm := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	buf := make([]byte, s.size)
	for i := 0; i < s.count; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if s.rate.Limit() > 0 {
			if err := s.rate.Wait(ctx); err != nil {
				return err
			}
		}
		rand.Read(buf)
		start := time.Now()
		if err := ch.PublishWithContext(ctx, s.ex, s.rk, false, false, amqp.Publishing{Body: buf}); err != nil {
			return err
		}
		c := <-confirm
		if !c.Ack {
			return amqp.ErrClosed
		}
		rtt := time.Since(start)
		metrics.Throughput.Inc()
		metrics.RTT.Observe(float64(rtt.Milliseconds()))
	}
	return nil
}
