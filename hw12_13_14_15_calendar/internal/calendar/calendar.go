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

type RemindStatus int

const (
	NotSentStatus RemindStatus = iota
	SentStatus
	CancelSendStatus
)

type Event struct {
	ID            string       `db:"id"`
	Title         string       `db:"title"`
	DateTimeStart time.Time    `db:"datetime_from"`
	DateTimeEnd   time.Time    `db:"datetime_to"`
	Description   string       `db:"description"`
	CreatedBy     UserID       `db:"created_by"`
	RemindFrom    time.Time    `db:"start_notify"`
	RemindStatus  RemindStatus `db:"notify_status"`
	CreatedAt     time.Time    `db:"created_at"`
}

type UserID int

type Notification struct {
	ID       string
	Title    string
	Datetime time.Time
	UserTo   UserID
}

func ToCalendarEvent(eventDto *EventDto) *Event {
	statusRemind := NotSentStatus
	if eventDto.RemindFrom == eventDto.DateTimeStart {
		statusRemind = CancelSendStatus
	}
	return &Event{
		Title:         eventDto.Title,
		DateTimeStart: eventDto.DateTimeStart,
		DateTimeEnd:   eventDto.DateTimeEnd,
		Description:   eventDto.Description,
		CreatedBy:     UserID(eventDto.CreatedBy),
		RemindFrom:    eventDto.RemindFrom,
		RemindStatus:  statusRemind,
		CreatedAt:     time.Now(),
	}
}

func ToCalendarEventDto(e *Event) *EventDto {
	return &EventDto{
		ID:            e.ID,
		Title:         e.Title,
		DateTimeStart: e.DateTimeStart,
		DateTimeEnd:   e.DateTimeEnd,
		Description:   e.Description,
		CreatedBy:     int32(e.CreatedBy),
		RemindFrom:    e.RemindFrom,
	}
}
