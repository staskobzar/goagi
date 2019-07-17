package goagi

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestCmdAnswerOk(t *testing.T) {
	resp := "200 result=0\n"

	r, w := io.Pipe()
	agi := &AGI{input: r, output: w}
	go func() {
		b := make([]byte, 1024)
		r.Read(b)
		w.Write([]byte(resp))
	}()

	ok, err := agi.Answer()
	assert.Nil(t, err)
	assert.True(t, ok)
}
