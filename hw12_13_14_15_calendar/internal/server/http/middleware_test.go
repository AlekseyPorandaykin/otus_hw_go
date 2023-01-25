package internalhttp

import (
	"net/http"
	"strings"
	"testing"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/stretchr/testify/require"
)

type MockResponse struct {
	headers http.Header
	body    strings.Builder
	status  int
}

func NewMockResponse() *MockResponse {
	return &MockResponse{
		headers: make(http.Header),
		body:    strings.Builder{},
	}
}

func (m *MockResponse) Header() http.Header {
	return m.headers
}

func (m *MockResponse) Write(bytes []byte) (int, error) {
	return m.body.Write(bytes)
}

func (m *MockResponse) WriteHeader(statusCode int) {
	m.status = statusCode
}

func (m *MockResponse) getBody() string {
	return m.body.String()
}

func TestMiddleware_loggingMiddleware(t *testing.T) {
	type fields struct {
		logger *logger.MockLogger
	}
	type args struct {
		next http.HandlerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Correct execute",
			fields: fields{
				logger: logger.NewMockLogger(),
			},
			args: args{
				next: func(writer http.ResponseWriter, request *http.Request) {
					writer.WriteHeader(200)
					writer.Write([]byte("test"))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Middleware{
				logger: tt.fields.logger,
			}
			resp := NewMockResponse()
			headers := make(http.Header)
			headers["User-Agent"] = []string{"test"}
			req := &http.Request{
				Method:     "POST",
				RemoteAddr: "0.0.0.0:80",
				RequestURI: "/list",
				Header:     headers,
				Proto:      "HTTP/1.1",
			}

			m.loggingMiddleware(tt.args.next)(resp, req)
			require.True(t, tt.fields.logger.HasMessage("info", "HTTP request", "ip", "0.0.0.0"))
			require.True(t, tt.fields.logger.HasMessage("info", "HTTP request", "method", "POST"))
			require.True(t, tt.fields.logger.HasMessage("info", "HTTP request", "path", "/list"))
			require.True(t, tt.fields.logger.HasMessage("info", "HTTP request", "user-agent", "test"))
			require.True(t, tt.fields.logger.HasMessage("info", "HTTP request", "protocol", "HTTP/1.1"))
			require.True(t, tt.fields.logger.HasMessage("info", "HTTP request", "http-status", "200"))
		})
	}
}
