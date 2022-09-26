package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/fixme_my_friend/hw09_struct_validator/rules"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Complex struct {
		Name string   `validate:"regexp:\\d+|len:5"`
		Res  Response `validate:"nested"`
	}

	SliceValues struct {
		A []string `validate:"len:3"`
		B []int    `validate:"in:100,200,300,400"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: &User{
				ID:     "9ac87c60-352d-11ed-a1fe-baed42207ac3",
				Name:   "Ivan",
				Age:    23,
				Email:  "ivan@email.com",
				Role:   "admin",
				Phones: []string{"18294593857"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in:          &App{Version: "param"},
			expectedErr: nil,
		},
		{
			in:          &Token{Header: []byte("one"), Payload: []byte("two"), Signature: []byte("three")},
			expectedErr: nil,
		},
		{
			in:          &Response{Code: 200, Body: "test message"},
			expectedErr: nil,
		},
		{
			in: &SliceValues{
				A: []string{"value", "from", "there"},
				B: []int{100, 200, 400},
			},
			expectedErr: nil,
		},
		{
			in: &User{
				ID:     "9ac87c60-352d-11ed-a1fe-baed42207ac3",
				Name:   "Ivan",
				Age:    10,
				Email:  "ivan@email.com",
				Role:   "admin",
				Phones: []string{"18294593857"},
				meta:   nil,
			},
			expectedErr: rules.ErrMinValueLess,
		},
		{
			in: &User{
				ID:     "9ac87c60-352d-11ed-a1fe-baed42207ac3",
				Name:   "Ivan",
				Age:    56,
				Email:  "ivan@email.com",
				Role:   "admin",
				Phones: []string{"18294593857"},
				meta:   nil,
			},
			expectedErr: rules.ErrMaxValueLarge,
		},
		{
			in: &User{
				ID:     "9ac87c60-352d-11ed-a1fe-baed42207ac3",
				Name:   "Ivan",
				Age:    56,
				Email:  "ivan@.com",
				Role:   "admin",
				Phones: []string{"18294593857"},
				meta:   nil,
			},
			expectedErr: rules.ErrRegexpNotValidRule,
		},
		{
			in: &User{
				ID:     "9ac87c60-352d-11ed-a1fe-baed42207ac3",
				Name:   "Ivan",
				Age:    56,
				Email:  "ivan@mail.com",
				Role:   "guest",
				Phones: []string{"18294593857"},
				meta:   nil,
			},
			expectedErr: rules.ErrInNotInRange,
		},
		{
			in: &User{
				ID:     "9ac87c60-352d-11ed-a1fe-baed42207ac3",
				Name:   "Ivan",
				Age:    56,
				Email:  "ivan@mail.com",
				Role:   "guest",
				Phones: []string{"182957"},
				meta:   nil,
			},
			expectedErr: rules.ErrLenStringValue,
		},
		{
			in:          &Response{Code: 209, Body: "test message"},
			expectedErr: rules.ErrInNotInRange,
		},
		{
			in:          &Complex{Name: "10050", Res: Response{Code: 210, Body: "test message"}},
			expectedErr: rules.ErrInNotInRange,
		},
		{
			in: &SliceValues{
				A: []string{"value", "from", "the"},
				B: []int{100, 200, 400},
			},
			expectedErr: rules.ErrLenStringValue,
		},
		{
			in: &SliceValues{
				A: []string{"value", "from", "there"},
				B: []int{10, 200, 400},
			},
			expectedErr: rules.ErrLenStringValue,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			errValidate := Validate(tt.in)
			if tt.expectedErr != nil {
				require.ErrorIs(t, errValidate, tt.expectedErr)
			}
		})
	}
}
