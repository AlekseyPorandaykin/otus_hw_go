package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
)

var ErrNotExist = errors.New("event not exist")

type Storage struct {
	mu     sync.RWMutex
	events map[string]*calendar.Event
}

func (m *Storage) CreateEvent(ctx context.Context, e *calendar.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events[e.ID] = e
	return nil
}

func (m *Storage) UpdateEvent(ctx context.Context, e *calendar.Event) error {
	e, err := m.EventByID(ctx, e.ID)
	if err != nil {
		return err
	}
	if e == nil {
		return ErrNotExist
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events[e.ID] = e
	return nil
}

func (m *Storage) DeleteEvent(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.events, id)
	return nil
}

func (m *Storage) EventByID(ctx context.Context, id string) (*calendar.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.events[id], nil
}

func (m *Storage) Connect(ctx context.Context) error {
	return nil
}

func (m *Storage) Close(ctx context.Context) error {
	m.events = map[string]*calendar.Event{}
	return nil
}

func (m *Storage) EventsByPeriod(ctx context.Context, start, end time.Time, limit int) ([]*calendar.Event, error) {
	var res []*calendar.Event
	for _, e := range m.events {
		if len(res) >= limit {
			break
		}
		if e.DateTimeStart.After(start) && e.DateTimeEnd.Before(end) {
			res = append(res, e)
		}
	}
	return res, nil
}

func New() *Storage {
	return &Storage{
		events: map[string]*calendar.Event{},
	}
}
