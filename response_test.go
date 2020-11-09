package goagi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespOkZero(t *testing.T) {
	str := "200 result=0\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 0, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespInvalCode(t *testing.T) {
	str := "a200 result=0\n"
	r, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalResp)
	assert.Equal(t, 0, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespTooShort(t *testing.T) {
	str := "200\n"
	r, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalResp)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespNoSpaceAfterCode(t *testing.T) {
	str := "200result=1\n"
	r, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalResp)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespInvalidResponseResult(t *testing.T) {
	str := "200 foo=1\n"
	r, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalResp)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespResultMissing(t *testing.T) {
	str := "200 result=\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, -3, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespOkOne(t *testing.T) {
	str := "200 result=1\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespOkMinusOne(t *testing.T) {
	str := "200 result=-1\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespOkWithShortData(t *testing.T) {
	str := "200 result=1 (timeout)\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "timeout", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespOkDTMFendpos(t *testing.T) {
	str := "200 result=5 (dtmf) endpos=123456\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 5, r.result)
	assert.Equal(t, "dtmf", r.value)
	assert.EqualValues(t, 123456, r.endpos)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespOkHangupEndpos(t *testing.T) {
	str := "200 result=-1 (hangup) endpos=554687\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.Equal(t, "hangup", r.value)
	assert.EqualValues(t, 554687, r.endpos)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespOkValueEndposAndData(t *testing.T) {
	str := "200 result=1 (fooBar) endpos=55468 Additional message\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 1, r.result)
	assert.Equal(t, "fooBar", r.value)
	assert.EqualValues(t, 55468, r.endpos)
	assert.Equal(t, "Additional message", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespOkDataOnly(t *testing.T) {
	str := "200 result=1 Gosub failed\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 1, r.result)
	assert.Equal(t, "", r.value)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "Gosub failed", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespOk100Trying(t *testing.T) {
	str := "100 result=0 Trying...\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 100, r.code)
	assert.EqualValues(t, 0, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "Trying...", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespError520(t *testing.T) {
	str := "520 Invalid command syntax.  Proper usage not available.\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 520, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "Invalid command syntax.  Proper usage not available.", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespError520Long(t *testing.T) {
	data := "Invalid command syntax.  Proper usage follows:\n" +
		"Usage: DATABASE GET\n" +
		"Retrieves an entry in the Asterisk database for a\n" +
		"given family and key.\n" +
		"Returns 0 if is not set. Returns 1 if \n" +
		"is set and returns the variable in parentheses.\n" +
		"Example return code: 200 result=1 (testvariable)\n"
	str := "520-" + data +
		"520 End of proper usage.\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 520, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, data, r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespError511(t *testing.T) {
	str := "511 Command Not Permitted on a dead channel\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 511, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "Command Not Permitted on a dead channel", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespError510(t *testing.T) {
	str := "510 Invalid or unknown command\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 510, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "Invalid or unknown command", r.data)
	assert.Equal(t, str, r.raw)
}

func TestRespHangup(t *testing.T) {
	str := "HANGUP\n"
	r, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EHangUp)
	assert.EqualValues(t, 1, r.result)
	assert.EqualValues(t, -1, r.endpos)
	assert.Equal(t, "", r.value)
	assert.Equal(t, "", r.data)
	assert.Equal(t, str, r.raw)
}

func BenchmarkParseAGIResponse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str := "200 result=1 (timeout)\n"
		_, err := parseResponse(str)
		if err != nil {
			panic(err)
		}
	}
}
