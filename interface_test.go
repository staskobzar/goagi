package goagi

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func dummyReadWrite(input string) *bufio.ReadWriter {
	reader := strings.NewReader(input)
	writer := bufio.NewWriter(ioutil.Discard)
	return bufio.NewReadWriter(bufio.NewReader(reader), writer)
}

func dummyReadWriteWError() *bufio.ReadWriter {
	r, w := io.Pipe()
	w.Close()
	return bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))
}

func dummyReadWriteRError() *bufio.ReadWriter {
	r, w := io.Pipe()
	r.Close()
	return bufio.NewReadWriter(bufio.NewReader(r), bufio.NewWriter(w))
}

func TestIFaceCommandCompile(t *testing.T) {
	resp := compileCmd("ANSWER")
	assert.Equal(t, "ANSWER\n", resp)
	resp = compileCmd("CHANNEL STATUS", "SIP/00007-1234")
	assert.Equal(t, "CHANNEL STATUS SIP/00007-1234\n", resp)
	resp = compileCmd("CHANNEL STATUS", "")
	assert.Equal(t, "CHANNEL STATUS \"\"\n", resp)
	resp = compileCmd("CONTROL STREAM FILE", "welcome", "", 3000, 9, "#", "", 1600)
	assert.Equal(t, "CONTROL STREAM FILE welcome \"\" 3000 9 # \"\" 1600\n", resp)
	resp = compileCmd("DATABASE DEL", "user", `"alice"`)
	assert.Equal(t, "DATABASE DEL user \"alice\"\n", resp)
}

func TestIFaceInitSuccessful(t *testing.T) {
	input := "agi_network: yes\n" +
		"agi_network_script: foo?\n" +
		"agi_request: agi://127.0.0.1/foo?\n" +
		"agi_channel: SIP/2222@default-00000023\n" +
		"agi_language: en\n" +
		"agi_type: SIP\n" +
		"agi_uniqueid: 1397044468.0\n" +
		"agi_version: 0.1\n" +
		"agi_callerid: 5001\n" +
		"agi_calleridname: Alice\n" +
		"agi_callingpres: 67\n" +
		"agi_callingani2: 0\n" +
		"agi_callington: 0\n" +
		"agi_callingtns: 0\n" +
		"agi_dnid: 123456\n" +
		"agi_rdnis: unknown\n" +
		"agi_context: default\n" +
		"agi_extension: 2222\n" +
		"agi_priority: 1\n" +
		"agi_enhanced: 0.0\n" +
		"agi_accountcode: 0\n" +
		"agi_threadid: 140536028174080\n" +
		"agi_arg_1: argument1\n" +
		"agi_arg_2: bar=123\n" +
		"agi_arg_3: 3\n" +
		"\n"

	rw := dummyReadWrite(input)
	agi, err := newInterface(rw)

	assert.Nil(t, err)
	assert.Equal(t, agi.Env("network"), "yes")
	assert.Equal(t, agi.Env("network_script"), "foo?")
	assert.Equal(t, agi.Env("request"), "agi://127.0.0.1/foo?")
	assert.Equal(t, agi.Env("channel"), "SIP/2222@default-00000023")
	assert.Equal(t, agi.Env("language"), "en")
	assert.Equal(t, agi.Env("type"), "SIP")
	assert.Equal(t, agi.Env("uniqueid"), "1397044468.0")
	assert.Equal(t, agi.Env("version"), "0.1")
	assert.Equal(t, agi.Env("callerid"), "5001")
	assert.Equal(t, agi.Env("calleridname"), "Alice")
	assert.Equal(t, agi.Env("callingpres"), "67")
	assert.Equal(t, agi.Env("callingani2"), "0")
	assert.Equal(t, agi.Env("callington"), "0")
	assert.Equal(t, agi.Env("callingtns"), "0")
	assert.Equal(t, agi.Env("dnid"), "123456")
	assert.Equal(t, agi.Env("rdnis"), "unknown")
	assert.Equal(t, agi.Env("context"), "default")
	assert.Equal(t, agi.Env("extension"), "2222")
	assert.Equal(t, agi.Env("priority"), "1")
	assert.Equal(t, agi.Env("enhanced"), "0.0")
	assert.Equal(t, agi.Env("accountcode"), "0")
	assert.Equal(t, agi.Env("threadid"), "140536028174080")
	assert.Equal(t, 3, len(agi.EnvArgs()))
	assert.ElementsMatch(t, agi.EnvArgs(), []string{"argument1", "bar=123", "3"})

	// not existing header returns empty string
	assert.Equal(t, agi.Env("unknown"), "")
}

