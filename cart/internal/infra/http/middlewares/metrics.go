package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"route256/cart/internal/metrics"
)

// ResponseWriterWrapper оборачивает http.ResponseWriter и захватывает статус код.
type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

type MetricsMiddleware struct {
	h http.Handler
}

func NewMetricsMiddleware(h http.Handler) http.Handler {
	return &MetricsMiddleware{h: h}
}

func (m *MetricsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Создаем обертку для ResponseWriter, чтобы захватить статус код.
	wrapper := &ResponseWriterWrapper{ResponseWriter: w}

	route := mux.CurrentRoute(r)

	urlTemplate, err := route.GetPathTemplate()
	if err != nil {
		urlTemplate = "unknown_url"
	}

	start := time.Now()
	m.h.ServeHTTP(wrapper, r)
	duration := time.Since(start)

	handler := r.Method + " " + urlTemplate

	metrics.RequestCounterInc(handler, strconv.Itoa(wrapper.StatusCode))
	metrics.RequestHandlerDuration(handler, duration)
}
