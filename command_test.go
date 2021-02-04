package goagi

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stubReaderWriter(response string) (*bytes.Buffer, *stubReader, *stubWriter) {
	buf := new(bytes.Buffer)
	response += "\n"
	reader := &stubReader{strings.NewReader(response), 0}
	writer := &stubWriter{buf, 0}

	return buf, reader, writer
}

func mockAGI(response string) (*AGI, *bytes.Buffer) {
	buf, r, w := stubReaderWriter(response)
	return &AGI{reader: r, writer: w, rwtout: rwDefaultTimeout}, buf
}

const respOk = "200 result=1"

func TestCmdCommand(t *testing.T) {
	buf, r, w := stubReaderWriter(respOk)
	agi := &AGI{reader: r, writer: w, rwtout: rwDefaultTimeout * 5}

	resp, err := agi.Command("ANSWER")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.EqualValues(t, 0, r.tout, "No timeout when reading")
	assert.Equal(t, "ANSWER\n", buf.String())
}

func TestCmdAnswer(t *testing.T) {
	agi, buf := mockAGI("511 Command Not Permitted")
	resp, err := agi.Answer()
	assert.Nil(t, err)
	assert.Equal(t, 511, resp.Code())
	assert.Equal(t, "ANSWER\n", buf.String())
}

func TestCmdAsyncAGIBreak(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.AsyncAGIBreak()
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, "ASYNCAGI BREAK\n", buf.String())
}

func TestCmdChannelStatus(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.ChannelStatus("SIP/0001-FA878E0")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, "CHANNEL STATUS SIP/0001-FA878E0\n", buf.String())
}

func TestCmdControlStreamFile(t *testing.T) {
	agi, buf := mockAGI("200 result=1 endpos=10")
	resp, err := agi.ControlStreamFile("welcome", "")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.EqualValues(t, 10, resp.EndPos())
	assert.Equal(t, "CONTROL STREAM FILE welcome \"\"\n", buf.String())

	agi, buf = mockAGI("200 result=1 endpos=998877")
	resp, _ = agi.ControlStreamFile("welcome", "123", "1500")
	assert.EqualValues(t, 998877, resp.EndPos())
	assert.Equal(t, `CONTROL STREAM FILE welcome "123" "1500"`+"\n", buf.String())

	agi, buf = mockAGI("200 result=1 endpos=10")
	resp, _ = agi.ControlStreamFile("welcome", "123", "1500", "#", "0", "*", "1600")
	assert.Equal(t, `CONTROL STREAM FILE welcome "123" "1500" "#" "0" "*" "1600"`+"\n",
		buf.String())

	_, err = agi.ControlStreamFile("welcome", "123", "1500", "#", "0", "*", "1600", "a", "1")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "[a 1]")
}

func TestCmdDatabaseDel(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.DatabaseDel("channel/sip", "foo")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, "DATABASE DEL channel/sip foo\n", buf.String())
}

func TestCmdDatabaseDelTree(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.DatabaseDelTree("channel", "sip")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, "DATABASE DELTREE channel sip\n", buf.String())
}

func TestCmdDatabaseGet(t *testing.T) {
	agi, buf := mockAGI("200 result=1 (SIP/router01-000e57a5)\n")
	resp, err := agi.DatabaseGet("channel", "sip/111")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, "SIP/router01-000e57a5", resp.Value())
	assert.Equal(t, "DATABASE GET channel sip/111\n", buf.String())
}

func TestCmdDatabasePut(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.DatabasePut("callwait", "sip/1111", "on")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, "DATABASE PUT callwait sip/1111 on\n", buf.String())
}

func TestCmdExec(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.Exec("GoTo", "default,s,1")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, "EXEC GoTo \"default,s,1\"\n", buf.String())
}

