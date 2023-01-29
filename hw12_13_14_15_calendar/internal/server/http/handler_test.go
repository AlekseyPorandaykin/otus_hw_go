package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/stretchr/testify/require"
)

type mockApp struct {
	createEventRes struct {
		s string
		e error
	}
	readEventRes struct {
		c *calendar.EventDto
		e error
	}
	deleteEventRes error
	updateEventRes error
	eventsOnDay    struct {
		c []*calendar.EventDto
		e error
	}
	eventsOnWeek struct {
		c []*calendar.EventDto
		e error
	}
	eventsOnMonth struct {
		c []*calendar.EventDto
		e error
	}
}

func (m *mockApp) CreateEvent(ctx context.Context, event *calendar.EventDto) (string, error) {
	return m.createEventRes.s, m.createEventRes.e
}

func (m *mockApp) ReadEvent(ctx context.Context, id string) (*calendar.EventDto, error) {
	return m.readEventRes.c, m.readEventRes.e
}

func (m *mockApp) UpdateEvent(ctx context.Context, uuid string, event *calendar.EventDto) error {
	return m.updateEventRes
}

func (m *mockApp) DeleteEvent(ctx context.Context, id string) error {
	return m.deleteEventRes
}

func (m *mockApp) GetEventsOnDay(ctx context.Context, day time.Time) ([]*calendar.EventDto, error) {
	return m.eventsOnDay.c, m.eventsOnDay.e
}

func (m *mockApp) GetEventsOnWeek(ctx context.Context, fromDay time.Time) ([]*calendar.EventDto, error) {
	return m.eventsOnWeek.c, m.eventsOnWeek.e
}

func (m *mockApp) GetEventsOnMonth(ctx context.Context, fromDay time.Time) ([]*calendar.EventDto, error) {
	return m.eventsOnMonth.c, m.eventsOnMonth.e
}

func TestHandler_AddMiddleware(t *testing.T) {
	type args struct {
		m middlewareFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success test",
			args: args{
				func(handlerFunc http.HandlerFunc) http.HandlerFunc {
					return handlerFunc
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{middlewares: make([]middlewareFunc, 0)}
			h.AddMiddleware(tt.args.m)
			require.Equal(t, 1, len(h.middlewares))
		})
	}
}

var (
	testUUID = "ac1b717e-9b4d-4048-b869-ea3be369a5e2"
	testTime = time.Date(2022, 2, 19, 11, 15, 0, 0, time.UTC)
	testDate = "2022-02-19"
)

func TestHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		app calendar.Application
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedBody string
		expectedCode int
	}{
		{
			name: "Main page success",
			fields: fields{
				app: &mockApp{},
			},
			args: args{
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path: "/",
					},
				},
			},
			expectedBody: `{"code":"Success execute","message":"hello-world","data":null}`,
			expectedCode: 200,
		},
		{
			name: "Not found page",
			fields: fields{
				app: &mockApp{},
			},
			args: args{
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Path: "/other",
					},
				},
			},
			expectedBody: `{"code":"Handler not fount","message":"Not found","data":null}`,
			expectedCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewMockResponse()
			h := NewHandler(tt.fields.app, logger.NewMockLogger())
			h.ServeHTTP(resp, tt.args.r)
			require.Equal(t, tt.expectedBody, resp.getBody())
			require.Contains(t, resp.headers["Content-Type"], "application/json")
			require.Equal(t, tt.expectedCode, resp.status)
		})
	}
}

func TestHandler_EventActionsOnDay(t *testing.T) {
	tests := []serveHTTP{
		{
			name: "EventsOnDay success",
			app: &mockApp{eventsOnDay: struct {
				c []*calendar.EventDto
				e error
			}{c: []*calendar.EventDto{{
				ID:            testUUID,
				Title:         "Test title",
				Description:   "Test description",
				DateTimeStart: testTime,
				DateTimeEnd:   testTime,
				CreatedBy:     1,
				RemindFrom:    testTime,
			}}, e: nil}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/day/%s", testDate),
				},
			},
			expectedBody: `{"code":"Events read success","message":"Events read success",` +
				`"data":[{"id":"ac1b717e-9b4d-4048-b869-ea3be369a5e2","title":"Test title",` +
				`"dateTimeStart":"2022-02-19T11:15:00Z","dateTimeEnd":"2022-02-19T11:15:00Z",` +
				`"description":"Test description","createdBy":1,"remindFrom":"2022-02-19T11:15:00Z"}]}`,
			expectedCode: 200,
		},
		{
			name: "EventsOnDay error not found",
			app: &mockApp{eventsOnDay: struct {
				c []*calendar.EventDto
				e error
			}{c: nil, e: nil}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/day/%s", testDate),
				},
			},
			expectedBody: `{"code":"Events not found","message":"Events not found","data":null}`,
			expectedCode: 404,
		},
		{
			name: "EventsOnDay error app",
			app: &mockApp{eventsOnDay: struct {
				c []*calendar.EventDto
				e error
			}{c: nil, e: errors.New("test error")}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/day/%s", testDate),
				},
			},
			expectedBody: `{"code":"Error get events","message":"test error","data":null}`,
			expectedCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewMockResponse()
			h := NewHandler(tt.app, logger.NewMockLogger())
			h.ServeHTTP(resp, tt.req)
			require.Equal(t, tt.expectedBody, resp.getBody())
			require.Contains(t, resp.headers["Content-Type"], "application/json")
			require.Equal(t, tt.expectedCode, resp.status)
		})
	}
}

