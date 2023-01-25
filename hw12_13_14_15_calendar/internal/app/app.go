package app

import (
	"context"
	"errors"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrCreateExistEvent = errors.New("error create exist event")
	ErrEventNotExist    = errors.New("error event not exist")
	ErrDateBusy         = errors.New("error date for event busy")
)

const LimitEventsOnQuery = 1000

type App struct {
	logger  logger.Logger
	storage storage.Storage
}

func New(logger logger.Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *calendar.EventDto) (string, error) {
	events, err := a.storage.EventsByPeriod(ctx, event.DateTimeStart, event.DateTimeEnd, 1)
	if err != nil {
		a.logger.Debug("Found events when create event", zap.Int("count_events", len(events)))
		return "", err
	}
	if len(events) > 0 {
		return "", ErrCreateExistEvent
	}
	e := toCalendarEvent(event)
	e.ID = uuid.New().String()

	return e.ID, a.storage.CreateEvent(ctx, e)
}

func (a *App) UpdateEvent(ctx context.Context, uuid string, event *calendar.EventDto) error {
	eventsInPeriod, err := a.storage.EventsByPeriod(ctx, event.DateTimeStart, event.DateTimeEnd, 1)
	if err != nil {
		a.logger.Error("Error found event", zap.Error(err))
		return err
	}
	var events []*calendar.Event
	for _, e := range eventsInPeriod {
		if e.ID != uuid {
			events = append(events, e)
		}
	}
	if len(events) > 0 {
		return ErrDateBusy
	}
	eventDto := toCalendarEvent(event)
	eventDto.ID = uuid
	return a.storage.UpdateEvent(ctx, eventDto)
}

func (a *App) ReadEvent(ctx context.Context, id string) (*calendar.EventDto, error) {
	event, err := a.storage.EventByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, nil
	}
	return toCalendarEventDto(event), nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	event, err := a.storage.EventByID(ctx, id)
	if err != nil {
		return err
	}
	if event == nil {
		return ErrEventNotExist
	}

	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEventsOnDay(ctx context.Context, day time.Time) ([]*calendar.EventDto, error) {
	events := make([]*calendar.EventDto, 0)
	eventsInPeriod, err := a.storage.EventsByPeriod(ctx, day, day.AddDate(0, 0, 1), LimitEventsOnQuery)
	if err != nil {
		return events, err
	}
	for _, e := range eventsInPeriod {
		events = append(events, toCalendarEventDto(e))
	}
	return events, nil
}

func (a *App) GetEventsOnWeek(ctx context.Context, day time.Time) ([]*calendar.EventDto, error) {
	events := make([]*calendar.EventDto, 0)
	eventsInPeriod, err := a.storage.EventsByPeriod(ctx, day, day.AddDate(0, 0, 7), LimitEventsOnQuery)
	if err != nil {
		return events, err
	}
	for _, e := range eventsInPeriod {
		events = append(events, toCalendarEventDto(e))
	}
	return events, nil
}

func (a *App) GetEventsOnMonth(ctx context.Context, day time.Time) ([]*calendar.EventDto, error) {
	events := make([]*calendar.EventDto, 0)
	eventsInPeriod, err := a.storage.EventsByPeriod(ctx, day, day.AddDate(0, 1, 0), LimitEventsOnQuery)
	if err != nil {
		return events, err
	}
	for _, e := range eventsInPeriod {
		events = append(events, toCalendarEventDto(e))
	}
	return events, nil
}

func toCalendarEvent(eventDto *calendar.EventDto) *calendar.Event {
	return &calendar.Event{
		Title:         eventDto.Title,
		DateTimeStart: eventDto.DateTimeStart,
		DateTimeEnd:   eventDto.DateTimeEnd,
		Description:   eventDto.Description,
		CreatedBy:     calendar.UserID(eventDto.CreatedBy),
		RemindFrom:    eventDto.RemindFrom,
	}
}

func toCalendarEventDto(e *calendar.Event) *calendar.EventDto {
	return &calendar.EventDto{
		ID:            e.ID,
		Title:         e.Title,
		DateTimeStart: e.DateTimeStart,
		DateTimeEnd:   e.DateTimeEnd,
		Description:   e.Description,
		CreatedBy:     int32(e.CreatedBy),
		RemindFrom:    e.RemindFrom,
	}
}