func TestCmdGetData(t *testing.T) {
	tests := []struct {
		response string
		file     string
		tout     int
		max      int
		code     int
		result   int
		value    string
		data     string
		exec     string
	}{
		{
			"200 result= (timeout)\n",
			"mainmenu", 0, 1,
			200, 0, "timeout", "",
			"GET DATA mainmenu 0 1\n",
		}, {
			"200 result=*123 (timeout)\n",
			"mainmenu", 1500, 4,
			200, 0, "timeout", "*123",
			"GET DATA mainmenu 1500 4\n",
		}, {
			"200 result=23\n",
			"hollidays", 2000, 2,
			200, 0, "", "23",
			"GET DATA hollidays 2000 2\n",
		}, {
			"200 result=-1\n",
			"hollidays", 3000, 3,
			200, -1, "", "",
			"GET DATA hollidays 3000 3\n",
		}, {
			"200 result=\n", // user press # only
			"hollidays", 3000, 3,
			200, 0, "", "",
			"GET DATA hollidays 3000 3\n",
		}, {
			"511 Command Not Permitted on a dead channel or intercept routine\n",
			"hollidays", 3000, 3,
			511, 0, "Command Not Permitted on a dead channel or intercept routine",
			"Command Not Permitted on a dead channel or intercept routine",
			"GET DATA hollidays 3000 3\n",
		},
	}

	for _, tc := range tests {
		agi, buf := mockAGI(tc.response)
		resp, err := agi.GetData(tc.file, tc.tout, tc.max)
		assert.Nil(t, err)
		assert.Equal(t, tc.code, resp.Code(), "Code:"+tc.response)
		assert.Equal(t, tc.result, resp.Result(), "Result:"+tc.response)
		assert.Equal(t, tc.value, resp.Value(), "Value:"+tc.response)
		assert.Equal(t, tc.data, resp.Data(), "Data:"+tc.response)
		assert.Equal(t, tc.exec, buf.String(), tc.response)
	}

	agi, _ := mockAGI("foo")
	resp, err := agi.GetData("hello", 1000, 3)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invalid input")
	assert.Nil(t, resp)
}

func TestCmdGetFullVariable(t *testing.T) {
	agi, buf := mockAGI("200 result=1 (fr)")
	resp, err := agi.GetFullVariable("CHANNEL(language)", "")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "fr", resp.Value())
	assert.Equal(t, "GET FULL VARIABLE CHANNEL(language)\n", buf.String())

	agi, buf = mockAGI("200 result=1 (Today is sunny)")
	resp, err = agi.GetFullVariable("WEATHER", "SIP/2222-12-00000008")
	assert.Nil(t, err)
	assert.Equal(t, "Today is sunny", resp.Value())
	assert.Equal(t, "GET FULL VARIABLE WEATHER SIP/2222-12-00000008\n", buf.String())
}

func TestCmdGetOption(t *testing.T) {
	agi, buf := mockAGI("200 result=0 endpos=7680")
	resp, err := agi.GetOption("welcome_menu", "", 500)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, 0, resp.Result())
	assert.Equal(t, "GET OPTION welcome_menu \"\" 500\n", buf.String())
}

func TestCmdGetVariable(t *testing.T) {
	agi, buf := mockAGI(`200 result=1 ("Alice Foo" <5145553322>)`)
	resp, err := agi.GetVariable("CALLERID(all)")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code())
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "\"Alice Foo\" <5145553322>", resp.Value())
	assert.Equal(t, "GET VARIABLE CALLERID(all)\n", buf.String())
}

func TestCmdHangup(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.Hangup()
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "HANGUP\n", buf.String())

	agi, buf = mockAGI(respOk)
	resp, err = agi.Hangup("SIP/111-222")
	assert.Nil(t, err)
	assert.Equal(t, "HANGUP SIP/111-222\n", buf.String())

	agi, buf = mockAGI(respOk)
	resp, err = agi.Hangup("IAX/333-444", "SIP/111-222")
	assert.Nil(t, err)
	assert.Equal(t, "HANGUP IAX/333-444\n", buf.String())
}

