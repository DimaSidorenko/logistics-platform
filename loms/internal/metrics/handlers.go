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
			ConstLabels: prometheus.Labels{"service": "loms"},
		},
		[]string{"handler", "code"},
	)

	handlerHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   "app",
			Name:        "handler_request_duration_histogram",
			Help:        "Total duration of processing request",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{"service": "loms"},
		},
		[]string{"handler"},
	)

	analyzeFileContentHistogram = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace:   "app",
			Name:        "analyzer_filecontent_histogram",
			Help:        "Total duration of processing text",
			Buckets:     []float64{.5, 1, 5, 10, 30, 60},
			ConstLabels: prometheus.Labels{"service": "loms"},
		},
	)
)

func RequestCounterInc(handler string, code string) {
	requestCounter.WithLabelValues(handler, code).Inc()
}

func RequestHandlerDuration(handler string, since time.Duration) {
	handlerHistogram.WithLabelValues(handler).Observe(since.Seconds())
}

func AnalyzeFileContentDuration(since time.Duration) {
	analyzeFileContentHistogram.Observe(since.Seconds())
}
