package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	HttpRequests      *prometheus.CounterVec
	HttpDuration      *prometheus.HistogramVec
	ActiveConnections prometheus.Gauge
}

func NewMetrics() *Metrics {
	httpRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total"},
		[]string{"method", "endpoint", "status_code"},
	)

	httpDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	activeConnections := prometheus.NewGauge(
		prometheus.GaugeOpts{Name: "http_active_connections"},
	)

	prometheus.MustRegister(httpRequests, httpDuration, activeConnections)

	return &Metrics{
		HttpRequests:      httpRequests,
		HttpDuration:      httpDuration,
		ActiveConnections: activeConnections,
	}
}
