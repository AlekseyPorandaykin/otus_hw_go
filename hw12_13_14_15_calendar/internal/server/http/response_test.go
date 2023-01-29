package internalhttp

import (
	"reflect"
	"testing"
)

func TestNewResponse(t *testing.T) {
	type args struct {
		code    string
		message string
		data    interface{}
	}
	tests := []struct {
		name string
		args args
		want *Response
	}{
		{
			name: "Correct response with message",
			args: args{
				code:    eventNotFoundCode,
				message: "Test message",
				data:    "data",
			},
			want: &Response{
				Code:    eventNotFoundCode,
				Message: "Test message",
				Data:    "data",
			},
		},
		{
			name: "Correct response without message",
			args: args{
				code:    eventNotFoundCode,
				message: "",
				data:    "data",
			},
			want: &Response{
				Code:    eventNotFoundCode,
				Message: eventNotFoundCode,
				Data:    "data",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponse(tt.args.code, tt.args.message, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_getStatus(t *testing.T) {
	tests := []struct {
		name string
		want int
		Code string
	}{
		{
			name: "Test validate",
			Code: validateErrorCode,
			want: 422,
		},
		{
			name: "Test errorCreateEvent",
			Code: errorCreateEvent,
			want: 500,
		},
		{
			name: "Test errorGetEvent",
			Code: errorGetEvent,
			want: 500,
		},
		{
			name: "Test eventDeletedErrorCode",
			Code: eventDeletedErrorCode,
			want: 500,
		},
		{
			name: "Test eventUpdateErrorCode",
			Code: eventUpdateErrorCode,
			want: 500,
		},
		{
			name: "Test errorGetEventsCode",
			Code: errorGetEventsCode,
			want: 500,
		},
		{
			name: "Test eventCreateSuccessCode",
			Code: eventCreateSuccessCode,
			want: 201,
		},
		{
			name: "Test eventNotFoundCode",
			Code: eventNotFoundCode,
			want: 404,
		},
		{
			name: "Test handlerNotFound",
			Code: handlerNotFound,
			want: 404,
		},
		{
			name: "Test eventsNotFoundCode",
			Code: eventsNotFoundCode,
			want: 404,
		},
		{
			name: "Test eventAlreadyExistCode",
			Code: eventAlreadyExistCode,
			want: 409,
		},
		{
			name: "Test eventDeletedSuccessCode",
			Code: eventDeletedSuccessCode,
			want: 200,
		},
		{
			name: "Test eventReadSuccessCode",
			Code: eventReadSuccessCode,
			want: 200,
		},
		{
			name: "Test eventUpdateSuccessCode",
			Code: eventUpdateSuccessCode,
			want: 200,
		},
		{
			name: "Test eventsReadSuccessCode",
			Code: eventsReadSuccessCode,
			want: 200,
		},
		{
			name: "Test unsupportedActionCode",
			Code: unsupportedActionCode,
			want: 405,
		},
		{
			name: "Test eventAlreadyDeleted",
			Code: eventAlreadyDeleted,
			want: 410,
		},
		{
			name: "Test errorReadRequest",
			Code: errorReadRequest,
			want: 400,
		},
		{
			name: "Test unsupported code",
			Code: "other code",
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Response{
				Code: tt.Code,
			}
			if got := r.getStatus(); got != tt.want {
				t.Errorf("getStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
