package sender

import (
	"context"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"go.uber.org/zap"
)

type loggingSaver struct {
	log  logger.Logger
	next Saver
}

func NewLoggingSaver(log logger.Logger, next Saver) Saver {
	return &loggingSaver{
		log:  log,
		next: next,
	}
}

func (s *loggingSaver) Save(ctx context.Context, eventID string, body []byte, date time.Time) error {
	s.log.Info("Message from broker:", zap.Any("event", string(body)))
	if s.next == nil {
		return nil
	}
	if err := s.next.Save(ctx, eventID, body, date); err != nil {
		return err
	}
	return nil
}
