package internalhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"go.uber.org/zap"
)

type middlewareFunc func(http.HandlerFunc) http.HandlerFunc

type Handler struct {
	app              calendar.Application
	logger           logger.Logger
	urlEvent         *regexp.Regexp
	urlEventsOnDay   *regexp.Regexp
	urlEventsOnWeek  *regexp.Regexp
	urlEventsOnMonth *regexp.Regexp
	middlewares      []middlewareFunc
}

func (h *Handler) AddMiddleware(m middlewareFunc) {
	h.middlewares = append(h.middlewares, m)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := h.getHandler(r)
	for _, m := range h.middlewares {
		handler = m(handler)
	}
	handler.ServeHTTP(w, r)
}

func NewHandler(app calendar.Application, logger logger.Logger) *Handler {
	return &Handler{
		app:              app,
		logger:           logger,
		urlEvent:         regexp.MustCompile(fmt.Sprintf(`/events/(?m)(?P<uuid>%s)`, UUIDRegular)),
		urlEventsOnDay:   regexp.MustCompile(fmt.Sprintf(`/events/day/(?m)(?P<date>%s)`, DateRegular)),
		urlEventsOnWeek:  regexp.MustCompile(fmt.Sprintf(`/events/week/(?m)(?P<date>%s)`, DateRegular)),
		urlEventsOnMonth: regexp.MustCompile(fmt.Sprintf(`/events/month/(?m)(?P<date>%s)`, DateRegular)),
		middlewares:      make([]middlewareFunc, 0),
	}
}

func (h *Handler) getHandler(r *http.Request) http.HandlerFunc {
	var c http.HandlerFunc
	switch r.URL.String() {
	case "/":
		c = h.mainPage
	default:
		if c = h.getEventHandler(r); c != nil {
			break
		}
		c = h.notFound
	}

	return c
}

func (h *Handler) getEventHandler(r *http.Request) http.HandlerFunc {
	if !strings.Contains(r.URL.String(), "/events") {
		return nil
	}

	if r.URL.String() == "/events" && http.MethodPost == r.Method {
		return h.createEvent
	}
	if h.urlEvent.Match([]byte(r.URL.String())) {
		switch r.Method {
		case http.MethodGet:
			return h.getEvent
		case http.MethodDelete:
			return h.deleteEvent
		case http.MethodPut:
			return h.updateEvent
		default:
			return h.actionNotSupported
		}
	}
	if r.Method != http.MethodGet {
		return nil
	}
	if h.urlEventsOnDay.Match([]byte(r.URL.String())) {
		return h.getEventsOnDay
	}
	if h.urlEventsOnWeek.Match([]byte(r.URL.String())) {
		return h.getEventsOnWeek
	}
	if h.urlEventsOnMonth.Match([]byte(r.URL.String())) {
		return h.getEventsOnMonth
	}
	return nil
}

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(NewResponse(successExecute, "hello-world", nil), w)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(NewResponse(handlerNotFound, "Not found", nil), w)
}

func (h *Handler) actionNotSupported(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(NewResponse(unsupportedActionCode, "", nil), w)
}

func (h *Handler) createEvent(w http.ResponseWriter, r *http.Request) {
	req := h.parserEventRequest(w, r)
	if req == nil {
		return
	}
	eventDto, errE := toEventDto(req)
	if errE != nil {
		h.sendResponse(NewResponse(errorCreateEvent, errE.Error(), nil), w)
		return
	}
	eventUUID, err := h.app.CreateEvent(r.Context(), eventDto)
	if err != nil {
		h.sendResponse(NewResponse(errorCreateEvent, err.Error(), nil), w)
		return
	}
	h.sendResponse(NewResponse(eventCreateSuccessCode, fmt.Sprintf("Event created with uuid=%s", eventUUID), nil), w)
}

func (h *Handler) getEvent(w http.ResponseWriter, r *http.Request) {
	e, err := h.app.ReadEvent(r.Context(), h.getEventUUID(r))
	if err != nil {
		h.sendResponse(NewResponse(errorGetEvent, err.Error(), nil), w)
		return
	}
	if e == nil {
		h.sendResponse(NewResponse(eventNotFoundCode, "", nil), w)
		return
	}

	h.sendResponse(NewResponse(eventReadSuccessCode, "", e), w)
}

func (h *Handler) deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventUUID := h.getEventUUID(r)
	if err := h.app.DeleteEvent(r.Context(), eventUUID); err != nil {
		h.sendResponse(NewResponse(eventDeletedErrorCode, err.Error(), nil), w)
		return
	}

	h.sendResponse(NewResponse(
		eventDeletedSuccessCode,
		fmt.Sprintf("Event deleted with uuid=%s", eventUUID),
		nil,
	), w)
}

