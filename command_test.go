package goagi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// command Command
func TestCmdCommandOk(t *testing.T) {
	resp := "200 result=25 endpos=542268\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	code, result, respStr, err := agi.Command("STREAM FILE welcome 09# \"554\"")
	assert.Nil(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, 25, result)
	assert.Equal(t, resp, respStr)
}

func TestCmdCommandFail(t *testing.T) {
	rw := dummyReadWriteWError()
	agi := &AGI{io: rw}
	code, result, respStr, err := agi.Command("STREAM FILE welcome 09# \"554\"")
	assert.NotNil(t, err)
	assert.Equal(t, -1, code)
	assert.Equal(t, -1, result)
	assert.Equal(t, "", respStr)
}

// command Answer
func TestCmdAnswerOk(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)

	agi := &AGI{io: rw}
	ok, err := agi.Answer()
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestCmdAnswerFail(t *testing.T) {
	resp := "200 result=-1\n"
	rw := dummyReadWrite(resp)

	agi := &AGI{io: rw}
	ok, err := agi.Answer()
	assert.Nil(t, err)
	assert.False(t, ok)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	ok, err = agi.Answer()
	assert.NotNil(t, err)
	assert.False(t, ok)
}

// command AsyncAGIBreak
func TestCmdAsyncAGIBreakOk(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)

	agi := &AGI{io: rw}
	ok, err := agi.AsyncAGIBreak()
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestCmdAsyncAGIBreakFail(t *testing.T) {
	resp := "200 result=-1\n"
	rw := dummyReadWrite(resp)

	agi := &AGI{io: rw}
	ok, err := agi.AsyncAGIBreak()
	assert.Nil(t, err)
	assert.False(t, ok)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	ok, err = agi.AsyncAGIBreak()
	assert.NotNil(t, err)
	assert.False(t, ok)
}

// command ChannelStatus
func TestCmdChannelStatusOk(t *testing.T) {
	resp := "200 result=6\n"
	rw := dummyReadWrite(resp)

	agi := &AGI{io: rw}
	status, err := agi.ChannelStatus("")
	assert.Nil(t, err)
	assert.EqualValues(t, 6, status)
}

func TestCmdChannelStatusFail(t *testing.T) {
	resp := "200 result=-1\n"
	rw := dummyReadWrite(resp)

	agi := &AGI{io: rw}
	status, err := agi.ChannelStatus("SIP/00001-44330")
	assert.NotNil(t, err)
	assert.EqualValues(t, -1, status)
	assert.Contains(t, err.Error(), "No channel name matched")

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	status, err = agi.ChannelStatus("")
	assert.NotNil(t, err)
	assert.Equal(t, -1, status)
}

// command ControlStreamFile
func TestCmdControlStreamFileOk(t *testing.T) {
	resp := "200 result=0 endpos=2541236\n"
	rw := dummyReadWrite(resp)

	agi := &AGI{io: rw}
	endpos, err := agi.ControlStreamFile("welcome", "")
	assert.Nil(t, err)
	assert.EqualValues(t, 2541236, endpos)
}

func TestCmdControlStreamFileFail(t *testing.T) {
	resp := "200 result=-1 endpos=2541236\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	_, err := agi.ControlStreamFile("welcome", "123")
	assert.NotNil(t, err)

	rw = dummyReadWriteRError()
	agi = &AGI{io: rw}
	_, err = agi.ControlStreamFile("welcome", "")
	assert.NotNil(t, err)
}

// command DatabaseDel
func TestCmdDatabaseDelOk(t *testing.T) {
	resp := "200 result=1\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	r, err := agi.DatabaseDel("channel", "foo")
	assert.Nil(t, err)
	assert.True(t, r)
}

func TestCmdDatabaseDelFail(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	r, err := agi.DatabaseDel("channel", "foo")
	assert.Nil(t, err)
	assert.False(t, r)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	r, err = agi.DatabaseDel("hello", "world")
	assert.NotNil(t, err)
	assert.False(t, r)
}

// command DatabaseDel
func TestCmdDatabaseDelTreeOk(t *testing.T) {
	resp := "200 result=1\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	r, err := agi.DatabaseDelTree("channel", "foo")
	assert.Nil(t, err)
	assert.True(t, r)
}

