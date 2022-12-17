package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	s := New()
	datetime := time.Now()
	eventOne := &calendar.Event{
		ID:            "00-01",
		Title:         "Event title 1",
		DateTimeStart: time.Now(),
		DateTimeEnd:   time.Now(),
		Description:   "Event description",
		CreatedBy:     calendar.UserID(1),
		RemindFrom:    time.Now(),
	}

	eventTwo := &calendar.Event{
		ID:            "00-02",
		Title:         "Event title 2",
		DateTimeStart: datetime.Add(-time.Minute * 15),
		DateTimeEnd:   datetime,
		Description:   "Event description",
		CreatedBy:     calendar.UserID(2),
		RemindFrom:    datetime,
	}
	eventThree := &calendar.Event{
		ID:            "00-03",
		Title:         "Event title 3",
		DateTimeStart: datetime,
		DateTimeEnd:   datetime.Add(time.Minute * 15),
		Description:   "Event description",
		CreatedBy:     calendar.UserID(3),
		RemindFrom:    datetime,
	}
	eventFour := &calendar.Event{
		ID:            "00-04",
		Title:         "Event title 4",
		DateTimeStart: datetime,
		DateTimeEnd:   datetime,
		Description:   "Event description",
		CreatedBy:     calendar.UserID(4),
		RemindFrom:    datetime,
	}

	require.Nil(t, s.Connect(ctx))
	require.Nil(t, s.CreateEvent(ctx, eventOne))
	require.Nil(t, s.CreateEvent(ctx, eventTwo))
	require.Nil(t, s.CreateEvent(ctx, eventThree))
	eventActual, err := s.EventByID(ctx, "00-01")
	require.Nil(t, err)
	require.Equal(t, eventOne, eventActual)

	eventsOne, err := s.EventsByPeriod(ctx, time.Now().Add(-time.Minute*5), time.Now().Add(time.Minute*5), 10)
	require.Nil(t, err)
	require.Equal(t, []*calendar.Event{eventOne}, eventsOne)
	eventsTwo, err := s.EventsByPeriod(ctx, time.Now().Add(-time.Hour*5), time.Now().Add(time.Hour*5), 2)
	require.Nil(t, err)
	require.Equal(t, 2, len(eventsTwo))

	require.Nil(t, s.DeleteEvent(ctx, "00-01"))
	eventOneActual, err := s.EventByID(ctx, "00-01")
	require.Nil(t, eventOneActual)
	require.Nil(t, err)

	eventTwo.Title = "Other title"
	require.Nil(t, s.UpdateEvent(ctx, eventTwo))

	eventTwoActual, err := s.EventByID(ctx, "00-02")
	require.Nil(t, err)
	require.Equal(t, "Other title", eventTwoActual.Title)

	require.ErrorIs(t, s.UpdateEvent(ctx, eventFour), ErrNotExist)
	require.Nil(t, s.Close(ctx))
}
