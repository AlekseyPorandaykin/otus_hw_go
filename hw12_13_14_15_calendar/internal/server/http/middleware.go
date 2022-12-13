package internalhttp

import (
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (s *StatusResponseWriter) WriteHeader(status int) {
	s.Status = status
	s.ResponseWriter.WriteHeader(status)
}

func (h *Handler) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &StatusResponseWriter{ResponseWriter: w}
		next(sw, r)
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		h.logger.Info(
			"Request",
			zap.String("ip", ip),
			zap.String("datetime", time.Now().String()),
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.String("user-agent", r.UserAgent()),
			zap.Any("latency", time.Since(start).String()),
			zap.String("protocol", r.Proto),
			zap.Any("http-status", sw.Status),
		)
	}
}
