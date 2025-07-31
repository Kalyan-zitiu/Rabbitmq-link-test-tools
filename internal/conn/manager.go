package conn

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/example/rabbitprobe/internal/metrics"
)

// Manager handles connection and reconnection.
type Manager struct {
	addrs     []string
	vhost     string
	mu        sync.Mutex
	conn      *amqp.Connection
	log       *logrus.Logger
	downAt    time.Time
	connected bool
}

// New creates a new Manager.
func New(addrs []string, vhost string, log *logrus.Logger) *Manager {
	return &Manager{addrs: addrs, vhost: vhost, log: log}
}

func (m *Manager) dial(addr string) (*amqp.Connection, error) {
	cfg := amqp.Config{
		Vhost:     m.vhost,
		Heartbeat: time.Second,
		Locale:    "en_US",
		Dial:      (&net.Dialer{Timeout: time.Second, KeepAlive: 30 * time.Second}).DialContext,
	}
	return amqp.DialConfig(addr, cfg)
}

// Connect establishes connection with retry.
func (m *Manager) Connect(ctx context.Context) {
	backoff := 100 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		addr := m.addrs[rand.Intn(len(m.addrs))]
		conn, err := m.dial(addr)
		if err != nil {
			m.log.WithError(err).Error("connect failed")
			time.Sleep(backoff)
			if backoff < 5*time.Second {
				backoff *= 2
			}
			continue
		}

		m.mu.Lock()
		m.conn = conn
		m.connected = true
		m.mu.Unlock()
		metrics.SetConnected(true)
		m.log.WithField("addr", addr).Info("ConnUp")
		if !m.downAt.IsZero() {
			downtime := time.Since(m.downAt)
			metrics.Downtime.Observe(downtime.Seconds())
			m.log.WithField("downtime_ms", downtime.Milliseconds()).Info("ConnRecovered")
			m.downAt = time.Time{}
		}

		m.waitClose(conn)
		metrics.SetConnected(false)
		metrics.DisconnectTotal.Inc()
		m.mu.Lock()
		m.conn = nil
		m.connected = false
		m.mu.Unlock()
		m.downAt = time.Now()
	}
}

func (m *Manager) waitClose(conn *amqp.Connection) {
	errCh := make(chan *amqp.Error, 1)
	conn.NotifyClose(errCh)
	err := <-errCh
	if err != nil {
		m.log.WithError(err).Error("ConnDown")
	} else {
		m.log.Info("Conn closed")
	}
}

// Channel returns a new channel using current connection.
func (m *Manager) Channel() (*amqp.Channel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conn == nil {
		return nil, amqp.ErrClosed
	}
	return m.conn.Channel()
}

// Connected returns connection state.
func (m *Manager) Connected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.connected
}
