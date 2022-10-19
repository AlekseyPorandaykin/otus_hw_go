package hw09structvalidator

import (
	"errors"
	"fmt"
)

var (
	ErrVarNotStruct = errors.New(`value not struct`)
)

type ValidationError struct {
	NestedError ValidationErrors
	Field       string
	Err         error
}

func (e *ValidationError) Wrap(err ValidationError) {
	e.NestedError = append(e.NestedError, err)
}

func (e *ValidationError) Unwrap() error {
	return e.NestedError
}

func (e *ValidationError) Is(target error) bool {
	if errors.Is(e.Err, target) {
		return true
	}
	for _, err := range e.NestedError {
		if errors.Is(err.Err, target) {
			return true
		}
	}
	return false
}

type ValidationErrors []ValidationError

//TODO need rewrite
func (v ValidationErrors) Error() string {
	errStr := ""
	for _, err := range v {
		errStr += fmt.Sprintf("Field: %s, Error: %s %s", err.Field, err.Err, EOF)
	}
	return errStr
}

func (v ValidationErrors) Is(target error) bool {
	for _, e := range v {
		if e.Is(target) {
			return true
		}
	}
	return false
}
