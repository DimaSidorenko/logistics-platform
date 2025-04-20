package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "app",
			Name:        "handler_request_total_counter",
			Help:        "Total amount of request by handler",
			ConstLabels: prometheus.Labels{"service": "cart"},
		},
		[]string{"handler", "code"},
	)

	handlerHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   "app",
			Name:        "handler_request_duration_histogram",
			Help:        "Total duration of processing request",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{"service": "cart"},
		},
		[]string{"handler"},
	)
)

func RequestCounterInc(handler string, code string) {
	requestCounter.WithLabelValues(handler, code).Inc()
}

func RequestHandlerDuration(handler string, since time.Duration) {
	handlerHistogram.WithLabelValues(handler).Observe(since.Seconds())
}
