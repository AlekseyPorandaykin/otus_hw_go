package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
)

var ErrNotExist = errors.New("event not exist")

type log struct {
	body []byte
	date time.Time
}

type Storage struct {
	muEvents sync.RWMutex
	events   map[string]*calendar.Event

	muLogs sync.RWMutex
	logs   map[string]log
}

func (m *Storage) Save(ctx context.Context, eventID string, body []byte, date time.Time) error {
	m.muLogs.Lock()
	defer m.muLogs.Unlock()
	m.logs[eventID] = log{body: body, date: date}
	return nil
}

func (m *Storage) CreateEvent(ctx context.Context, e *calendar.Event) error {
	m.muEvents.Lock()
	defer m.muEvents.Unlock()
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
	m.muEvents.Lock()
	defer m.muEvents.Unlock()
	m.events[e.ID] = e
	return nil
}

func (m *Storage) DeleteEvent(ctx context.Context, id string) error {
	m.muEvents.Lock()
	defer m.muEvents.Unlock()
	delete(m.events, id)
	return nil
}

func (m *Storage) EventByID(ctx context.Context, id string) (*calendar.Event, error) {
	m.muEvents.Lock()
	defer m.muEvents.Unlock()
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

func (m *Storage) GetEventsWithRemindStatus(
	ctx context.Context, time time.Time, status calendar.RemindStatus,
) ([]*calendar.Event, error) {
	events := make([]*calendar.Event, 0)
	for _, e := range m.events {
		if e.RemindStatus == calendar.NotSentStatus && (e.RemindFrom.After(time) || e.RemindFrom.Equal(time)) {
			events = append(events, e)
		}
	}
	return events, nil
}

func (m *Storage) UpdateRemindStatus(ctx context.Context, id string, status calendar.RemindStatus) error {
	m.muEvents.Lock()
	defer m.muEvents.Unlock()
	if e, ok := m.events[id]; ok && e != nil {
		e.RemindStatus = status
	}
	return nil
}

func (m *Storage) GetOldEventIDs(ctx context.Context, oldTime time.Time) ([]string, error) {
	events := make([]string, 0)
	for _, e := range m.events {
		if e.DateTimeEnd.Before(oldTime) {
			events = append(events, e.ID)
		}
	}
	return events, nil
}

func (m *Storage) DeleteEventByIDs(ctx context.Context, ids []string) error {
	m.muEvents.Lock()
	defer m.muEvents.Unlock()
	for _, id := range ids {
		delete(m.events, id)
	}
	return nil
}

func New() *Storage {
	return &Storage{
		events: map[string]*calendar.Event{},
		logs:   map[string]log{},
	}
}