func TestCmdDatabaseDelTreeFail(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	r, err := agi.DatabaseDelTree("channel", "foo")
	assert.Nil(t, err)
	assert.False(t, r)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	r, err = agi.DatabaseDelTree("hello", "world")
	assert.NotNil(t, err)
	assert.False(t, r)
}

// command DatabaseDel
func TestCmdDatabaseGetOk(t *testing.T) {
	resp := "200 result=1 (SIP/router01-000e57a5)\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	val, err := agi.DatabaseGet("channel", "foo")
	assert.Nil(t, err)
	assert.Equal(t, "SIP/router01-000e57a5", val)
}

func TestCmdDatabaseGetFail(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	val, err := agi.DatabaseGet("channel", "foo")
	assert.NotNil(t, err)
	assert.Equal(t, "", val)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	val, err = agi.DatabaseGet("hello", "world")
	assert.NotNil(t, err)
	assert.Equal(t, "", val)
}

// command DatabasePut
func TestCmdDatabasePutOk(t *testing.T) {
	resp := "200 result=1\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	r, err := agi.DatabasePut("channel", "foo", "bar")
	assert.Nil(t, err)
	assert.Equal(t, true, r)
}

func TestCmdDatabasePutFail(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	r, err := agi.DatabasePut("channel", "foo", "bar")
	assert.Nil(t, err)
	assert.Equal(t, false, r)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	r, err = agi.DatabasePut("channel", "foo", "bar")
	assert.NotNil(t, err)
	assert.Equal(t, false, r)
}

// command Exec
func TestCmdExecOk(t *testing.T) {
	resp := "200 result=1\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	r, err := agi.Exec("MusicOnHold", "default,15")
	assert.Nil(t, err)
	assert.Equal(t, 1, r)
}

func TestCmdExecFail(t *testing.T) {
	resp := "200 result=-2\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	r, err := agi.Exec("Dial", "PJSIP/bob,,Q(NO_ANSWER)")
	assert.NotNil(t, err)
	assert.Equal(t, -2, r)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	r, err = agi.Exec("Dial", "PJSIP/bob,,Q(NO_ANSWER)")
	assert.NotNil(t, err)
	assert.Equal(t, -1, r)
}

// command GetData
func TestCmdGetDataOk(t *testing.T) {
	resp := "200 result=42 (timeout)\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	res, tout, err := agi.GetData("prompt", 1000, 3)
	assert.Nil(t, err)
	assert.Equal(t, "*", res)
	assert.True(t, tout)

	resp = "200 result=49\n"
	rw = dummyReadWrite(resp)
	agi = &AGI{io: rw}
	res, tout, err = agi.GetData("prompt", 1000, 3)
	assert.Nil(t, err)
	assert.Equal(t, "1", res)
	assert.False(t, tout)
}

func TestCmdGetDataFail(t *testing.T) {
	resp := "200 result=-1\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	res, tout, err := agi.GetData("prompt", 1000, 3)
	assert.NotNil(t, err)
	assert.Equal(t, "", res)
	assert.False(t, tout)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	res, tout, err = agi.GetData("prompt", 1000, 3)
	assert.NotNil(t, err)
	assert.Equal(t, "", res)
	assert.False(t, tout)
}

// command GetFullVariable
func TestCmdGetFullVariableOk(t *testing.T) {
	resp := "200 result=1 (\"John Dow\" <12345>)\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	val, err := agi.GetFullVariable("CALLERID")
	assert.Nil(t, err)
	assert.Equal(t, "\"John Dow\" <12345>", val)

	resp = "200 result=1 (107.5.2.224)\n"
	rw = dummyReadWrite(resp)
	agi = &AGI{io: rw}
	val, err = agi.GetFullVariable("CHANNEL(rtp,dest)", "SIP/112003430-44432")
	assert.Nil(t, err)
	assert.Equal(t, "107.5.2.224", val)
}

func TestCmdGetFullVariableFail(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	val, err := agi.GetFullVariable("CALLERID(null)")
	assert.NotNil(t, err)
	assert.Equal(t, "", val)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	_, err = agi.GetFullVariable("CALLERID(null)")
	assert.NotNil(t, err)
}

