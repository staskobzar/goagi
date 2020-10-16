package goagi

import (
	"fmt"
)

// Command sends command as string to the AGI and returns response valus with
// text response
func (agi *AGI) Command(cmd string) (code int, result int, respStr string, err error) {
	resp, err := agi.execute(cmd)
	if err != nil {
		return -1, -1, "", err
	}
	return resp.code, int(resp.result), resp.raw, nil
}

// Answer executes AGI command "ANSWER"
// Answers channel if not already in answer state.
func (agi *AGI) Answer() (bool, error) {
	resp, err := agi.execute("ANSWER")
	if err != nil {
		return false, err
	}
	return resp.isOk(), nil
}

// AsyncAGIBreak Interrupts Async AGI
//	Interrupts expected flow of Async AGI commands and returns control
// to previous source (typically, the PBX dialplan).
func (agi *AGI) AsyncAGIBreak() (bool, error) {
	resp, err := agi.execute("ASYNCAGI BREAK")
	if err != nil {
		return false, err
	}
	// Asterisk res_agi always returns 200 result=0
	// but for the future try to check response.
	return resp.isOk(), nil
}

// ChannelStatus returns status of the connected channel.
//
// If no channel name is given (empty line) then returns the status of the current channel.
//
//Return values:
//	0 - Channel is down and available.
//	1 - Channel is down, but reserved.
//	2 - Channel is off hook.
//	3 - Digits (or equivalent) have been dialed.
//	4 - Line is ringing.
//	5 - Remote end is ringing.
//	6 - Line is up.
//	7 - Line is busy.
func (agi *AGI) ChannelStatus(channel string) (int, error) {
	resp, err := agi.execute("CHANNEL STATUS " + channel)
	if err != nil {
		return -1, err
	}
	if resp.result == -1 {
		return -1, errorNew("No channel name matched the argument given.")
	}
	return int(resp.result), nil
}

// ControlStreamFile sends audio file on channel and allows the listener
// to control the stream.
//	Send the given file, allowing playback to be controlled by the given digits, if any.
// Use double quotes for the digits if you wish none to be permitted. If offsetms
// is provided then the audio will seek to offsetms before play starts.
//	Returns 0 if playback completes without a digit being pressed, or the ASCII numerical
// value of the digit if one was pressed, or -1 on error or if the channel was
// disconnected.
//	Returns the position where playback was terminated as endpos.
//	Example:
//	agi.ControlStreamFile("prompt_en", "19", "3000", "#", "0", "#", "1600")
//	agi.ControlStreamFile("prompt_en", "")
//	agi.ControlStreamFile("prompt_en", "19", "", "", "", "#", "1600")
func (agi *AGI) ControlStreamFile(filename, digits string, args ...interface{}) (int32, error) {
	cmd := "CONTROL STREAM FILE " + filename
	if len(digits) > 0 {
		cmd += " " + digits
	} else {
		cmd += " \"\""
	}
	resp, err := agi.execute(cmd, args...)
	if err != nil {
		return 0, err
	}
	if resp.result == -1 {
		return resp.result, errorNew("Error or channel disconnected.")
	}

	return resp.endpos, nil
}

// DatabaseDel deletes an entry in the Asterisk database for a given family and key.
//	Returns status and error if fails.
func (agi *AGI) DatabaseDel(family, key string) (bool, error) {
	resp, err := agi.execute("DATABASE DELETE", family, key)
	if err != nil {
		return false, err
	}
	ok := resp.code == 200 && resp.result == 1
	return ok, nil
}

// DatabaseDelTree deletes a family or specific keytree within a family in the Asterisk database.
func (agi *AGI) DatabaseDelTree(family, keytree string) (bool, error) {
	resp, err := agi.execute("DATABASE DELTREE", family, keytree)
	if err != nil {
		return false, err
	}
	ok := resp.code == 200 && resp.result == 1
	return ok, nil
}

// DatabaseGet Retrieves an entry in the Asterisk database for a given family and key.
//	Returns value as string or error if failed or value not set
func (agi *AGI) DatabaseGet(family, key string) (string, error) {
	resp, err := agi.execute("DATABASE GET", family, key)
	if err != nil {
		return "", err
	}
	if resp.result == 0 {
		return "", errorNew("Value not set.")
	}
	return resp.value, nil
}

