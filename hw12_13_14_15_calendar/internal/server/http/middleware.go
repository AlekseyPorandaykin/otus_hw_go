package internalhttp

import (
	"net"
	"net/http"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
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

type Middleware struct {
	logger logger.Logger
}

func (m *Middleware) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &StatusResponseWriter{ResponseWriter: w}
		next(sw, r)
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		m.logger.Info(
			"HTTP request",
			zap.String("ip", ip),
			zap.String("datetime", time.Now().String()),
			zap.String("method", r.Method),
			zap.String("path", r.RequestURI),
			zap.String("user-agent", r.UserAgent()),
			zap.Any("latency", time.Since(start).String()),
			zap.String("protocol", r.Proto),
			zap.Int("http-status", sw.Status),
		)
	}
}