func TestHandler_EventActionsOnWeek(t *testing.T) {
	tests := []serveHTTP{
		{
			name: "EventsOnWeek success",
			app: &mockApp{eventsOnWeek: struct {
				c []*calendar.EventDto
				e error
			}{c: []*calendar.EventDto{{
				ID:            testUUID,
				Title:         "Test title",
				Description:   "Test description",
				DateTimeStart: testTime,
				DateTimeEnd:   testTime,
				CreatedBy:     1,
				RemindFrom:    testTime,
			}}, e: nil}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/week/%s", testDate),
				},
			},
			expectedBody: `{"code":"Events read success","message":"Events read success",` +
				`"data":[{"id":"ac1b717e-9b4d-4048-b869-ea3be369a5e2","title":"Test title",` +
				`"dateTimeStart":"2022-02-19T11:15:00Z","dateTimeEnd":"2022-02-19T11:15:00Z",` +
				`"description":"Test description","createdBy":1,"remindFrom":"2022-02-19T11:15:00Z"}]}`,
			expectedCode: 200,
		},
		{
			name: "EventsOnWeek error not found",
			app: &mockApp{eventsOnWeek: struct {
				c []*calendar.EventDto
				e error
			}{c: nil, e: nil}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/week/%s", testDate),
				},
			},
			expectedBody: `{"code":"Events not found","message":"Events not found","data":null}`,
			expectedCode: 404,
		},
		{
			name: "EventsOnWeek error app",
			app: &mockApp{eventsOnWeek: struct {
				c []*calendar.EventDto
				e error
			}{c: nil, e: errors.New("test error")}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/week/%s", testDate),
				},
			},
			expectedBody: `{"code":"Error get events","message":"test error","data":null}`,
			expectedCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewMockResponse()
			h := NewHandler(tt.app, logger.NewMockLogger())
			h.ServeHTTP(resp, tt.req)
			require.Equal(t, tt.expectedBody, resp.getBody())
			require.Contains(t, resp.headers["Content-Type"], "application/json")
			require.Equal(t, tt.expectedCode, resp.status)
		})
	}
}

func TestHandler_EventActionsOnMonth(t *testing.T) {
	tests := []serveHTTP{
		{
			name: "EventsOnMonth success",
			app: &mockApp{eventsOnMonth: struct {
				c []*calendar.EventDto
				e error
			}{c: []*calendar.EventDto{{
				ID:            testUUID,
				Title:         "Test title",
				Description:   "Test description",
				DateTimeStart: testTime,
				DateTimeEnd:   testTime,
				CreatedBy:     1,
				RemindFrom:    testTime,
			}}, e: nil}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/month/%s", testDate),
				},
			},
			expectedBody: `{"code":"Events read success","message":"Events read success",` +
				`"data":[{"id":"ac1b717e-9b4d-4048-b869-ea3be369a5e2","title":"Test title","dateTimeStart":` +
				`"2022-02-19T11:15:00Z","dateTimeEnd":"2022-02-19T11:15:00Z","description":"Test description",` +
				`"createdBy":1,"remindFrom":"2022-02-19T11:15:00Z"}]}`,
			expectedCode: 200,
		},
		{
			name: "EventsOnMonth error not found",
			app: &mockApp{eventsOnMonth: struct {
				c []*calendar.EventDto
				e error
			}{c: nil, e: nil}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/month/%s", testDate),
				},
			},
			expectedBody: `{"code":"Events not found","message":"Events not found","data":null}`,
			expectedCode: 404,
		},
		{
			name: "EventsOnMonth error app",
			app: &mockApp{eventsOnMonth: struct {
				c []*calendar.EventDto
				e error
			}{c: nil, e: errors.New("test error")}},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/month/%s", testDate),
				},
			},
			expectedBody: `{"code":"Error get events","message":"test error","data":null}`,
			expectedCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewMockResponse()
			h := NewHandler(tt.app, logger.NewMockLogger())
			h.ServeHTTP(resp, tt.req)
			require.Equal(t, tt.expectedBody, resp.getBody())
			require.Contains(t, resp.headers["Content-Type"], "application/json")
			require.Equal(t, tt.expectedCode, resp.status)
		})
	}
}

