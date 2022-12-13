package app

import (
	"context"
	"errors"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	internalhttp "github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
)

var (
	ErrCreateExistEvent = errors.New("error create exist event")
	ErrEventNotExist    = errors.New("error event not exist")
	ErrDateBusy         = errors.New("error date for event busy")
)

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

func (a *App) CreateEvent(ctx context.Context, event internalhttp.EventDto) error {
	events, err := a.storage.EventsByPeriod(ctx, event.DateTimeStart, event.DateTimeEnd, 1)
	if err != nil {
		return err
	}
	if len(events) > 0 {
		return ErrCreateExistEvent
	}

	return a.storage.CreateEvent(ctx, assemble(event))
}

func (a *App) UpdateEvent(ctx context.Context, event internalhttp.EventDto) error {
	events, err := a.storage.EventsByPeriod(ctx, event.DateTimeStart, event.DateTimeEnd, 1)
	if err != nil {
		return err
	}
	if len(events) > 0 {
		return ErrDateBusy
	}

	return a.storage.UpdateEvent(ctx, assemble(event))
}

func (a *App) ReadEvent(ctx context.Context, id string) (*calendar.Event, error) {
	return a.storage.EventByID(ctx, id)
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

func assemble(eventDto internalhttp.EventDto) *calendar.Event {
	return &calendar.Event{
		Title:         eventDto.Title,
		DateTimeStart: eventDto.DateTimeStart,
		DateTimeEnd:   eventDto.DateTimeEnd,
		Description:   eventDto.Description,
		CreatedBy:     calendar.UserID(eventDto.CreatedBy),
		RemindFrom:    eventDto.RemindFrom,
	}
}