// DatabasePut adds or updates an entry in the Asterisk database for
// a given family, key, and value.
func (agi *AGI) DatabasePut(family, key, val string) (bool, error) {
	resp, err := agi.execute("DATABASE PUT", family, key, val)
	if err != nil {
		return false, err
	}
	ok := resp.code == 200 && resp.result == 1
	return ok, nil
}

// Exec executes application with given options.
func (agi *AGI) Exec(app, opts string) (int, error) {
	resp, err := agi.execute("EXEC", app, opts)
	if err != nil {
		return -1, err
	}

	if resp.result == -2 {
		return -2, errorNew("Could not find application " + app)
	}
	return int(resp.result), nil
}

// GetData Stream the given file, and receive DTMF data.
func (agi *AGI) GetData(file string, args ...interface{}) (digit string, timeout bool, err error) {
	cmd := "GET DATA " + file
	resp, err := agi.execute(cmd, args...)
	if err != nil {
		return "", false, err
	}
	if resp.result < 0 {
		return "", false, errorNew("Failed get data.")
	}
	timeout = resp.value == "timeout"
	digit = string(resp.result)
	return
}

// GetFullVariable evaluates a channel expression
func (agi *AGI) GetFullVariable(name string, channel ...string) (string, error) {
	cmd := "GET FULL VARIABLE " + name
	var resp *agiResp
	var err error
	if len(channel) > 0 {
		resp, err = agi.execute(cmd, channel[0])
	} else {
		resp, err = agi.execute(cmd)
	}
	if err != nil {
		return "", err
	}
	if resp.result == 0 {
		return "", errorNew("Variable is not set.")
	}

	return resp.value, nil
}

// GetOption Stream file, prompt for DTMF, with timeout.
//	Behaves similar to STREAM FILE but used with a timeout option.
//	Returns digit pressed, offset and error
func (agi *AGI) GetOption(filename, digits string, timeout int32) (int, int32, error) {
	cmd := "GET OPTION " + filename
	resp, err := agi.execute(cmd, digits, timeout)
	if err != nil {
		return -1, 0, err
	}

	if resp.result == -1 {
		return -1, 0, errorNew("Command failure")
	}

	if resp.result == 0 && resp.endpos == 0 {
		return -1, 0, errorNew("Failure on open")
	}

	return int(resp.result), resp.endpos, nil
}

// GetVariable Gets a channel variable.
func (agi *AGI) GetVariable(name string) (string, error) {
	resp, err := agi.execute("GET VARIABLE", name)
	if err != nil {
		return "", err
	}
	if resp.result == 0 {
		return "", errorNew("Variable is not set.")
	}

	return resp.value, nil
}

// Hangup a channel.
//	Hangs up the specified channel. If no channel name is given, hangs up the current channel
func (agi *AGI) Hangup(channel ...string) (bool, error) {
	cmd := "HANGUP"
	if len(channel) > 0 {
		cmd += " " + channel[0]
	}
	resp, err := agi.execute(cmd)
	if err != nil {
		return false, err
	}
	if resp.result == -1 {
		return false, errorNew("Failed hangup")
	}
	return true, nil
}

// Noop Does nothing.
func (agi *AGI) Noop() error {
	_, err := agi.execute("NOOP")
	return err
}

// ReceiveChar Receives one character from channels supporting it.
//	Most channels do not support the reception of text. Returns the decimal value of
// the character if one is received, or 0 if the channel does not support text reception.
//	timeout - The maximum time to wait for input in milliseconds, or 0 for infinite. Most channels
//	Returns -1 only on error/hangup.
func (agi *AGI) ReceiveChar(timeout int32) (int, error) {
	resp, err := agi.execute("RECEIVE CHAR", timeout)
	if err != nil {
		return -1, err
	}
	if resp.result == -1 {
		return -1, errorNew("Channel error or hangup.")
	}
	if resp.result == 0 {
		return -1, errorNew("Channel does not support text reception.")
	}
	return int(resp.result), nil
}