func TestHandler_EventActionCreate(t *testing.T) {
	tests := []serveHTTP{
		{
			name: "Create event success",
			app: &mockApp{
				createEventRes: struct {
					s string
					e error
				}{s: testUUID, e: nil},
			},
			req: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/events",
				},
				Body: io.NopCloser(
					strings.NewReader(`{
									"title": "Test title 1",
									"dateTimeStart": "2022-02-19 12:24:05",
									"dateTimeEnd":"2022-02-19 13:24:05",
									"description":"Test desciption",
									"createdBy":1,
									"remindFrom": "2022-01-10 12:10:05"
							}`)),
			},
			expectedBody: `{"code":"Event created success","message":"Event created ` +
				`with uuid=ac1b717e-9b4d-4048-b869-ea3be369a5e2","data":null}`,
			expectedCode: 201,
		},
		{
			name: "Create event error app",
			app: &mockApp{
				createEventRes: struct {
					s string
					e error
				}{s: "", e: errors.New("test error")},
			},
			req: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/events",
				},
				Body: io.NopCloser(
					strings.NewReader(`{
									"title": "Test title 1",
									"dateTimeStart": "2022-02-19 12:24:05",
									"dateTimeEnd":"2022-02-19 13:24:05",
									"description":"Test desciption",
									"createdBy":1,
									"remindFrom": "2022-01-10 12:10:05"
							}`),
				),
			},
			expectedBody: `{"code":"Error create event","message":"test error","data":null}`,
			expectedCode: 500,
		},
		{
			name: "Create event error validate",
			app: &mockApp{
				createEventRes: struct {
					s string
					e error
				}{s: "", e: errors.New("test error")},
			},
			req: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/events",
				},
				Body: io.NopCloser(
					strings.NewReader(`{
									"title": "Test title 1",
									"dateTimeStart": "",
									"dateTimeEnd":"2022-02-19 13:24:05",
									"description":"Test desciption",
									"createdBy":1,
									"remindFrom": "2022-01-10 12:10:05"
							}`)),
			},
			expectedBody: `{"code":"Validate error","message":"incorrect dateTimeStart","data":null}`,
			expectedCode: 422,
		},
		{
			name: "Create event error incorrect requestBody",
			app: &mockApp{
				createEventRes: struct {
					s string
					e error
				}{s: "", e: errors.New("test error")},
			},
			req: &http.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/events",
				},
				Body: io.NopCloser(strings.NewReader("")),
			},
			expectedBody: `{"code":"Error read request","message":"unexpected end of JSON input","data":null}`,
			expectedCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewMockResponse()
			h := NewHandler(tt.app, logger.NewMockLogger())
			h.ServeHTTP(resp, tt.req)
			require.Equal(t, tt.expectedBody, resp.getBody())
			require.Contains(t, resp.headers["Content-Type"], "application/json")
			require.Equal(t, tt.expectedCode, resp.status)
		})
	}
}

