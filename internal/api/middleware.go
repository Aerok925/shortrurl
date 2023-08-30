package api

import (
	"fmt"
	"net/http"
	"time"
)

type loggingMiddleware struct {
	r          http.ResponseWriter
	statusCode int
	bufferSize int
}

func (l *loggingMiddleware) Header() http.Header {
	return l.r.Header()
}

func (l *loggingMiddleware) Write(data []byte) (int, error) {
	l.bufferSize = len(data)
	return l.r.Write(data)
}

func (l *loggingMiddleware) WriteHeader(statusCode int) {
	l.statusCode = statusCode
	l.r.WriteHeader(statusCode)
}

func (api *API) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var l loggingMiddleware
		timeStart := time.Now()
		l.r = w
		next.ServeHTTP(&l, r)
		api.logger.Debug(fmt.Sprintf("uri - %s, method - %s, timeDuration - %d, statusCode - %d, bufferSize - %d", r.RequestURI, r.Method, time.Since(timeStart).Milliseconds(), l.statusCode, l.bufferSize))
	})
}
