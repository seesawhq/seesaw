package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqConextWithRequestID := context.WithValue(r.Context(), "request_id", "kjkjkj")
		starttime := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r.WithContext(reqConextWithRequestID))
		slog.Info("", "method", r.Method, "path", r.URL.Path, "http_status", rw.status, "time", fmt.Sprintf("%v", time.Since(starttime)))
	})
}