func TestHandler_EventActionGet(t *testing.T) {
	tests := []serveHTTP{
		{
			name: "Get event success",
			app: &mockApp{
				readEventRes: struct {
					c *calendar.EventDto
					e error
				}{c: &calendar.EventDto{
					ID:            testUUID,
					Title:         "Test title",
					Description:   "Test description",
					DateTimeStart: testTime,
					DateTimeEnd:   testTime,
					CreatedBy:     1,
					RemindFrom:    testTime,
				}, e: nil},
			},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
			},
			expectedBody: `{"code":"Event read success","message":"Event read success",` +
				`"data":{"id":"ac1b717e-9b4d-4048-b869-ea3be369a5e2","title":"Test title",` +
				`"dateTimeStart":"2022-02-19T11:15:00Z","dateTimeEnd":"2022-02-19T11:15:00Z",` +
				`"description":"Test description","createdBy":1,"remindFrom":"2022-02-19T11:15:00Z"}}`,
			expectedCode: 200,
		},
		{
			name: "Get event error not found",
			app: &mockApp{
				readEventRes: struct {
					c *calendar.EventDto
					e error
				}{c: nil, e: nil},
			},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
			},
			expectedBody: `{"code":"Event not found","message":"Event not found","data":null}`,
			expectedCode: 404,
		},
		{
			name: "Get event error app",
			app: &mockApp{
				readEventRes: struct {
					c *calendar.EventDto
					e error
				}{c: nil, e: errors.New("test error")},
			},
			req: &http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
			},
			expectedBody: `{"code":"Error get event","message":"test error","data":null}`,
			expectedCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewMockResponse()
			h := NewHandler(tt.app, logger.NewMockLogger())
			h.ServeHTTP(resp, tt.req)
			require.Equal(t, tt.expectedBody, resp.getBody())
			require.Contains(t, resp.headers["Content-Type"], "application/json")
			require.Equal(t, tt.expectedCode, resp.status)
		})
	}
}

func TestHandler_EventAction(t *testing.T) {
	tests := []serveHTTP{
		{
			name: "Delete event success",
			app:  &mockApp{},
			req: &http.Request{
				Method: "DELETE",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
			},
			expectedBody: `{"code":"Event deleted success","message":"Event deleted ` +
				`with uuid=ac1b717e-9b4d-4048-b869-ea3be369a5e2","data":null}`,
			expectedCode: 200,
		},
		{
			name: "Delete event error app",
			app:  &mockApp{deleteEventRes: errors.New("test error")},
			req: &http.Request{
				Method: "DELETE",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
			},
			expectedBody: `{"code":"Error delete event","message":"test error","data":null}`,
			expectedCode: 500,
		},
		{
			name: "Update event success",
			app:  &mockApp{},
			req: &http.Request{
				Method: "PUT",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
				Body: io.NopCloser(
					strings.NewReader(`{
									"title": "Test title 1",
									"dateTimeStart": "2022-02-19 12:24:05",
									"dateTimeEnd":"2022-02-19 13:24:05",
									"description":"Test desciption",
									"createdBy":1,
									"remindFrom": "2022-01-10 12:10:05"
							}`),
				),
			},
			expectedBody: `{"code":"Event updated success","message":"Event updated ` +
				`with uuid=ac1b717e-9b4d-4048-b869-ea3be369a5e2","data":null}`,
			expectedCode: 200,
		},
		{
			name: "Update event error app",
			app:  &mockApp{updateEventRes: errors.New("test error")},
			req: &http.Request{
				Method: "PUT",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
				Body: io.NopCloser(
					strings.NewReader(`{
									"title": "Test title 1",
									"dateTimeStart": "2022-02-19 12:24:05",
									"dateTimeEnd":"2022-02-19 13:24:05",
									"description":"Test desciption",
									"createdBy":1,
									"remindFrom": "2022-01-10 12:10:05"
							}`),
				),
			},
			expectedBody: `{"code":"Error update event","message":"test error","data":null}`,
			expectedCode: 500,
		},
		{
			name: "Update event error incorrect requestBody",
			app:  &mockApp{},
			req: &http.Request{
				Method: "PUT",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
				Body: io.NopCloser(
					strings.NewReader(`{
									"title": "Test title 1",
									"dateTimeStart": "",
									"dateTimeEnd":"2022-02-19 13:24:05",
									"description":"Test desciption",
									"createdBy":1,
									"remindFrom": "2022-01-10 12:10:05"
							}`),
				),
			},
			expectedBody: `{"code":"Validate error","message":"incorrect dateTimeStart","data":null}`,
			expectedCode: 422,
		},
		{
			name: "Event actionNotSupported",
			app:  &mockApp{},
			req: &http.Request{
				Method: "PATCH",
				URL: &url.URL{
					Path: fmt.Sprintf("/events/%s", testUUID),
				},
			},
			expectedBody: `{"code":"Unsupported action","message":"Unsupported action","data":null}`,
			expectedCode: 405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := NewMockResponse()
			h := NewHandler(tt.app, logger.NewMockLogger())
			h.ServeHTTP(resp, tt.req)
			require.Equal(t, tt.expectedBody, resp.getBody())
			require.Contains(t, resp.headers["Content-Type"], "application/json")
			require.Equal(t, tt.expectedCode, resp.status)
		})
	}
}

type serveHTTP struct {
	name         string
	app          calendar.Application
	req          *http.Request
	expectedBody string
	expectedCode int
}
