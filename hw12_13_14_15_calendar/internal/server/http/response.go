package internalhttp

import "net/http"

const (
	validateErrorCode = "Validate error"
	errorReadRequest  = "Error read request"
	errorCreateEvent  = "Error create event"
	errorGetEvent     = "Error get event"

	eventCreateSuccessCode  = "Event created success"
	eventReadSuccessCode    = "Event read success"
	eventUpdateSuccessCode  = "Event updated success"
	eventNotFoundCode       = "Event not found"
	eventDeletedSuccessCode = "Event deleted success"
	eventsReadSuccessCode   = "Events read success"
	eventsNotFoundCode      = "Events not found"

	eventDeletedErrorCode = "Error delete event"
	eventUpdateErrorCode  = "Error update event"
	eventAlreadyDeleted   = "Event already deleted"
	eventAlreadyExistCode = "Event already exist"
	errorGetEventsCode    = "Error get events"

	unsupportedActionCode = "Unsupported action"
	handlerNotFound       = "Handler not fount"
	successExecute        = "Success execute"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse(code, message string, data interface{}) *Response {
	if message == "" {
		message = code
	}
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (r *Response) getStatus() int {
	switch r.Code {
	case validateErrorCode:
		return http.StatusUnprocessableEntity
	case errorCreateEvent, errorGetEvent, eventDeletedErrorCode, eventUpdateErrorCode, errorGetEventsCode:
		return http.StatusInternalServerError
	case eventCreateSuccessCode:
		return http.StatusCreated
	case eventNotFoundCode, handlerNotFound, eventsNotFoundCode:
		return http.StatusNotFound
	case eventAlreadyExistCode:
		return http.StatusConflict
	case eventDeletedSuccessCode, eventReadSuccessCode, eventUpdateSuccessCode, eventsReadSuccessCode:
		return http.StatusOK
	case unsupportedActionCode:
		return http.StatusMethodNotAllowed
	case eventAlreadyDeleted:
		return http.StatusGone
	case errorReadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusOK
	}
}
