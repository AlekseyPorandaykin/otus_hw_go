package internalhttp

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventRequest_GetDateTimeEnd(t *testing.T) {
	type fields struct {
		Title         string
		DateTimeStart string
		DateTimeEnd   string
		Description   string
		CreatedBy     int32
		RemindFrom    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name: "Correct DateTimeEnd",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			want:    time.Date(2022, 2, 19, 12, 24, 5, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Error DateTimeEnd",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:15:00",
				DateTimeEnd:   "2022-02-19--12:24:05",
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventRequest{
				Title:         tt.fields.Title,
				DateTimeStart: tt.fields.DateTimeStart,
				DateTimeEnd:   tt.fields.DateTimeEnd,
				Description:   tt.fields.Description,
				CreatedBy:     tt.fields.CreatedBy,
				RemindFrom:    tt.fields.RemindFrom,
			}
			got, err := e.GetDateTimeEnd()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDateTimeEnd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, fmt.Sprintf("GetDateTimeEnd() got = %v, want %v", got, tt.want))
		})
	}
}

func TestEventRequest_GetDateTimeStart(t *testing.T) {
	type fields struct {
		Title         string
		DateTimeStart string
		DateTimeEnd   string
		Description   string
		CreatedBy     int32
		RemindFrom    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name: "Correct DateTimeStart",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			want:    time.Date(2022, 2, 19, 11, 15, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Error DateTimeStart",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19-fsd11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventRequest{
				Title:         tt.fields.Title,
				DateTimeStart: tt.fields.DateTimeStart,
				DateTimeEnd:   tt.fields.DateTimeEnd,
				Description:   tt.fields.Description,
				CreatedBy:     tt.fields.CreatedBy,
				RemindFrom:    tt.fields.RemindFrom,
			}
			got, err := e.GetDateTimeStart()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDateTimeStart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDateTimeStart() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventRequest_GetRemindFrom(t *testing.T) {
	type fields struct {
		Title         string
		DateTimeStart string
		DateTimeEnd   string
		Description   string
		CreatedBy     int32
		RemindFrom    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		{
			name: "Correct RemindFrom",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			want:    time.Date(2022, 2, 18, 12, 0o0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Error RemindFrom",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19-fsd11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18fds 12:00:00",
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventRequest{
				Title:         tt.fields.Title,
				DateTimeStart: tt.fields.DateTimeStart,
				DateTimeEnd:   tt.fields.DateTimeEnd,
				Description:   tt.fields.Description,
				CreatedBy:     tt.fields.CreatedBy,
				RemindFrom:    tt.fields.RemindFrom,
			}
			got, err := e.GetRemindFrom()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemindFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRemindFrom() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventRequest_Validate(t *testing.T) {
	type fields struct {
		Title         string
		DateTimeStart string
		DateTimeEnd   string
		Description   string
		CreatedBy     int32
		RemindFrom    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		errMsg  string
	}{
		{
			name: "Correct validate",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			wantErr: false,
			errMsg:  "",
		}, {
			name: "Incorrect title validate",
			fields: fields{
				Title:         "",
				DateTimeStart: "2022-02-19 11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			wantErr: true,
			errMsg:  "empty title",
		}, {
			name: "Incorrect description validate",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "",
				CreatedBy:     2,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			wantErr: true,
			errMsg:  "empty description",
		}, {
			name: "Incorrect createdBy validate",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:15:00",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     0,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			wantErr: true,
			errMsg:  "empty createdBy",
		}, {
			name: "Incorrect dateTimeStart validate",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "",
				DateTimeEnd:   "2022-02-19 12:24:05",
				Description:   "Test description",
				CreatedBy:     1,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			wantErr: true,
			errMsg:  "incorrect dateTimeStart",
		}, {
			name: "Incorrect dateTimeEnd validate",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:24:05",
				DateTimeEnd:   "",
				Description:   "Test description",
				CreatedBy:     1,
				RemindFrom:    "2022-02-18 12:00:00",
			},
			wantErr: true,
			errMsg:  "incorrect dateTimeEnd",
		}, {
			name: "Incorrect remindFrom validate",
			fields: fields{
				Title:         "Test title",
				DateTimeStart: "2022-02-19 11:24:05",
				DateTimeEnd:   "2022-02-19 11:24:05",
				Description:   "Test description",
				CreatedBy:     1,
				RemindFrom:    "",
			},
			wantErr: true,
			errMsg:  "incorrect remindFrom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventRequest{
				Title:         tt.fields.Title,
				DateTimeStart: tt.fields.DateTimeStart,
				DateTimeEnd:   tt.fields.DateTimeEnd,
				Description:   tt.fields.Description,
				CreatedBy:     tt.fields.CreatedBy,
				RemindFrom:    tt.fields.RemindFrom,
			}
			err := e.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				require.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func Test_toEventDto(t *testing.T) {
	type args struct {
		req *EventRequest
	}
	tests := []struct {
		name string
		args args
		want *calendar.EventDto
	}{
		{
			name: "Correct eventDto",
			args: args{
				req: &EventRequest{
					Title:         "Test title",
					DateTimeStart: "2022-02-19 11:15:00",
					DateTimeEnd:   "2022-02-19 12:24:05",
					Description:   "Test description",
					CreatedBy:     2,
					RemindFrom:    "2022-02-18 12:00:00",
				},
			},
			want: &calendar.EventDto{
				Title:         "Test title",
				DateTimeStart: time.Date(2022, 2, 19, 11, 15, 0, 0, time.UTC),
				DateTimeEnd:   time.Date(2022, 2, 19, 12, 24, 5, 0, time.UTC),
				Description:   "Test description",
				CreatedBy:     2,
				RemindFrom:    time.Date(2022, 2, 18, 12, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := toEventDto(tt.args.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toEventDto() = %v, want %v", got, tt.want)
			}
		})
	}
}
