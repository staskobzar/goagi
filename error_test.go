package goagi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorNew(t *testing.T) {
	err := newError("Foo")
	err.Msg("bar")
	assert.Equal(t, "Foo", err.s)
	assert.Equal(t, "bar", err.e)
}

func TestErrorMessage(t *testing.T) {
	NetErr := newError("Network Error")

	err := NetErr.Msg("connection lost")

	assert.Equal(t, "Network Error: connection lost", err.Error())
}

func TestErrorMessageArgs(t *testing.T) {
	ArgErr := newError("EArg")

	err := ArgErr.Msg("has invalid value %d", -5)

	assert.Equal(t, "EArg: has invalid value -5", err.Error())
	assert.Equal(t, ArgErr, err)
}
