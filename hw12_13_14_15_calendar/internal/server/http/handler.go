package internalhttp

import (
	"net/http"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	app    Application
	logger logger.Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var c func(http.ResponseWriter, *http.Request)

	switch r.URL.String() {
	case "/":
		c = h.mainPage
	case "/ping":
		c = h.ping
	default:
		c = h.notFound
	}
	h.loggingMiddleware(c).ServeHTTP(w, r)
}

func NewHandler(app Application, logger logger.Logger) *Handler {
	return &Handler{app: app, logger: logger}
}

func (h *Handler) mainPage(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"message":"hello-world"}"`, http.StatusOK, w)
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"message":"pong"}"`, http.StatusOK, w)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) {
	h.sendResponse(`{"message":"Not found"}"`, http.StatusNotFound, w)
}

func (h *Handler) sendResponse(data string, status int, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write([]byte(data)); err != nil {
		status = http.StatusInternalServerError
		h.logger.Error("error send data to client", zap.Error(err))
	}
	w.WriteHeader(status)
}
