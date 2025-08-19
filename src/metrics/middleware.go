package metrics

import (
	"net/http"
	"strconv"
	"time"
)

func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.ActiveConnections.Inc()
		defer m.ActiveConnections.Dec()

		recorder := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(recorder, r)

		duration := time.Since(start).Seconds()
		m.HttpDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		m.HttpRequests.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(recorder.status)).Inc()
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