// command GetOption
func TestCmdGetOptionOk(t *testing.T) {
	resp := "200 result=0 endpos=10245\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	dig, offset, err := agi.GetOption("welcome_prompt", "", 0)
	assert.Nil(t, err)
	assert.Equal(t, 0, dig)
	assert.EqualValues(t, 10245, offset)

	resp = "200 result=5 endpos=52417854\n"
	rw = dummyReadWrite(resp)
	agi = &AGI{io: rw}
	dig, offset, err = agi.GetOption("welcome_prompt", "12", 1800)
	assert.Nil(t, err)
	assert.Equal(t, 5, dig)
	assert.EqualValues(t, 52417854, offset)
}

func TestCmdGetOptionFail(t *testing.T) {
	resp := "200 result=-1 endpos=0\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	dig, offset, err := agi.GetOption("welcome_prompt", "12", 0)
	assert.NotNil(t, err)
	assert.Equal(t, -1, dig)
	assert.EqualValues(t, 0, offset)

	resp = "200 result=0 endpos=0\n"
	rw = dummyReadWrite(resp)
	agi = &AGI{io: rw}
	dig, offset, err = agi.GetOption("welcome_prompt", "12", 0)
	assert.NotNil(t, err)
	assert.Equal(t, -1, dig)
	assert.EqualValues(t, 0, offset)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	dig, offset, err = agi.GetOption("welcome_prompt", "12", 0)
	assert.NotNil(t, err)
	assert.Equal(t, -1, dig)
	assert.EqualValues(t, 0, offset)
}

// command GetVariable
func TestCmdGetVariableOk(t *testing.T) {
	resp := "200 result=1 (\"John Dow\" <12345>)\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	val, err := agi.GetVariable("CALLERID")
	assert.Nil(t, err)
	assert.Equal(t, "\"John Dow\" <12345>", val)

	resp = "200 result=1 (107.5.2.224)\n"
	rw = dummyReadWrite(resp)
	agi = &AGI{io: rw}
	val, err = agi.GetVariable("CHANNEL(rtp,dest)")
	assert.Nil(t, err)
	assert.Equal(t, "107.5.2.224", val)
}

func TestCmdGetVariableFail(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	val, err := agi.GetVariable("CALLERID(null)")
	assert.NotNil(t, err)
	assert.Equal(t, "", val)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	_, err = agi.GetVariable("CALLERID(null)")
	assert.NotNil(t, err)
}

// command Hangup
func TestCmdHangupOk(t *testing.T) {
	resp := "200 result=1\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	ok, err := agi.Hangup()
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestCmdHangupFail(t *testing.T) {
	resp := "200 result=-1\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	ok, err := agi.Hangup("SIP/0001-4578")
	assert.NotNil(t, err)
	assert.False(t, ok)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	ok, err = agi.Hangup("SIP/0001-4578")
	assert.NotNil(t, err)
	assert.False(t, ok)
}

// command Noop
func TestCmdNoopOk(t *testing.T) {
	resp := "200 result=0\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	err := agi.Noop()
	assert.Nil(t, err)
}

func TestCmdNoopFail(t *testing.T) {
	rw := dummyReadWriteWError()
	agi := &AGI{io: rw}
	err := agi.Noop()
	assert.NotNil(t, err)
}

// command ReceiveChar
func TestCmdReceiveCharOk(t *testing.T) {
	rw := dummyReadWrite("200 result=5 (timeout)\n")
	agi := &AGI{io: rw}
	chr, err := agi.ReceiveChar(0)
	assert.Nil(t, err)
	assert.Equal(t, 5, chr)

	rw = dummyReadWrite("200 result=9\n")
	agi = &AGI{io: rw}
	chr, err = agi.ReceiveChar(500)
	assert.Nil(t, err)
	assert.Equal(t, 9, chr)
}

func TestCmdReceiveCharFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1 (hangup)\n")
	agi := &AGI{io: rw}
	chr, err := agi.ReceiveChar(1000)
	assert.NotNil(t, err)
	assert.Equal(t, -1, chr)

	rw = dummyReadWrite("200 result=0\n")
	agi = &AGI{io: rw}
	chr, err = agi.ReceiveChar(1000)
	assert.NotNil(t, err)
	assert.Equal(t, -1, chr)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	chr, err = agi.ReceiveChar(1000)
	assert.NotNil(t, err)
	assert.Equal(t, -1, chr)
}

