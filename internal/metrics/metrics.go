package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Connected = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "rabbitprobe_connected",
		Help: "Connection status 0/1",
	})
	connectedValue float64
	RTT            = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "rabbitprobe_rtt_ms",
		Help:    "Probe RTT in milliseconds",
		Buckets: prometheus.ExponentialBuckets(1, 2, 10),
	})
	DisconnectTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rabbitprobe_disconnect_total",
		Help: "Total disconnects",
	})
	Downtime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "rabbitprobe_downtime_seconds",
		Help:    "Downtime seconds",
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
	})
	Throughput = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "rabbitprobe_throughput_msgs_total",
		Help: "Total messages sent",
	})
)

// Start starts metrics HTTP server.
func Start(addr string) error {
	reg := prometheus.NewRegistry()
	reg.MustRegister(Connected, RTT, DisconnectTotal, Downtime, Throughput)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	return http.ListenAndServe(addr, nil)
}

// SetConnected sets gauge and value.
func SetConnected(v bool) {
	if v {
		Connected.Set(1)
		connectedValue = 1
	} else {
		Connected.Set(0)
		connectedValue = 0
	}
}

// ConnectedValue returns current value.
func ConnectedValue() float64 {
	return connectedValue
}
