package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	event "github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/server/grpc/pb/api/grpc"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Server struct {
	logger  logger.Logger
	host    string
	port    int
	handler *Handler
}

func NewServer(logger logger.Logger, app calendar.Application, conf *Config) *Server {
	return &Server{
		logger:  logger,
		host:    conf.Host,
		port:    conf.Port,
		handler: &Handler{app: app},
	}
}

func (s *Server) Listen(ctx context.Context) error {
	s.logger.Debug(fmt.Sprintf("Start grpc server %s:%d \n", s.host, s.port))
	errC := make(chan error)
	l, errL := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if errL != nil {
		return errL
	}
	h := NewMiddleware(s.logger)
	ser := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(h.loggingMiddleware)),
	)
	event.RegisterEventServiceServer(ser, s.handler)
	defer ser.Stop()
	go func() {
		if err := ser.Serve(l); err != nil {
			errC <- err
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errC:
			return err
		}
	}
}
