package goagi

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubReader struct {
	io.Reader
}

type stubWriter struct {
	io.Writer
}

var agiSetupInput = []string{
	"agi_network: yes",
	"agi_network_script: foo?",
	"agi_request: agi://127.0.0.1/foo?",
	"agi_channel: SIP/2222@default-00000023",
	"agi_language: en",
	"agi_type: SIP",
	"agi_uniqueid: 1397044468.0",
	"agi_version: 0.1",
	"agi_callerid: 5001",
	"agi_calleridname: Alice",
	"agi_callingpres: 67",
	"agi_callingani2: 0",
	"agi_callington: 0",
	"agi_callingtns: 0",
	"agi_dnid: 123456",
	"agi_rdnis: unknown",
	"agi_context: default",
	"agi_extension: 2222",
	"agi_priority: 1",
	"agi_enhanced: 0.0",
	"agi_accountcode: 0",
	"agi_threadid: 140536028174080",
	"agi_arg_1: argument1",
	"agi_arg_2: argument2",
}

func TestNew(t *testing.T) {
	input := strings.Join(agiSetupInput, "\n")
	input += "\n\n"
	logBuffer := new(bytes.Buffer)
	logger := log.New(logBuffer, "agi: ", log.Lmicroseconds)
	reader := &stubReader{strings.NewReader(input)}
	writer := &stubWriter{ioutil.Discard}

	agi, err := New(reader, writer, logger)
	assert.Nil(t, err)
	assert.Equal(t, "SIP/2222@default-00000023", agi.Env("channel"))
	assert.Equal(t, "2222", agi.Env("extension"))
	assert.Empty(t, agi.Env("invalid_name"))
	assert.Equal(t, 2, len(agi.EnvArgs()))
	assert.Contains(t, logBuffer.String(), "agi_network: yes")
	assert.Contains(t, logBuffer.String(), "agi_threadid: 140536028174080")
}

func TestNewFail(t *testing.T) {
	reader, writer, err := os.Pipe()
	assert.Nil(t, err)

	reader.Close()
	agi, err := New(reader, writer, nil)
	assert.Nil(t, agi)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "closed")
}

func TestNewIOPipe(t *testing.T) {
	reader, writer, err := os.Pipe()
	assert.Nil(t, err)
	input := strings.Join(agiSetupInput, "\n")
	input += "\n\n"

	writer.WriteString(input)
	agi, err := New(reader, writer, nil)
	assert.Nil(t, err)
	assert.Equal(t, "2222", agi.Env("extension"))
	assert.Equal(t, 2, len(agi.EnvArgs()))
}

func TestNewNetPipe(t *testing.T) {
	reader, writer := net.Pipe()
	input := strings.Join(agiSetupInput, "\n")
	input += "\n\n"

	go writer.Write([]byte(input))
	agi, err := New(reader, writer, nil)
	assert.Nil(t, err)
	assert.Equal(t, "2222", agi.Env("extension"))
	assert.Equal(t, 2, len(agi.EnvArgs()))
}

func TestSessionInitDeviceFail(t *testing.T) {
	stdin, _, err := os.Pipe()
	assert.Nil(t, err)
	agi := &AGI{reader: stdin}
	stdin.Close()
	_, err = agi.sessionInit()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "closed")
}

func TestReadResponse(t *testing.T) {
	tests := []struct {
		input string
		code  int
		ishup bool
	}{
		{"100 result=0 Trying...\n", codeEarly, false},
		{"200 result=1\n", codeSucc, false},
		{"503 result=-2 Memory allocation failure\n", codeE503, false},
		{"510 Invalid or unknown command\n", codeE510, false},
		{"511 Command Not Premitted\n", codeE511, false},
		{"520 Invalid command syntax.\n", codeE520, false},
	}

	// run tests
	for _, test := range tests {
		reader := &stubReader{strings.NewReader(test.input)}
		agi := &AGI{reader: reader}

		res, code, err := agi.read()
		assert.Nil(t, err, test.input)
		assert.Equal(t, test.input, res)
		assert.Equal(t, test.code, code, test.input)
		assert.Equal(t, test.ishup, agi.isHUP, test.input)
	}
}

