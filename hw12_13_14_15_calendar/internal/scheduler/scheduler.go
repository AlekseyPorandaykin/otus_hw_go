package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue"
	"go.uber.org/zap"
)

type Config struct {
	Duration time.Duration `mapstructure:"duration"`
}

type EventRepository interface {
	GetEventsWithRemindStatus(
		ctx context.Context, time time.Time, status calendar.RemindStatus,
	) ([]*calendar.Event, error)
	UpdateRemindStatus(ctx context.Context, id string, status calendar.RemindStatus) error
	GetOldEventIDs(ctx context.Context, oldTime time.Time) ([]string, error)
	DeleteEventByIDs(ctx context.Context, ids []string) error
}

type Scheduler struct {
	rep       EventRepository
	publisher queue.Producer
	cfg       *Config
	log       logger.Logger
}

func New(log logger.Logger, rep EventRepository, sender queue.Producer, cfg *Config) *Scheduler {
	return &Scheduler{
		rep:       rep,
		publisher: sender,
		cfg:       cfg,
		log:       log,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	if errP := s.publisher.Connect(ctx); errP != nil {
		return errP
	}
	ticker := time.NewTicker(s.cfg.Duration)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if errN := s.remindEvents(ctx); errN != nil {
				return errN
			}
			if errD := s.deleteOldEvents(ctx); errD != nil {
				return errD
			}
		}
	}
}

func (s *Scheduler) remindEvents(ctx context.Context) error {
	events, errR := s.rep.GetEventsWithRemindStatus(ctx, time.Now(), calendar.NotSentStatus)
	if errR != nil {
		return errR
	}
	if len(events) == 0 {
		return nil
	}
	for _, e := range events {
		data, errM := json.Marshal(calendar.ToCalendarEventDto(e))
		if errM != nil {
			return errM
		}
		if errS := s.publisher.Publish(data); errS != nil {
			return errS
		}
		if errU := s.rep.UpdateRemindStatus(ctx, e.ID, calendar.SentStatus); errU != nil {
			return errU
		}
		s.log.Debug("Notify events", zap.String("id", e.ID))
	}

	return nil
}

func (s *Scheduler) deleteOldEvents(ctx context.Context) error {
	ids, err := s.rep.GetOldEventIDs(ctx, time.Now().AddDate(-1, 0, 0))
	if err != nil {
		return err
	}
	if errD := s.rep.DeleteEventByIDs(ctx, ids); err != nil {
		return errD
	}
	if len(ids) > 0 {
		s.log.Debug("Deleted old events", zap.Strings("ids", ids))
	}
	return nil
}