// command ReceiveText
func TestCmdReceiveTextOk(t *testing.T) {
	rw := dummyReadWrite("200 result=1 (White fox is lost in the sea)\n")
	agi := &AGI{io: rw}
	text, err := agi.ReceiveText(0)
	assert.Nil(t, err)
	assert.Equal(t, "White fox is lost in the sea", text)
}

func TestCmdReceiveTextFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	text, err := agi.ReceiveText(1000)
	assert.NotNil(t, err)
	assert.Equal(t, "", text)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	text, err = agi.ReceiveText(1000)
	assert.NotNil(t, err)
	assert.Equal(t, "", text)
}

// command RecordFile
func TestCmdRecordFileOk(t *testing.T) {
	rw := dummyReadWrite("200 result=48 (dtmf) endpos=554879\n")
	agi := &AGI{io: rw}
	err := agi.RecordFile("new_rec", "wav", "029", -1, 0, false, 0)
	assert.Nil(t, err)
}

func TestCmdRecordFileFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1 (hangup) endpos=554879\n")
	agi := &AGI{io: rw}
	err := agi.RecordFile("new_rec", "wav", "09", 1000, 1800, true, 500)
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.RecordFile("new_rec", "wav", "09", 1000, 1800, true, 500)
	assert.NotNil(t, err)
}

// command SayAlpha
func TestCmdSayAlphaOk(t *testing.T) {
	rw := dummyReadWrite("200 result=48\n")
	agi := &AGI{io: rw}
	err := agi.SayAlpha("abc", "10")
	assert.Nil(t, err)
}

func TestCmdSayAlphaFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SayAlpha("abc", "10")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SayAlpha("abc", "10")
	assert.NotNil(t, err)
}

// command SayDate
func TestCmdSayDateOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SayDate("1563844045", "0")
	assert.Nil(t, err)
}

// command SayDatetime
func TestCmdSayDatetimeOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SayDatetime("1563844045", "9", "", "")
	assert.Nil(t, err)
}

// command SayDigits
func TestCmdSayDigitsOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SayDigits("4045", "9")
	assert.Nil(t, err)
}

// command SayNumber
func TestCmdSayNumberOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SayNumber("4045", "9")
	assert.Nil(t, err)
}

// command SayPhonetic
func TestCmdSayPhoneticOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SayPhonetic("abcd", "9")
	assert.Nil(t, err)
}

// command SayTime
func TestCmdSayTimeOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SayTime("1563844046", "9")
	assert.Nil(t, err)
}

// command SendImage
func TestCmdSendImageOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SendImage("image_file")
	assert.Nil(t, err)
}

func TestCmdSendImageFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SendImage("image_file")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SendImage("image_file")
	assert.NotNil(t, err)
}

// command SendText
func TestCmdSendTextOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SendText("Hello there, friend!")
	assert.Nil(t, err)
}

func TestCmdSendTextFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SendText("Hello there, friend!")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SendText("Hello there, friend!")
	assert.NotNil(t, err)
}

// command SetAutoHangup
func TestCmdSetAutoHangupOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SetAutoHangup(0)
	assert.Nil(t, err)
}

func TestCmdSetAutoHangupFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SetAutoHangup(1800)
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SetAutoHangup(800)
	assert.NotNil(t, err)
}

// command SetCallerid
func TestCmdSetCalleridOk(t *testing.T) {
	rw := dummyReadWrite("200 result=1\n")
	agi := &AGI{io: rw}
	err := agi.SetCallerid("\"Achlie\" <5544>")
	assert.Nil(t, err)
}

func TestCmdSetCalleridFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SetCallerid("5552452211")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SetCallerid("5552452211")
	assert.NotNil(t, err)
}

// command SetContext
func TestCmdSetContextOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SetContext("default-inbound")
	assert.Nil(t, err)
}

func TestCmdSetContextFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SetContext("default-inbound")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SetContext("default-inbound")
	assert.NotNil(t, err)
}

