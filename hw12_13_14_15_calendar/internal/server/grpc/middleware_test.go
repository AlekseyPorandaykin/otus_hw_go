package grpc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"google.golang.org/grpc"
)

type mockResponse struct {
	status bool
}

func (r mockResponse) GetStatus() bool {
	return r.status
}

func TestMiddleware_loggingMiddleware(t *testing.T) {
	type fields struct {
		logger logger.Logger
	}
	type args struct {
		ctx     context.Context
		req     interface{}
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	var (
		respTrue  = &mockResponse{true}
		respFalse = &mockResponse{false}
	)
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp interface{}
		wantErr  bool
	}{
		{
			name: "Success test true",
			fields: fields{
				logger: logger.NewMockLogger(),
			},
			args: args{
				ctx:  context.TODO(),
				req:  nil,
				info: &grpc.UnaryServerInfo{FullMethod: "Test"},
				handler: func(ctx context.Context, req interface{}) (interface{}, error) {
					return respTrue, nil
				},
			},
			wantResp: respTrue,
			wantErr:  false,
		},
		{
			name: "Success test false",
			fields: fields{
				logger: logger.NewMockLogger(),
			},
			args: args{
				ctx:  context.TODO(),
				req:  nil,
				info: &grpc.UnaryServerInfo{FullMethod: "Test"},
				handler: func(ctx context.Context, req interface{}) (interface{}, error) {
					return respFalse, nil
				},
			},
			wantResp: respFalse,
			wantErr:  false,
		},
		{
			name: "Success test error",
			fields: fields{
				logger: logger.NewMockLogger(),
			},
			args: args{
				ctx:  context.TODO(),
				req:  nil,
				info: &grpc.UnaryServerInfo{FullMethod: "Test"},
				handler: func(ctx context.Context, req interface{}) (interface{}, error) {
					return respFalse, errors.New("test error")
				},
			},
			wantResp: respFalse,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Middleware{
				logger: tt.fields.logger,
			}
			gotResp, err := m.loggingMiddleware(tt.args.ctx, tt.args.req, tt.args.info, tt.args.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("loggingMiddleware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("loggingMiddleware() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_getStatus(t *testing.T) {
	type args struct {
		resp interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Success test true status",
			args: args{
				resp: &mockResponse{true},
			},
			want: true,
		},
		{
			name: "Success test false status",
			args: args{
				resp: &mockResponse{false},
			},
			want: false,
		},
		{
			name: "Success test incorrect response",
			args: args{
				resp: "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStatus(tt.args.resp); got != tt.want {
				t.Errorf("getStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
