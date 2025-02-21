package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

type LoggingMiddleware struct {
	h http.Handler
}

func NewLoggingMiddleware(h http.Handler) http.Handler {
	return &LoggingMiddleware{h: h}
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Print("request url = ", r.URL.Path, " body = ", string(body))

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	m.h.ServeHTTP(w, r)
}
