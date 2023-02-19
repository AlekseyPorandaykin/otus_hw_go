package storage

import (
	"context"
	"errors"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/sender"
	memorystorage "github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/storage/sql"
)

var ErrNotSupportedDriver = errors.New("not supported driver storage")

type Storage interface {
	EventStorage
	scheduler.EventRepository
	sender.Saver
}

type EventStorage interface {
	CreateEvent(ctx context.Context, e *calendar.Event) error
	UpdateEvent(ctx context.Context, e *calendar.Event) error
	DeleteEvent(ctx context.Context, id string) error
	EventByID(ctx context.Context, id string) (*calendar.Event, error)
	EventsByPeriod(ctx context.Context, start, end time.Time, limit int) ([]*calendar.Event, error)
	Close(ctx context.Context) error
}

func CreateStorage(conf *config.StorageConfig) (Storage, error) {
	switch conf.Storage {
	case "sql":
		s := sqlstorage.New(conf)
		if err := s.Connect(context.Background()); err != nil {
			return nil, err
		}

		return s, nil
	case "memory":
		return memorystorage.New(), nil
	}

	return nil, ErrNotSupportedDriver
}
