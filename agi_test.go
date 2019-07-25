package goagi

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"os"
	"testing"
	"time"
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

	ch := make(chan interface{})
	go func() {
		err := NewFastAGI("127.0.0.1:56111", func(agi *AGI) {
			defer close(ch)
			assert.Equal(t, "SIP/2222@default-00000023", agi.Env("channel"))
			err := agi.Verbose("Accept new connection.")
			assert.Nil(t, err)
		})
		assert.Nil(t, err)
	}()

	// this is very ugly way. TODO: find better way to sync listen/dial
	time.Sleep(10 * time.Millisecond)
	conn, err := net.Dial("tcp", "127.0.0.1:56111")
	assert.Nil(t, err)
	fmt.Fprintf(conn, input)
	fmt.Fprintf(conn, "200 result=1\n")
	<-ch
	status, err := bufio.NewReader(conn).ReadString('\n')
	assert.Nil(t, err)
	assert.Equal(t, "VERBOSE \"Accept new connection.\"\n", status)
}
