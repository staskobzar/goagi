package goagi

import "fmt"

type Error struct {
	s string
	e string
}

func newError(ctx string) *Error {
	return &Error{s: ctx}
}

func (e *Error) Msg(msg string, args ...interface{}) error {
	e.e = fmt.Sprintf(msg, args...)
	return e
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.s, e.e)
}