// ReceiveText Receives text from channels supporting it.
//	timeout - The timeout to be the maximum time to wait for input in milliseconds, or 0 for infinite.
func (agi *AGI) ReceiveText(timeout int32) (string, error) {
	resp, err := agi.execute("RECEIVE TEXT", timeout)
	if err != nil {
		return "", err
	}
	if resp.result == -1 {
		return "", errorNew("Failure, hangup or timeout.")
	}
	return resp.value, nil
}

// RecordFile Record to a file until a given dtmf digit in the sequence is received.
// The format will specify what kind of file will be recorded. The timeout is the
// maximum record time in milliseconds, or -1 for no timeout.
//	offset samples is optional, and, if provided, will seek to the offset without
// exceeding the end of the file.
//	beep causes Asterisk to play a beep to the channel that is about to be recorded.
//	silence is the number of seconds of silence allowed before the function returns
// despite the lack of dtmf digits or reaching timeout.
//	silence is the number of seconds of silence that are permitted before the
// recording is terminated, regardless of the escape_digits or timeout arguments
func (agi *AGI) RecordFile(file, format, escDigits string,
	timeout, offset int, beep bool, silence int) error {

	cmd := "RECORD FILE"
	cmd = fmt.Sprintf("%s %s %s %s", cmd, file, format, escDigits)
	if beep {
		cmd = fmt.Sprintf("%s BEEP", cmd)
	}
	if silence > 0 {
		cmd = fmt.Sprintf("%s s=%d", cmd, silence)
	}
	resp, err := agi.execute(cmd)
	if err != nil {
		return err
	}
	if resp.result <= 0 {
		return errorNew("Failed record file")
	}
	return nil
}

func (agi *AGI) say(cmd string, args ...interface{}) error {
	resp, err := agi.execute("SAY "+cmd, args...)
	if err != nil {
		return err
	}
	if resp.result < 0 {
		return errorNew("Failure")
	}
	return nil
}

// SayAlpha says a given character string, returning early if any of the given
// DTMF digits are received on the channel.
func (agi *AGI) SayAlpha(number, escDigits string) error {
	return agi.say("ALPHA", number, escDigits)
}

// SayDate say a given date, returning early if any of the given DTMF digits
// are received on the channel
func (agi *AGI) SayDate(date, escDigits string) error {
	return agi.say("DATE", date, escDigits)
}

// SayDatetime say a given time, returning early if any of the given DTMF
// digits are received on the channel
func (agi *AGI) SayDatetime(time, escDigits, format, timezone string) error {
	return agi.say("DATETIME", time, escDigits, format, timezone)
}

// SayDigits say a given digit string, returning early if any of the given
// DTMF digits are received on the channel
func (agi *AGI) SayDigits(number, escDigits string) error {
	return agi.say("DIGITS", number, escDigits)
}

// SayNumber say a given digit string, returning early if any of the given
// DTMF digits are received on the channel
func (agi *AGI) SayNumber(number, escDigits string) error {
	return agi.say("NUMBER", number, escDigits)
}

// SayPhonetic say a given character string with phonetics, returning early
// if any of the given DTMF digits are received on the channel
func (agi *AGI) SayPhonetic(str, escDigits string) error {
	return agi.say("PHONETIC", str, escDigits)
}

// SayTime say a given time, returning early if any of the given DTMF digits
// are received on the channel
func (agi *AGI) SayTime(time, escDigits string) error {
	return agi.say("TIME", time, escDigits)
}

