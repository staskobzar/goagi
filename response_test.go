package goagi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRespOkZero(t *testing.T) {
	str := "200 result=0\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 0, r.result)
	assert.Equal(t, "", r.data)
}

func TestRespInvalCode(t *testing.T) {
	str := "a200 result=0\n"
	_, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalResp)
}

func TestRespTooShort(t *testing.T) {
	str := "200\n"
	_, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalResp)
}

func TestRespNoSpaceAfterCode(t *testing.T) {
	str := "200result=1\n"
	_, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalResp)
}

func TestRespInvalidResponseResult(t *testing.T) {
	str := "200 foo=1\n"
	_, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EInvalResp)
}

func TestRespResultMissing(t *testing.T) {
	str := "200 result=\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.EqualValues(t, -3, r.result)
}

func TestRespOkOne(t *testing.T) {
	str := "200 result=1\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 1, r.result)
	assert.Equal(t, "", r.data)
}

func TestRespOkMinusOne(t *testing.T) {
	str := "200 result=-1\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.Equal(t, "", r.data)
}

func TestRespOkWithShortData(t *testing.T) {
	str := "200 result=1 (timeout)\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 1, r.result)
	assert.Equal(t, "(timeout)", r.data)
}

func TestRespOkWithLongData(t *testing.T) {
	str := "200 result=5 (dtmf) endpos=123456\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 200, r.code)
	assert.EqualValues(t, 5, r.result)
	assert.Equal(t, "(dtmf) endpos=123456", r.data)
}

func TestRespError520(t *testing.T) {
	str := "520 Invalid command syntax.  Proper usage not available.\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 520, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.Equal(t, "Invalid command syntax.  Proper usage not available.", r.data)
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
	assert.Equal(t, data, r.data)
}

func TestRespError511(t *testing.T) {
	str := "511 Command Not Permitted on a dead channel\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 511, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.Equal(t, "Command Not Permitted on a dead channel", r.data)
}

func TestRespError510(t *testing.T) {
	str := "510 Invalid or unknown command\n"
	r, err := parseResponse(str)
	assert.Nil(t, err)
	assert.Equal(t, 510, r.code)
	assert.EqualValues(t, -1, r.result)
	assert.Equal(t, "Invalid or unknown command", r.data)
}

func TestRespHangup(t *testing.T) {
	str := "HANGUP\n"
	_, err := parseResponse(str)
	assert.NotNil(t, err)
	assert.Equal(t, err, EHangUp)
}

func TestRespExtractValue(t *testing.T) {
	l := &lexer{input: "(var value)"}
	assert.Equal(t, "var value", l.extractResposeValue())
	l = &lexer{input: ""}
	assert.Equal(t, "", l.extractResposeValue())
	l = &lexer{input: "something wrong"}
	assert.Equal(t, "", l.extractResposeValue())
	l = &lexer{input: "(12454785415.001)"}
	assert.Equal(t, "12454785415.001", l.extractResposeValue())
}

func TestRespExtractEndpos(t *testing.T) {
	l := &lexer{input: "endpos=54568"}
	assert.EqualValues(t, 54568, l.extractEndpos())
	l = &lexer{input: "foo"}
	assert.EqualValues(t, -1, l.extractEndpos())
	l = &lexer{input: "endpos=0"}
	assert.EqualValues(t, 0, l.extractEndpos())
	l = &lexer{input: "endpos="}
	assert.EqualValues(t, -1, l.extractEndpos())
}

func BenchmarkParseAGIResponse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str := "200 result=1 (timeout)\n"
		parseResponse(str)
	}
}
