package sender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue"
)

type Saver interface {
	Save(ctx context.Context, eventID string, body []byte, date time.Time) error
}

type Sender struct {
	consumer queue.Consumer
	saver    Saver
}

func NewSender(consumer queue.Consumer, saver Saver) *Sender {
	return &Sender{
		consumer: consumer,
		saver:    saver,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	return s.consumer.Listen(ctx, s)
}

func (s *Sender) Handle(ctx context.Context, d *queue.Message) error {
	var event calendar.EventDto
	if err := json.Unmarshal(d.Body, &event); err != nil {
		return err
	}
	if err := s.saver.Save(ctx, event.ID, d.Body, time.Now()); err != nil {
		return err
	}
	return nil
}
