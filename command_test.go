package goagi

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestCmdAnswerOk(t *testing.T) {
	input := "agi_network: yes\n" +
		"agi_network_script: foo?\n" +
		"agi_request: agi://127.0.0.1/foo?\n" +
		"agi_channel: SIP/2222@default-00000023\n" +
		"\n"
	resp := "200 result=0\n"

	r, w := io.Pipe()
	go func() {
		w.Write([]byte(input))
		reader := bufio.NewReader(r)
		reader.ReadString('\n')
		w.Write([]byte(resp))
	}()
	agi, err := newInterface(r, w)
	assert.Nil(t, err)

	ok, err := agi.Answer()
	assert.Nil(t, err)
	assert.True(t, ok)
}
