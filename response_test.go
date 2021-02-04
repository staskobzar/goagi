package goagi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionSetup(t *testing.T) {
	agi := &AGI{}

	envLen := len(agiSetupInput) - 2
	agi.sessionSetup(agiSetupInput)
	assert.Equal(t, envLen, len(agi.env))
	assert.Equal(t, 2, len(agi.arg))

	input := agiSetupInput[:5]
	input = append(input, "fooo: bar")
	agi.sessionSetup(input)
	assert.Equal(t, 5, len(agi.env))

	input = agiSetupInput[:5]
	input = append(input, "agi_foo bar")
	agi.sessionSetup(input)
	assert.Equal(t, 5, len(agi.env))

	input = agiSetupInput[:5]
	input = append(input, "agi_foo: bar")
	agi.sessionSetup(input)
	assert.Equal(t, 6, len(agi.env))
}

func TestParseResponse(t *testing.T) {
	tests := []struct {
		input  string
		code   int
		result int
		value  string
		data   string
		endpos int64
		digit  string
		sres   int
	}{
		{"100 result=1 Trying...\n", 100, 1, "Trying...", "Trying...", 0, "", 0},
		{"100 result=0\n", 100, 0, "", "", 0, "", 0},
		{"100 Trying\n", 100, 0, "Trying", "Trying", 0, "", 0},
		{"200 result=1\n", 200, 1, "", "", 0, "", 0},
		{"200 result=\n", 200, 0, "", "", 0, "", 0},

		{"200 result=1 (hangup)\n", 200, 1, "hangup", "", 0, "", 0},
		{"200 result=-1 endpos=11223344\n", 200, -1, "", "", 11223344, "", 0},
		{"200 result=-1 endpos=\n", 200, -1, "", "", 0, "", 0},
		{"200 result=-1 endpos=asf\n", 200, -1, "", "", 0, "", 0},
		{"200 result=1 (en)\n", 200, 1, "en", "", 0, "", 0},
		{
			"200 result=1 (\"Alice Johnson\" <2233>)\n",
			200, 1, "\"Alice Johnson\" <2233>", "", 0, "", 0,
		}, {
			"200 result=1 (Alice Johnson)\n",
			200, 1, "Alice Johnson", "", 0, "", 0,
		}, {
			"200 result=1 (SIP/9170-12-00000008)\n",
			200, 1, "SIP/9170-12-00000008", "", 0, "", 0,
		}, {
			"200 result=1 (digit) digit=* endpos=998877660\n",
			200, 1, "digit", "", 998877660, "*", 0,
		}, {
			"503 result=-2 Memory allocation failure\n",
			503, -2,
			"Memory allocation failure", "Memory allocation failure",
			0, "", 0,
		}, {
			"200 result=1 (speech) endpos=918273 results=123 \n",
			200, 1, "speech", "", 918273, "", 123,
		}, {
			"510 Invalid or unknown command\n",
			510, 0,
			"Invalid or unknown command", "Invalid or unknown command",
			0, "", 0,
		}, {
			"510 Error", 510, 0, "Error", "Error", 0, "", 0,
		}, {
			"511 Command Not Permitted on a dead channel or intercept routine\n",
			511, 0,
			"Command Not Permitted on a dead channel or intercept routine",
			"Command Not Permitted on a dead channel or intercept routine",
			0, "", 0,
		}, {
			"520 Invalid command syntax.  Proper usage not available.\n",
			520, 0,
			"Invalid command syntax.  Proper usage not available.",
			"Invalid command syntax.  Proper usage not available.",
			0, "", 0,
		}, {
			"520-Invalid command syntax.  Proper usage follows:\n" +
				"Usage: database put <family> <key> <value>\n" +
				"Adds or updates an entry in the Asterisk database for\n" +
				"a given family, key, and value.\n" +
				"520 End of proper usage.\n",
			520, 0,
			"Usage: database put <family> <key> <value>\n" +
				"Adds or updates an entry in the Asterisk database for\n" +
				"a given family, key, and value.",
			"Usage: database put <family> <key> <value>\n" +
				"Adds or updates an entry in the Asterisk database for\n" +
				"a given family, key, and value.",
			0, "", 0,
		},
	}
	agi := &AGI{}

	for _, tc := range tests {
		resp, err := agi.parseResponse(tc.input, tc.code)
		assert.Nil(t, err)
		assert.Equal(t, tc.code, resp.Code())
		assert.Equal(t, tc.result, resp.Result())
		assert.Equal(t, tc.value, resp.Value())
		assert.Equal(t, tc.data, resp.Data())
		assert.Equal(t, tc.input, resp.RawResponse())
		assert.Equal(t, tc.endpos, resp.EndPos())
		assert.Equal(t, tc.digit, resp.Digit())
		assert.Equal(t, tc.sres, resp.SResults())
	}
}
func TestParseResponseProblems(t *testing.T) {
	tests := []string{
		"", "510", "620 FOO", "30000",
	}
	agi := &AGI{}

	for _, test := range tests {
		_, err := agi.parseResponse(test, 0)
		assert.NotNil(t, err)
	}
	resp, err := agi.parseResponse("200 result=1 (hangup\n", 200)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Empty(t, resp.Value())

	resp, err = agi.parseResponse("200 result=1 (dtmf) endpos=3 results=q \n", 200)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "dtmf", resp.Value())
	assert.Equal(t, 0, resp.SResults())
}

func TestScanResult(t *testing.T) {
	tests := []struct {
		input  string
		remain string
		result int
	}{
		{"result=0 Trying...\n", "Trying...\n", 0},
		{"Command Not Premitted", "Command Not Premitted", 0},
		{"result=2", "", 2},
		{"result=", "", 0},
		{"result= More data", "More data", 0},
		{"result=-1 (timeout) Foo", "(timeout) Foo", -1},
		{"result=*45#3 some data", "some data", 0},
	}
	for _, test := range tests {
		remain, result := scanResult(test.input)
		assert.Equal(t, test.remain, remain, test.input)
		assert.Equal(t, test.result, result, test.input)
	}
}
