package calendar

import (
	"context"
	"time"
)

type EventDto struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	DateTimeStart time.Time `json:"dateTimeStart"`
	DateTimeEnd   time.Time `json:"dateTimeEnd"`
	Description   string    `json:"description"`
	CreatedBy     int32     `json:"createdBy"`
	RemindFrom    time.Time `json:"remindFrom"`
}

type Application interface {
	CreateEvent(ctx context.Context, event *EventDto) (string, error)
	ReadEvent(ctx context.Context, id string) (*EventDto, error)
	UpdateEvent(ctx context.Context, uuid string, event *EventDto) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventsOnDay(ctx context.Context, day time.Time) ([]*EventDto, error)
	GetEventsOnWeek(ctx context.Context, fromDay time.Time) ([]*EventDto, error)
	GetEventsOnMonth(ctx context.Context, fromDay time.Time) ([]*EventDto, error)
}

type Event struct {
	ID            string    `db:"id"`
	Title         string    `db:"title"`
	DateTimeStart time.Time `db:"datetime_from"`
	DateTimeEnd   time.Time `db:"datetime_to"`
	Description   string    `db:"description"`
	CreatedBy     UserID    `db:"created_by"`
	RemindFrom    time.Time `db:"start_notify"`
}

type UserID int

type Notification struct {
	ID       string
	Title    string
	Datetime time.Time
	UserTo   UserID
}