// command SetExtension
func TestCmdSetExtensionOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SetExtension("55588")
	assert.Nil(t, err)
}

func TestCmdSetExtensionFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SetExtension("s")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SetExtension("VOICEMAIL")
	assert.NotNil(t, err)
}

// command SetMusic
func TestCmdSetMusicOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SetMusic("on", "hip-pop")
	assert.Nil(t, err)
}

func TestCmdSetMusicFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SetMusic("off")
	assert.NotNil(t, err)

	agi = &AGI{io: rw}
	err = agi.SetMusic("foo", "bar")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SetMusic("on")
	assert.NotNil(t, err)
}

// command SetPriority
func TestCmdSetPriorityOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	err := agi.SetPriority("1")
	assert.Nil(t, err)
}

func TestCmdSetPriorityFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SetPriority("label")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SetPriority("12")
	assert.NotNil(t, err)
}

// command SetVariable
func TestCmdSetVariableOk(t *testing.T) {
	rw := dummyReadWrite("200 result=1\n")
	agi := &AGI{io: rw}
	err := agi.SetVariable("FOO", "1234")
	assert.Nil(t, err)
}

func TestCmdSetVariableFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.SetVariable("MOH", "hello")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.SetVariable("CALLERID(num)", "5558876541")
	assert.NotNil(t, err)
}

// command StreamFile
func TestCmdStreamFileOk(t *testing.T) {
	rw := dummyReadWrite("200 result=48 endpos=123654\n")
	agi := &AGI{io: rw}
	r, err := agi.StreamFile("prompt.en", "0", 0)
	assert.Nil(t, err)
	assert.Equal(t, 48, r)
}

func TestCmdStreamFileFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1 endpos=1235\n")
	agi := &AGI{io: rw}
	r, err := agi.StreamFile("prompt.en", "", 1000)
	assert.NotNil(t, err)
	assert.Equal(t, -1, r)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	r, err = agi.StreamFile("prompt.en", "123", 200)
	assert.NotNil(t, err)
	assert.Equal(t, -1, r)
}

// command TDDMode
func TestCmdTDDModeOk(t *testing.T) {
	rw := dummyReadWrite("200 result=1\n")
	agi := &AGI{io: rw}
	err := agi.TDDMode("on")
	assert.Nil(t, err)
}

func TestCmdTDDModeFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	err := agi.TDDMode("mate")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	err = agi.TDDMode("off")
	assert.NotNil(t, err)
}

// command Verbose
func TestCmdVerboseOk(t *testing.T) {
	rw := dummyReadWrite("200 result=1\n")
	agi := &AGI{io: rw}
	err := agi.Verbose("Hello World!")
	assert.Nil(t, err)

	rw = dummyReadWrite("200 result=1\n")
	agi = &AGI{io: rw}
	err = agi.Verbose("Hello World!", 4)
	assert.Nil(t, err)
}

func TestCmdVerboseFail(t *testing.T) {
	rw := dummyReadWriteWError()
	agi := &AGI{io: rw}
	err := agi.Verbose("gonna fail")
	assert.NotNil(t, err)
}

// command WaitForDigit
func TestCmdWaitForDigitOk(t *testing.T) {
	rw := dummyReadWrite("200 result=0\n")
	agi := &AGI{io: rw}
	dig, err := agi.WaitForDigit(-1)
	assert.Nil(t, err)
	assert.Equal(t, "", dig)

	rw = dummyReadWrite("200 result=48\n")
	agi = &AGI{io: rw}
	dig, err = agi.WaitForDigit(1000)
	assert.Nil(t, err)
	assert.Equal(t, "0", dig)

	rw = dummyReadWrite("200 result=42\n")
	agi = &AGI{io: rw}
	dig, err = agi.WaitForDigit(2000)
	assert.Nil(t, err)
	assert.Equal(t, "*", dig)
}

func TestCmdWaitForDigitFail(t *testing.T) {
	rw := dummyReadWrite("200 result=-1\n")
	agi := &AGI{io: rw}
	dig, err := agi.WaitForDigit(-1)
	assert.NotNil(t, err)
	assert.Equal(t, "", dig)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	dig, err = agi.WaitForDigit(-1)
	assert.NotNil(t, err)
	assert.Equal(t, "", dig)
}