func (h *Handler) updateEvent(w http.ResponseWriter, r *http.Request) {
	eventUUID := h.getEventUUID(r)
	req := h.parserEventRequest(w, r)
	if req == nil {
		return
	}
	eventDto, errE := toEventDto(req)
	if errE != nil {
		h.sendResponse(NewResponse(eventUpdateErrorCode, errE.Error(), nil), w)
		return
	}
	if err := h.app.UpdateEvent(r.Context(), eventUUID, eventDto); err != nil {
		h.sendResponse(NewResponse(eventUpdateErrorCode, err.Error(), nil), w)
		return
	}

	h.sendResponse(NewResponse(
		eventUpdateSuccessCode,
		fmt.Sprintf("Event updated with uuid=%s", eventUUID),
		nil,
	), w)
}

func (h *Handler) getEventsOnDay(w http.ResponseWriter, r *http.Request) {
	date := parseDate(h.urlEventsOnDay, r.URL.String())
	day, errD := time.Parse(DateFormat, string(date))
	if errD != nil {
		h.sendResponse(NewResponse(errorGetEventsCode, errD.Error(), nil), w)
		return
	}
	events, errA := h.app.GetEventsOnDay(r.Context(), day)
	if errA != nil {
		h.sendResponse(NewResponse(errorGetEventsCode, errA.Error(), nil), w)
		return
	}
	if events == nil {
		h.sendResponse(NewResponse(eventsNotFoundCode, "", events), w)
		return
	}
	h.sendResponse(NewResponse(eventsReadSuccessCode, "", events), w)
}

func (h *Handler) getEventsOnWeek(w http.ResponseWriter, r *http.Request) {
	date := parseDate(h.urlEventsOnWeek, r.URL.String())
	day, errD := time.Parse(DateFormat, string(date))
	if errD != nil {
		h.sendResponse(NewResponse(errorGetEventsCode, errD.Error(), nil), w)
		return
	}
	events, errA := h.app.GetEventsOnWeek(r.Context(), day)
	if errA != nil {
		h.sendResponse(NewResponse(errorGetEventsCode, errA.Error(), nil), w)
		return
	}
	if events == nil {
		h.sendResponse(NewResponse(eventsNotFoundCode, "", events), w)
		return
	}
	h.sendResponse(NewResponse(eventsReadSuccessCode, "", events), w)
}

func (h *Handler) getEventsOnMonth(w http.ResponseWriter, r *http.Request) {
	date := parseDate(h.urlEventsOnMonth, r.URL.String())
	day, errD := time.Parse(DateFormat, string(date))
	if errD != nil {
		h.sendResponse(NewResponse(errorGetEventsCode, errD.Error(), nil), w)
		return
	}
	events, errA := h.app.GetEventsOnMonth(r.Context(), day)
	if errA != nil {
		h.sendResponse(NewResponse(errorGetEventsCode, errA.Error(), nil), w)
		return
	}
	if events == nil {
		h.sendResponse(NewResponse(eventsNotFoundCode, "", events), w)
		return
	}
	h.sendResponse(NewResponse(eventsReadSuccessCode, "", events), w)
}

func (h *Handler) parserEventRequest(w http.ResponseWriter, r *http.Request) *EventRequest {
	req := &EventRequest{}
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		h.sendResponse(NewResponse(errorReadRequest, err.Error(), nil), w)
		return nil
	}
	if errU := json.Unmarshal(data, req); errU != nil {
		h.sendResponse(NewResponse(errorReadRequest, errU.Error(), nil), w)
		return nil
	}
	if errV := req.Validate(); errV != nil {
		h.sendResponse(NewResponse(validateErrorCode, errV.Error(), nil), w)
		return nil
	}
	return req
}

func (h *Handler) getEventUUID(r *http.Request) string {
	var uuidEvent []byte
	uuidEvent = h.urlEvent.Expand(
		uuidEvent,
		[]byte("$uuid"),
		[]byte(r.URL.String()), h.urlEvent.FindSubmatchIndex([]byte(r.URL.String())),
	)
	return string(uuidEvent)
}

func (h *Handler) sendResponse(response *Response, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("error marshal data", zap.Error(err))
		return
	}
	w.WriteHeader(response.getStatus())
	if _, err := w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Error("error send data to client", zap.Error(err))
		return
	}
}

func parseDate(reg *regexp.Regexp, url string) []byte {
	var date []byte
	date = reg.Expand(date, []byte("$date"), []byte(url), reg.FindSubmatchIndex([]byte(url)))

	return date
}
