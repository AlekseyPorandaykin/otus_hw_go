package sender

import (
	"context"
	"encoding/json"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue"
	"go.uber.org/zap"
)

type Sender struct {
	consumer queue.Consumer
	log      logger.Logger
}

func NewSender(consumer queue.Consumer, log logger.Logger) *Sender {
	return &Sender{
		consumer: consumer,
		log:      log,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	return s.consumer.Listen(ctx, s)
}

func (s *Sender) Handle(d *queue.Message) error {
	var event calendar.EventDto
	if err := json.Unmarshal(d.Body, &event); err != nil {
		return err
	}
	s.log.Info("Message from broker:", zap.Any("event", event))
	return nil
}