func TestReadResponseWithHangup(t *testing.T) {
	reader := &stubReader{strings.NewReader("HANGUP\n200 result=1\n")}
	agi := &AGI{reader: reader}
	res, code, err := agi.read()
	assert.Nil(t, err)
	assert.Equal(t, "200 result=1\n", res)
	assert.Equal(t, codeSucc, code)
	assert.True(t, agi.IsHungup())
}

func TestReadResponseLongError520(t *testing.T) {
	input := "HANGUP\n520-Invalid command syntax.  Proper usage follows:\n" +
		"Usage: DATABASE GET\n" +
		"Retrieves an entry in the Asterisk database for a\n" +
		"given family and key.\n" +
		"Returns 0 if is not set. Returns 1 if \n" +
		"is set and returns the variable in parentheses.\n" +
		"Example return code: 200 result=1 (testvariable)\n" +
		"520 End of proper usage.\n"
	reader := &stubReader{strings.NewReader(input)}
	agi := &AGI{reader: reader}
	res, code, err := agi.read()
	assert.Nil(t, err)
	assert.Equal(t, input[7:], res)
	assert.Equal(t, codeE520, code)
	assert.True(t, agi.IsHungup())
}

func TestReadResponseGarbage(t *testing.T) {
	tests := []string{
		"Usage: DATABASE GET\n" +
			"Retrieves an entry in the Asterisk database for a\n" +
			"given family and key.\n" +
			"Returns 0 if is not set. Returns 1 if \n" +
			"is set and returns the variable in parentheses.\n" +
			"Example return code: 200 result=1 (testvariable)\n" +
			"520 End of proper usage.\n",
		"\n", // empty string
		"2222222\n",
		"-1 result=4\n",
		"780 result=0\n",
	}

	for _, input := range tests {
		reader := &stubReader{strings.NewReader(input)}
		agi := &AGI{reader: reader}
		_, code, err := agi.read()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Invalid input")
		assert.Zero(t, code)
	}
}

func TestReadResponseFail(t *testing.T) {
	stdin, _, err := os.Pipe()
	assert.Nil(t, err)
	agi := &AGI{reader: stdin}
	stdin.Close()
	_, _, err = agi.read()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "closed")
}

func TestWriteStdIO(t *testing.T) {
	reader, writer, err := os.Pipe()
	assert.Nil(t, err)
	agi := &AGI{writer: writer}

	err = agi.write([]byte("NOOP\n"))
	assert.Nil(t, err)

	buf := make([]byte, 32)
	n, _ := reader.Read(buf)
	assert.Equal(t, "NOOP\n", string(buf[:n]))
}

func TestWriteNetConn(t *testing.T) {
	reader, writer := net.Pipe()

	agi := &AGI{writer: writer}

	go func() {
		agi.write([]byte("ANSWER\n"))
		writer.Close()
	}()

	buf, _ := ioutil.ReadAll(reader)
	assert.Equal(t, "ANSWER\n", string(buf))
}

func TestWriteFail(t *testing.T) {
	reader, writer, err := os.Pipe()
	assert.Nil(t, err)
	agi := &AGI{writer: writer}

	reader.Close()
	err = agi.write([]byte("HANGUP\n"))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "broken pipe")
}

func TestExecute(t *testing.T) {
	buf, reader, writer := stubReaderWriter("200 result=1")
	agi := &AGI{reader: reader, writer: writer}
	resp, err := agi.execute("NOOP\n")
	assert.Nil(t, err)
	assert.Equal(t, "NOOP\n", buf.String())
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, 1, resp.Result())

	resp, err = agi.execute("HANGUP\n")
	assert.NotNil(t, err)
	assert.Nil(t, resp)

	_, wr := net.Pipe()
	agi.writer = wr
	wr.Close()
	resp, err = agi.execute("HANGUP\n")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "closed")
}