func TestCmdReceiveChar(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.ReceiveChar(2000)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "RECEIVE CHAR 2000\n", buf.String())
}

func TestCmdReceiveText(t *testing.T) {
	agi, buf := mockAGI("200 result=1 (Hello Bob)")
	resp, err := agi.ReceiveText(2000)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "Hello Bob", resp.Value())
	assert.Equal(t, "RECEIVE TEXT 2000\n", buf.String())
}

func TestCmdRecordFile(t *testing.T) {
	tests := []struct {
		input              string
		file, ftype, digit string
		tout, offset       int
		beep               bool
		silence            int
		result             int
		value, data        string
		endpos             int64
		cmd                string
	}{
		{
			"200 result=1 (timeout) endpos=86435\n",
			"new_rec", "wav", "029", -1, 0, false, 0,
			1, "timeout", "", 86435,
			"RECORD FILE new_rec wav \"029\" -1\n",
		}, {
			"200 result=* (dtmf) endpos=1554\n",
			"new_rec", "wav", "*#", 1000, 600, false, 0,
			1, "dtmf", "*", 1554,
			"RECORD FILE new_rec wav \"*#\" 1000 600\n",
		}, {
			"200 result=-1 (hangup) endpos=0\n",
			"new_rec", "wav", "*#", 1000, 600, true, 500,
			-1, "hangup", "", 0,
			"RECORD FILE new_rec wav \"*#\" 1000 600 BEEP s=500\n",
		}, {
			"200 result=4 (dtmf) endpos=0\n",
			"new_rec", "wav", "*#4", 1000, 600, false, 500,
			1, "dtmf", "4", 0,
			"RECORD FILE new_rec wav \"*#4\" 1000 600 s=500\n",
		},
	}

	for _, tc := range tests {
		agi, buf := mockAGI(tc.input)
		resp, err := agi.RecordFile(tc.file, tc.ftype, tc.digit, tc.tout, tc.offset,
			tc.beep, tc.silence)
		assert.Nil(t, err)
		assert.Equal(t, tc.result, resp.Result())
		assert.Equal(t, tc.value, resp.Value())
		assert.Equal(t, tc.data, resp.Data())
		assert.Equal(t, tc.endpos, resp.EndPos())
		assert.Equal(t, tc.cmd, buf.String())
	}

	agi, buf := mockAGI("")
	resp, err := agi.RecordFile("new_rec", "wav", "", 100, 0, false, 0)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "RECORD FILE new_rec wav \"\" 100\n", buf.String())
}

func TestCmdSayAlpha(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SayAlpha("abc", "*1")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SAY ALPHA abc \"*1\"\n", buf.String())
}

func TestCmdSayDate(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SayDate("1563844045", "")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SAY DATE 1563844045 \"\"\n", buf.String())
}

func TestCmdSayDatetime(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SayDatetime("1563844045", "*#0", "dB", "UTC")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SAY DATETIME 1563844045 \"*#0\" \"dB\" \"UTC\"\n", buf.String())
}

func TestCmdSayDigits(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SayDigits("2234", "")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SAY DIGITS 2234 \"\"\n", buf.String())
}

func TestCmdSayNumber(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SayNumber("1000", "01")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SAY NUMBER 1000 \"01\"\n", buf.String())
}

func TestCmdSayPhonetic(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SayPhonetic("welcome", "")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SAY PHONETIC welcome \"\"\n", buf.String())
}

func TestCmdSayTime(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SayTime("1563844046", "*#0123456789")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SAY TIME 1563844046 \"*#0123456789\"\n", buf.String())
}

func TestCmdSendImage(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SendImage("logo.png")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SEND IMAGE \"logo.png\"\n", buf.String())
}

func TestCmdSendText(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SendText("Hello there!")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SEND TEXT \"Hello there!\"\n", buf.String())
}

