package goagi

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAGI(t *testing.T) {
	input := "agi_channel: SIP/2222@default-00000023\n" +
		"agi_language: en\n" +
		"agi_type: SIP\n" +
		"agi_uniqueid: 1397044468.0\n" +
		"agi_version: 0.1\n" +
		"agi_callerid: 5001\n" +
		"agi_calleridname: Alice\n" +
		"agi_callingpres: 67\n" +
		"agi_dnid: 123456\n" +
		"agi_context: default\n" +
		"agi_extension: 2222\n" +
		"agi_priority: 1\n" +
		"agi_accountcode: 0\n" +
		"agi_threadid: 140536028174080\n" +
		"agi_arg_2: bar=123\n" +
		"\n"

	ch := make(chan interface{})

	var agi *AGI
	var err error
	go func() {
		defer close(ch)
		agi, err = NewAGI()
	}()
	fmt.Fprintf(os.Stdin, input)
	<-ch
	assert.Nil(t, err)
	assert.NotNil(t, agi)
}

func TestNewFastAGI(t *testing.T) {
	input := "agi_channel: SIP/2222@default-00000023\n" +
		"agi_language: en\n" +
		"agi_type: SIP\n" +
		"agi_uniqueid: 1397044468.0\n" +
		"agi_version: 0.1\n" +
		"agi_callerid: 5001\n" +
		"agi_calleridname: Alice\n" +
		"agi_callingpres: 67\n" +
		"agi_dnid: 123456\n" +
		"agi_context: default\n" +
		"agi_extension: 2222\n" +
		"agi_priority: 1\n" +
		"agi_accountcode: 0\n" +
		"agi_threadid: 140536028174080\n" +
		"agi_arg_2: bar=123\n" +
		"\n"

	fagi, err := NewFastAGI("127.0.0.1:0")

	assert.Nil(t, err)

	conn, err := net.Dial("tcp", fagi.ln.Addr().String())
	assert.Nil(t, err)

	fmt.Fprintf(conn, input)
	agi := <-fagi.Conn()

	assert.Equal(t, "SIP/2222@default-00000023", agi.Env("channel"))
	fmt.Fprintf(conn, "200 result=1\n")
	err = agi.Verbose("Accept new connection.")
	assert.Nil(t, err)

	assert.Nil(t, fagi.Err())
	fagi.Close()
}
