package internalhttp

import (
	"errors"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
)

const (
	UUIDRegular = `[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}`
	DateRegular = `[0-9]{4}-[0-9]{2}-[0-9]{2}`
)

type RequestValidator interface {
	Validate() error
}

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
)

type EventRequest struct {
	Title         string `json:"title"`
	DateTimeStart string `json:"dateTimeStart"`
	DateTimeEnd   string `json:"dateTimeEnd"`
	Description   string `json:"description"`
	CreatedBy     int32  `json:"createdBy"`
	RemindFrom    string `json:"remindFrom"`
}

func (e *EventRequest) Validate() error {
	if e.Title == "" {
		return errors.New("empty title")
	}
	if e.Description == "" {
		return errors.New("empty description")
	}
	if e.CreatedBy == 0 {
		return errors.New("empty createdBy")
	}
	if _, err := e.GetDateTimeStart(); err != nil {
		return errors.New("incorrect dateTimeStart")
	}
	if _, err := e.GetDateTimeEnd(); err != nil {
		return errors.New("incorrect dateTimeEnd")
	}
	if _, err := e.GetRemindFrom(); err != nil {
		return errors.New("incorrect remindFrom")
	}
	return nil
}

func (e *EventRequest) GetDateTimeStart() (time.Time, error) {
	return time.Parse(DateTimeFormat, e.DateTimeStart)
}

func (e *EventRequest) GetDateTimeEnd() (time.Time, error) {
	return time.Parse(DateTimeFormat, e.DateTimeEnd)
}

func (e *EventRequest) GetRemindFrom() (time.Time, error) {
	return time.Parse(DateTimeFormat, e.RemindFrom)
}

func toEventDto(req *EventRequest) (*calendar.EventDto, error) {
	dateTimeStart, errS := req.GetDateTimeStart()
	if errS != nil {
		return nil, errS
	}
	dateTimeEnd, errE := req.GetDateTimeEnd()
	if errE != nil {
		return nil, errE
	}
	remindFrom, errR := req.GetRemindFrom()
	if errR != nil {
		return nil, errR
	}
	return &calendar.EventDto{
		Title:         req.Title,
		Description:   req.Description,
		CreatedBy:     req.CreatedBy,
		DateTimeStart: dateTimeStart,
		DateTimeEnd:   dateTimeEnd,
		RemindFrom:    remindFrom,
	}, nil
}
