package grpc

import (
	"context"
	"fmt"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	event "github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/server/grpc/pb/api/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	NotFountEventsError = "Not found events"
	NotFountEventError  = "Not found event"
)

type Handler struct {
	app calendar.Application
	event.UnimplementedEventServiceServer
}

func (h *Handler) Create(ctx context.Context, e *event.Event) (*event.Response, error) {
	uuid, err := h.app.CreateEvent(ctx, toEventDto(e))
	if err != nil {
		return &event.Response{
			Status: false,
			Reason: err.Error(),
		}, err
	}
	return &event.Response{
		Status: true,
		Reason: fmt.Sprintf("Event created uuid=%s", uuid),
	}, nil
}

func (h *Handler) Delete(ctx context.Context, r *event.DeleteRequest) (*event.Response, error) {
	err := h.app.DeleteEvent(ctx, r.Id)
	if err != nil {
		return &event.Response{
			Status: false,
			Reason: err.Error(),
		}, err
	}
	return &event.Response{
		Status: true,
		Reason: fmt.Sprintf("Event deleted uuid=%s", r.Id),
	}, nil
}

func (h *Handler) GetEventsOnDay(ctx context.Context, r *event.EventsRequest) (*event.FilterResponse, error) {
	events, err := h.app.GetEventsOnDay(ctx, r.GetDateFrom().AsTime())
	if err != nil {
		return &event.FilterResponse{
			Status: false,
			Reason: err.Error(),
		}, err
	}
	if events == nil {
		return &event.FilterResponse{
			Status: false,
			Reason: NotFountEventsError,
		}, nil
	}
	return &event.FilterResponse{
		Status: true,
		Events: &event.Events{Events: toEventsPb(events)},
	}, nil
}

func (h *Handler) GetEventsOnWeek(ctx context.Context, r *event.EventsRequest) (*event.FilterResponse, error) {
	events, err := h.app.GetEventsOnWeek(ctx, r.GetDateFrom().AsTime())
	if err != nil {
		return &event.FilterResponse{
			Status: false,
			Reason: err.Error(),
		}, err
	}
	if events == nil {
		return &event.FilterResponse{
			Status: false,
			Reason: NotFountEventsError,
		}, nil
	}

	return &event.FilterResponse{
		Status: true,
		Events: &event.Events{Events: toEventsPb(events)},
	}, nil
}

func (h *Handler) GetEventsOnMonth(ctx context.Context, r *event.EventsRequest) (*event.FilterResponse, error) {
	events, err := h.app.GetEventsOnMonth(ctx, r.GetDateFrom().AsTime())
	if err != nil {
		return &event.FilterResponse{
			Status: false,
			Reason: err.Error(),
		}, err
	}
	if events == nil {
		return &event.FilterResponse{
			Status: false,
			Reason: NotFountEventsError,
		}, nil
	}
	return &event.FilterResponse{
		Status: true,
		Events: &event.Events{Events: toEventsPb(events)},
	}, nil
}

func (h *Handler) Reader(ctx context.Context, eventUUID *event.EventUUID) (*event.EventResponse, error) {
	e, err := h.app.ReadEvent(ctx, eventUUID.Id)
	if err != nil {
		return &event.EventResponse{
			Status: false,
			Reason: err.Error(),
		}, err
	}
	if e == nil {
		return &event.EventResponse{
			Status: false,
			Reason: NotFountEventError,
		}, nil
	}
	return &event.EventResponse{
		Status: true,
		Event:  toEventPb(e),
	}, nil
}

func (h *Handler) Update(ctx context.Context, r *event.UpdateRequest) (*event.Response, error) {
	err := h.app.UpdateEvent(ctx, r.Id, toEventDto(r.Event))
	if err != nil {
		return &event.Response{
			Status: false,
			Reason: err.Error(),
		}, err
	}
	return &event.Response{
		Status: true,
		Reason: "Event updated",
	}, nil
}

func toEventsPb(appEvents []*calendar.EventDto) []*event.Event {
	eventsPb := make([]*event.Event, 0, len(appEvents))
	for _, e := range appEvents {
		eventsPb = append(eventsPb, toEventPb(e))
	}

	return eventsPb
}

func toEventPb(e *calendar.EventDto) *event.Event {
	return &event.Event{
		Title:         e.Title,
		DateTimeStart: &timestamppb.Timestamp{Seconds: e.DateTimeStart.Unix()},
		DateTimeEnd:   &timestamppb.Timestamp{Seconds: e.DateTimeEnd.Unix()},
		Description:   e.Description,
		UserId:        e.CreatedBy,
		RemindFrom:    &timestamppb.Timestamp{Seconds: e.RemindFrom.Unix()},
	}
}

func toEventDto(e *event.Event) *calendar.EventDto {
	return &calendar.EventDto{
		Title:         e.Title,
		DateTimeStart: e.DateTimeStart.AsTime(),
		DateTimeEnd:   e.DateTimeEnd.AsTime(),
		Description:   e.Description,
		CreatedBy:     e.UserId,
		RemindFrom:    e.RemindFrom.AsTime(),
	}
}
