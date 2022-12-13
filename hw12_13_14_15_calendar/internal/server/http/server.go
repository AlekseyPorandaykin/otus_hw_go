package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"go.uber.org/zap"
)

type Server struct {
	logger logger.Logger
	server *http.Server
}

var closeTimeout = time.Second * 3

type EventDto struct {
	Title         string
	DateTimeStart time.Time
	DateTimeEnd   time.Time
	Description   string
	CreatedBy     int
	RemindFrom    time.Time
}

type Application interface {
	CreateEvent(ctx context.Context, event EventDto) error
	ReadEvent(ctx context.Context, id string) (*calendar.Event, error)
	UpdateEvent(ctx context.Context, event EventDto) error
	DeleteEvent(ctx context.Context, id string) error
}

type Config struct {
	Host              string        `mapstructure:"host"`
	Port              string        `mapstructure:"port"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
}

func NewServer(logger logger.Logger, app Application, conf *Config) *Server {
	return &Server{
		server: &http.Server{
			Addr:              net.JoinHostPort(conf.Host, conf.Port),
			Handler:           NewHandler(app, logger),
			ReadTimeout:       conf.ReadTimeout,
			WriteTimeout:      conf.WriteTimeout,
			ReadHeaderTimeout: conf.ReadHeaderTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Start http server", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer func() {
		s.logger.Debug("Shutdown http server")
		cancel()
	}()

	if err := s.server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to stop http server", zap.Error(err))
	}
}