func TestCmdSetAutoHangup(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SetAutoHangup(0)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SET AUTOHANGUP 0\n", buf.String())
}

func TestCmdSetCallerid(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SetCallerid("5145553322")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SET CALLERID \"5145553322\"\n", buf.String())
}

func TestCmdSetContext(t *testing.T) {
	agi, buf := mockAGI("200 result=0")
	resp, err := agi.SetContext("default")
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Result())
	assert.Equal(t, "SET CONTEXT default\n", buf.String())
}

func TestCmdSetExtension(t *testing.T) {
	agi, buf := mockAGI("200 result=0")
	resp, err := agi.SetExtension("s")
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Result())
	assert.Equal(t, "SET EXTENSION s\n", buf.String())
}

func TestCmdSetMusic(t *testing.T) {
	tests := []struct {
		enable bool
		class  string
		cmd    string
	}{
		{true, "", "SET MUSIC on \"\"\n"},
		{true, "jazz", "SET MUSIC on \"jazz\"\n"},
		{false, "lounge", "SET MUSIC off \"lounge\"\n"},
	}

	for _, tc := range tests {
		agi, buf := mockAGI(respOk)
		resp, err := agi.SetMusic(tc.enable, tc.class)
		assert.Nil(t, err)
		assert.Equal(t, 1, resp.Result())
		assert.Equal(t, tc.cmd, buf.String())
	}
}

func TestCmdSetPriority(t *testing.T) {
	agi, buf := mockAGI("200 result=0")
	resp, err := agi.SetPriority("1")
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Result())
	assert.Equal(t, "SET PRIORITY 1\n", buf.String())
}

func TestCmdSetVariable(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.SetVariable("CITY", "Toronto")
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "SET VARIABLE CITY \"Toronto\"\n", buf.String())

	agi, buf = mockAGI(respOk)
	resp, err = agi.SetVariable("CALLERID(all)", `"Alice Foo" <5145553322>`)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, `SET VARIABLE CALLERID(all) "\"Alice Foo\" <5145553322>"`+"\n",
		buf.String())
}

func TestCmdStreamFile(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.StreamFile("rec109234", "", 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "STREAM FILE rec109234 \"\" 0\n", buf.String())
}

func TestCmdTDDMode(t *testing.T) {
	tests := []struct {
		mode, cmd string
	}{
		{"on", "TDD MODE on\n"},
		{"off", "TDD MODE off\n"},
		{"mate", "TDD MODE mate\n"},
		{"tdd", "TDD MODE tdd\n"},
		{"foo", "TDD MODE off\n"},
	}

	for _, tc := range tests {
		agi, buf := mockAGI(respOk)
		resp, err := agi.TDDMode(tc.mode)
		assert.Nil(t, err)
		assert.Equal(t, 1, resp.Result())
		assert.Equal(t, tc.cmd, buf.String())
	}
}

func TestCmdVerbose(t *testing.T) {
	agi, buf := mockAGI(respOk)
	agi.Verbose("Debug message", 3)
	assert.Equal(t, "VERBOSE \"Debug message\" 3\n", buf.String())

	agi, buf = mockAGI(respOk)
	agi.Verbose("Debug message")
	assert.Equal(t, "VERBOSE \"Debug message\" 1\n", buf.String())

	agi, buf = mockAGI(respOk)
	agi.Verbose("Debug message", 0)
	assert.Equal(t, "VERBOSE \"Debug message\" 1\n", buf.String())

	agi, buf = mockAGI(respOk)
	agi.Verbose("Debug message", 10)
	assert.Equal(t, "VERBOSE \"Debug message\" 1\n", buf.String())
}

func TestCmdWaitForDigit(t *testing.T) {
	agi, buf := mockAGI(respOk)
	resp, err := agi.WaitForDigit(1000)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp.Result())
	assert.Equal(t, "WAIT FOR DIGIT 1000\n", buf.String())
}
