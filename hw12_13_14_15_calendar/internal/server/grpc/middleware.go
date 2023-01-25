package grpc

import (
	"context"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Middleware struct {
	logger logger.Logger
}

type Response interface {
	GetStatus() bool
}

func NewMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{logger: logger}
}

func (m *Middleware) loggingMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()
	errStr := ""

	resp, err = handler(ctx, req)

	if err != nil {
		errStr = err.Error()
	}

	m.logger.Info("GRPC request",
		zap.String("datetime", time.Now().String()),
		zap.Any("latency", time.Since(start).String()),
		zap.String("method", info.FullMethod),
		zap.Bool("status", getStatus(resp)),
		zap.String("error", errStr),
	)

	return
}

func getStatus(resp interface{}) bool {
	if v, ok := resp.(Response); ok {
		return v.GetStatus()
	}

	return false
}
