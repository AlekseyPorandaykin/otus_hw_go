package queue

import "fmt"

var (
	NotConnectError     = NewError("error setup connect to broker")
	NotOpenChannelError = NewError("error open channel broker")
	ClosedError         = NewError("closed connect broker")
	ConsumeError        = NewError("error consume message broker")
)

type Error struct {
	msg string
	err error
}

func (e *Error) Wrap(err error) *Error {
	e.err = err

	return e
}

func (e *Error) Error() string {
	if e.err == nil {
		return e.msg
	}
	return fmt.Sprintf("%s : %s", e.msg, e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}

func NewError(msg string) *Error {
	return &Error{msg: msg}
}