// SendImage Sends the given image on a channel. Most channels do not support
// the transmission of images.
func (agi *AGI) SendImage(image string) error {
	resp, err := agi.execute("SEND IMAGE", image)
	if err != nil {
		return err
	}
	if resp.result < 0 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// SendText Sends the given text on a channel. Most channels do not support
// the transmission of text.
func (agi *AGI) SendText(text string) error {
	text = fmt.Sprintf("\"%s\"", text)
	resp, err := agi.execute("SEND TEXT", text)
	if err != nil {
		return err
	}
	if resp.result < 0 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// SetAutoHangup Cause the channel to automatically hangup at time seconds in the future.
// Setting to 0 will cause the autohangup feature to be disabled on this channel.
func (agi *AGI) SetAutoHangup(seconds int) error {
	resp, err := agi.execute("SET AUTOHANGUP", seconds)
	if err != nil {
		return err
	}
	if resp.result != 0 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// SetCallerid Changes the callerid of the current channel.
func (agi *AGI) SetCallerid(clid string) error {
	resp, err := agi.execute("SET CALLERID", clid)
	if err != nil {
		return err
	}
	if resp.result != 1 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// SetContext Sets the context for continuation upon exiting the application.
func (agi *AGI) SetContext(ctx string) error {
	resp, err := agi.execute("SET CONTEXT", ctx)
	if err != nil {
		return err
	}
	if resp.result != 0 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// SetExtension Changes the extension for continuation upon exiting the application.
func (agi *AGI) SetExtension(ext string) error {
	resp, err := agi.execute("SET EXTENSION", ext)
	if err != nil {
		return err
	}
	if resp.result != 0 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// SetMusic Enables/Disables the music on hold generator. If class is not specified,
// then the default music on hold class will be used.
//	Parameters: opt is "on" or "off", and music class as string
func (agi *AGI) SetMusic(opt string, class ...string) error {
	if opt != "on" && opt != "off" {
		return errorNew("Invalid opt: '" + opt + "'. Must be 'on' or 'off'.")
	}

	if class != nil {
		opt = fmt.Sprintf("%s %s", opt, class[0])
	}

	resp, err := agi.execute("SET MUSIC", opt)
	if err != nil {
		return err
	}
	if resp.result != 0 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// SetPriority Changes the priority for continuation upon exiting the application.
// The priority must be a valid priority or label.
func (agi *AGI) SetPriority(priority string) error {
	resp, err := agi.execute("SET PRIORITY", priority)
	if err != nil {
		return err
	}
	if resp.result != 0 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// SetVariable Sets a variable to the current channel.
func (agi *AGI) SetVariable(name, value string) error {
	value = fmt.Sprintf("\"%s\"", value)
	resp, err := agi.execute("SET VARIABLE", name, value)
	if err != nil {
		return err
	}
	if resp.result != 1 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// StreamFile Send the given file, allowing playback to be interrupted by the given
// digits, if any.
func (agi *AGI) StreamFile(file, escDigits string, offset int) (int, error) {
	resp, err := agi.execute("STREAM FILE", file, escDigits, offset)
	if err != nil {
		return -1, err
	}
	if resp.result == -1 {
		return -1, errorNew("Failure or hangup.")
	}
	return int(resp.result), nil
}

// TDDMode Enable/Disable TDD transmission/reception on a channel.
//	Modes: on, off, mate, tdd
func (agi *AGI) TDDMode(mode string) error {
	resp, err := agi.execute("TDD MODE", mode)
	if err != nil {
		return err
	}
	if resp.result != 1 {
		return errorNew("Failure or hangup.")
	}
	return nil
}

// Verbose Sends message to the console via verbose message system.
// level is the verbose level (1-4)
func (agi *AGI) Verbose(msg string, level ...int) error {
	var err error
	msg = fmt.Sprintf("\"%s\"", msg)
	if level == nil {
		_, err = agi.execute("VERBOSE", msg)
		return err
	}

	lvl := level[0]
	if lvl < 1 && lvl > 4 {
		lvl = 1
	}
	_, err = agi.execute("VERBOSE", msg, lvl)
	return err
}

// WaitForDigit Waits up to timeout *milliseconds* for channel to receive a DTMF digit.
// Use -1 for the timeout value if you desire the call to block indefinitely.
//	Return digit pressed as string or error
func (agi *AGI) WaitForDigit(timeout int) (string, error) {
	resp, err := agi.execute("WAIT FOR DIGIT", timeout)
	if err != nil {
		return "", err
	}
	if resp.result == -1 {
		return "", errorNew("Failed run command")
	}
	if resp.result == 0 {
		return "", nil
	}
	return string(resp.result), nil
}