func TestIFaceInitInvalidEnvName(t *testing.T) {
	input := "agi_network: yes\n" +
		"agi_network_script: foo?\n" +
		"not_agi_env: bar\n" +
		"agi_language: en\n" +
		"\n"

	rw := dummyReadWrite(input)
	_, err := newInterface(rw)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalEnv)
	assert.Contains(t, err.Error(), "not_agi_env")
}

func TestIFaceInitInvalidHeader(t *testing.T) {
	input := "agi_network: yes\n" +
		"agi_network_script: foo?\n" +
		"agi_env_no_delim bar\n" +
		"agi_language: en\n" +
		"\n"

	rw := dummyReadWrite(input)
	_, err := newInterface(rw)

	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalEnv)
	assert.Contains(t, err.Error(), "agi_env_no_delim")
}

func TestIFaceInitScannerError(t *testing.T) {
	rw := dummyReadWriteRError()
	_, err := newInterface(rw)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "io: read/write")
}

func TestIFaceReaderOK(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	chstr, err := agi.read()
	str := <-chstr
	assert.Nil(t, <-err)
	assert.Equal(t, "200 result=0\n", str)
}

func TestIFaceReaderFail(t *testing.T) {
	rw := dummyReadWriteRError()
	agi := &AGI{io: rw}
	_, cerr := agi.read()
	err := <-cerr
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "io: read/write")
}

func TestIFaceReaderMultiline(t *testing.T) {
	input := "520-Invalid command syntax.  Proper usage follows:\n" +
		"Usage: DATABASE GET\n" +
		"Example return code: 200 result=1 (testvariable)\n" +
		"520 End of proper usage.\n"
	rw := dummyReadWrite(input)
	agi := &AGI{io: rw}
	chstr, err := agi.read()
	str := <-chstr
	assert.Nil(t, <-err)
	assert.Contains(t, str, "520 End of proper usage.\n")
}

func TestIFaceExecute(t *testing.T) {
	rw := dummyReadWrite("200 result=-1 (dtmf)\n")
	agi := &AGI{io: rw}
	resp, err := agi.execute("ANSWER\n")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.code)
	assert.EqualValues(t, -1, resp.result)
	assert.Equal(t, "dtmf", resp.value)
}

func TestIFaceExecuteFailWrite(t *testing.T) {
	rw := dummyReadWriteWError()
	agi := &AGI{io: rw}
	resp, err := agi.execute("CHANNEL STATUS\n")
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func BenchmarkIFaceExecute(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rw := dummyReadWrite("200 result=0 (timeout)\n")
		agi := &AGI{io: rw}
		agi.execute("ANSWER\n")
	}
}

func BenchmarkAGInterfaceInit(b *testing.B) {
	input := "agi_network: yes\n" +
		"agi_network_script: foo?\n" +
		"agi_request: agi://127.0.0.1/foo?\n" +
		"agi_channel: SIP/2222@default-00000023\n" +
		"agi_language: en\n" +
		"agi_type: SIP\n" +
		"agi_uniqueid: 1397044468.0\n" +
		"agi_version: 0.1\n" +
		"agi_callerid: 5001\n" +
		"agi_calleridname: Alice\n" +
		"agi_callingpres: 67\n" +
		"agi_callingani2: 0\n" +
		"agi_callington: 0\n" +
		"agi_callingtns: 0\n" +
		"agi_dnid: 123456\n" +
		"agi_rdnis: unknown\n" +
		"agi_context: default\n" +
		"agi_extension: 2222\n" +
		"agi_priority: 1\n" +
		"agi_enhanced: 0.0\n" +
		"agi_accountcode: 0\n" +
		"agi_threadid: 140536028174080\n" +
		"agi_arg_1: argument1\n" +
		"agi_arg_2: bar=123\n" +
		"agi_arg_3: 3\n" +
		"\n"

	for i := 0; i < b.N; i++ {
		rw := dummyReadWrite(input)
		agi, _ := newInterface(rw)
		agi.Env("network")
	}
}
