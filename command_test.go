package goagi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
	_, err := agi.ControlStreamFile("welcome", "")
	assert.NotNil(t, err)

	rw = dummyReadWriteWError()
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
	resp := "200 result=23 (timeout)\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	res, tout, err := agi.GetData("prompt", 1000, 3)
	assert.Nil(t, err)
	assert.Equal(t, 23, res)
	assert.True(t, tout)

	resp = "200 result= (timeout)\n"
	rw = dummyReadWrite(resp)
	agi = &AGI{io: rw}
	res, tout, err = agi.GetData("prompt", 1000, 3)
	assert.Nil(t, err)
	assert.Equal(t, -3, res)
	assert.True(t, tout)

	resp = "200 result=358\n"
	rw = dummyReadWrite(resp)
	agi = &AGI{io: rw}
	res, tout, err = agi.GetData("prompt", 1000, 3)
	assert.Nil(t, err)
	assert.Equal(t, 358, res)
	assert.False(t, tout)
}

func TestCmdGetDataFail(t *testing.T) {
	resp := "200 result=-1\n"
	rw := dummyReadWrite(resp)
	agi := &AGI{io: rw}
	res, tout, err := agi.GetData("prompt", 1000, 3)
	assert.NotNil(t, err)
	assert.Equal(t, -1, res)
	assert.False(t, tout)

	rw = dummyReadWriteWError()
	agi = &AGI{io: rw}
	res, tout, err = agi.GetData("prompt", 1000, 3)
	assert.NotNil(t, err)
	assert.Equal(t, -1, res)
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
