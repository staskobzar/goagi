package goagi

import "fmt"

// Error object for goagi library
type Error struct {
	s string
	e string
}

func newError(ctx string) *Error {
	return &Error{s: ctx}
}

// Msg append message to main context message
func (e *Error) Msg(msg string, args ...interface{}) error {
	e.e = fmt.Sprintf(msg, args...)
	return e
}

// Error messge for the Error object
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.s, e.e)
}
