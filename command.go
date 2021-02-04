package goagi

import (
	"fmt"
)

// Command sends command as string to the AGI and returns response valus with
// text response
func (agi *AGI) Command(cmd string) (Response, error) {
	return agi.execute(cmd + "\n")
}

// Answer executes AGI command "ANSWER"
// Answers channel if not already in answer state.
func (agi *AGI) Answer() (Response, error) {
	return agi.execute("ANSWER\n")
}

// AsyncAGIBreak Interrupts Async AGI
//	Interrupts expected flow of Async AGI commands and returns control
// to previous source (typically, the PBX dialplan).
func (agi *AGI) AsyncAGIBreak() (Response, error) {
	return agi.execute("ASYNCAGI BREAK\n")
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
func (agi *AGI) ChannelStatus(channel string) (Response, error) {
	cmd := fmt.Sprintf("CHANNEL STATUS %s\n", channel)
	return agi.execute(cmd)
}

// ControlStreamFile sends audio file on channel and allows the listener
// to control the stream.
//	Send the given file, allowing playback to be controlled by the given digits, if any.
// Use double quotes for the digits if you wish none to be permitted. If offsetms
// is provided then the audio will seek to offsetms before play starts.
//	Example:
//	agi.ControlStreamFile("prompt_en", "19", "3000", "#", "0", "#", "1600")
//	agi.ControlStreamFile("prompt_en", "")
//	agi.ControlStreamFile("prompt_en", "19", "", "", "", "#", "1600")
//CONTROL STREAM FILE FILENAME ESCAPE_DIGITS
//SKIPMS FFCHAR REWCHR PAUSECHR OFFSETMS
func (agi *AGI) ControlStreamFile(filename, digits string, args ...string) (Response, error) {
	cmd := fmt.Sprintf("CONTROL STREAM FILE %s %q", filename, digits)

	if len(args) > 5 {
		return nil, ErrAGI.Msg("Too many arguments. Unknown args: %v", args[5:])
	}

	for _, v := range args {
		cmd = fmt.Sprintf("%s %q", cmd, v)
	}

	cmd = fmt.Sprintf("%s\n", cmd)
	return agi.execute(cmd)
}

// DatabaseDel deletes an entry in the Asterisk database for a given family and key.
//	Returns status and error if fails.
func (agi *AGI) DatabaseDel(family, key string) (Response, error) {
	cmd := fmt.Sprintf("DATABASE DEL %s %s\n", family, key)
	return agi.execute(cmd)
}

// DatabaseDelTree deletes a family or specific keytree within a family in the Asterisk database.
func (agi *AGI) DatabaseDelTree(family, keytree string) (Response, error) {
	cmd := fmt.Sprintf("DATABASE DELTREE %s %s\n", family, keytree)
	return agi.execute(cmd)
}

// DatabaseGet Retrieves an entry in the Asterisk database for a given family and key.
//	Returns value as string or error if failed or value not set
//  Response.Value() for result
func (agi *AGI) DatabaseGet(family, key string) (Response, error) {
	cmd := fmt.Sprintf("DATABASE GET %s %s\n", family, key)
	return agi.execute(cmd)
}

// DatabasePut adds or updates an entry in the Asterisk database for
// a given family, key, and value.
func (agi *AGI) DatabasePut(family, key, val string) (Response, error) {
	cmd := fmt.Sprintf("DATABASE PUT %s %s %s\n", family, key, val)
	return agi.execute(cmd)
}

// Exec executes application with given options.
func (agi *AGI) Exec(app, opts string) (Response, error) {
	cmd := fmt.Sprintf("EXEC %s %q\n", app, opts)
	return agi.execute(cmd)
}

// GetData Stream the given file, and receive DTMF data.
// Note: when timeout is 0 then Asterisk will use 6 secods.
// Note: Asterisk has strange way to handle get data response.
// Contrary to other responses, where result has numeric value,
// here asterisk puts DTMF to sent by user to result and this value
// may contain "#" and "*".
// To get DTMF sent by user use Response.Data()
// Response.Value() will contain "timeout" if user has not terminated
// input with "#"
func (agi *AGI) GetData(file string, timeout, maxdigit int) (Response, error) {
	cmd := fmt.Sprintf("GET DATA %s %d %d\n", file, timeout, maxdigit)
	resp, err := agi.execute(cmd)
	if err != nil {
		return nil, err
	}
	// special get data result treatment
	if resp.Result() == -1 || resp.Code() != codeSucc {
		return resp, nil
	}

	r := resp.(*responseSuccess)
	r.result = 0
	r.data = scanResultStrFromRaw(r.raw)
	return r, nil
}

// GetFullVariable evaluates a channel expression
func (agi *AGI) GetFullVariable(name, channel string) (Response, error) {
	cmd := "GET FULL VARIABLE"
	if channel == "" {
		cmd = fmt.Sprintf("%s %s\n", cmd, name)
	} else {
		cmd = fmt.Sprintf("%s %s %s\n", cmd, name, channel)
	}
	return agi.execute(cmd)
}

// GetOption Stream file, prompt for DTMF, with timeout.
//	Behaves similar to STREAM FILE but used with a timeout option.
//	Returns digit pressed, offset and error
func (agi *AGI) GetOption(filename, digits string, timeout int32) (Response, error) {
	cmd := fmt.Sprintf("GET OPTION %s %q %d\n", filename, digits, timeout)
	return agi.execute(cmd)
}

// GetVariable Gets a channel variable.
func (agi *AGI) GetVariable(name string) (Response, error) {
	cmd := fmt.Sprintf("GET VARIABLE %s\n", name)
	return agi.execute(cmd)
}

// Hangup a channel.
//	Hangs up the specified channel. If no channel name is given, hangs up the current channel
func (agi *AGI) Hangup(channel ...string) (Response, error) {
	cmd := "HANGUP"

	if len(channel) > 0 {
		cmd = fmt.Sprintf("%s %s\n", cmd, channel[0])
	} else {
		cmd = fmt.Sprintf("%s\n", cmd)
	}
	return agi.execute(cmd)
}

// ReceiveChar Receives one character from channels supporting it.
//	Most channels do not support the reception of text. Returns the decimal value of
// the character if one is received, or 0 if the channel does not support text reception.
//	timeout - The maximum time to wait for input in milliseconds, or 0 for infinite.
// Returns result -1 on error or char byte
func (agi *AGI) ReceiveChar(timeout int) (Response, error) {
	cmd := fmt.Sprintf("RECEIVE CHAR %d\n", timeout)
	return agi.execute(cmd)
}

// ReceiveText Receives text from channels supporting it.
//	timeout - The timeout to be the maximum time to wait for input in milliseconds, or 0 for infinite.
func (agi *AGI) ReceiveText(timeout int) (Response, error) {
	cmd := fmt.Sprintf("RECEIVE TEXT %d\n", timeout)
	return agi.execute(cmd)
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
// If interupted by DTMF, digits will be available in Response.Data()
func (agi *AGI) RecordFile(file, format, escDigits string,
	timeout, offset int, beep bool, silence int) (Response, error) {

	cmd := "RECORD FILE"
	cmd = fmt.Sprintf("%s %s %s %q %d", cmd, file, format, escDigits, timeout)
	if offset > 0 {
		cmd = fmt.Sprintf("%s %d", cmd, offset)
	}
	if beep {
		cmd = fmt.Sprintf("%s BEEP", cmd)
	}
	if silence > 0 {
		cmd = fmt.Sprintf("%s s=%d", cmd, silence)
	}
	cmd = fmt.Sprintf("%s\n", cmd)

	resp, err := agi.execute(cmd)
	if err != nil {
		return nil, err
	}

	if resp.Value() != "dtmf" {
		return resp, nil
	}

	r := resp.(*responseSuccess)
	r.result = 1
	r.data = scanResultStrFromRaw(r.raw)

	return r, nil
}

// SayAlpha says a given character string, returning early if any of the given
// DTMF digits are received on the channel.
func (agi *AGI) SayAlpha(number, escDigits string) (Response, error) {
	cmd := fmt.Sprintf("SAY ALPHA %s %q\n", number, escDigits)
	return agi.execute(cmd)
}

// SayDate say a given date, returning early if any of the given DTMF digits
// are received on the channel
func (agi *AGI) SayDate(date, escDigits string) (Response, error) {
	cmd := fmt.Sprintf("SAY DATE %s %q\n", date, escDigits)
	return agi.execute(cmd)
}

// SayDatetime say a given time, returning early if any of the given DTMF
// digits are received on the channel
func (agi *AGI) SayDatetime(time, escDigits, format, timezone string) (Response, error) {
	cmd := fmt.Sprintf("SAY DATETIME %s %q %q %q\n", time, escDigits, format, timezone)
	return agi.execute(cmd)
}

// SayDigits say a given digit string, returning early if any of the given
// DTMF digits are received on the channel
func (agi *AGI) SayDigits(number, escDigits string) (Response, error) {
	cmd := fmt.Sprintf("SAY DIGITS %s %q\n", number, escDigits)
	return agi.execute(cmd)
}

// SayNumber say a given digit string, returning early if any of the given
// DTMF digits are received on the channel
func (agi *AGI) SayNumber(number, escDigits string) (Response, error) {
	cmd := fmt.Sprintf("SAY NUMBER %s %q\n", number, escDigits)
	return agi.execute(cmd)
}

// SayPhonetic say a given character string with phonetics, returning early
// if any of the given DTMF digits are received on the channel
func (agi *AGI) SayPhonetic(str, escDigits string) (Response, error) {
	cmd := fmt.Sprintf("SAY PHONETIC %s %q\n", str, escDigits)
	return agi.execute(cmd)
}

// SayTime say a given time, returning early if any of the given DTMF digits
// are received on the channel
func (agi *AGI) SayTime(time, escDigits string) (Response, error) {
	cmd := fmt.Sprintf("SAY TIME %s %q\n", time, escDigits)
	return agi.execute(cmd)
}

// SendImage Sends the given image on a channel. Most channels do not support
// the transmission of images.
func (agi *AGI) SendImage(image string) (Response, error) {
	cmd := fmt.Sprintf("SEND IMAGE %q\n", image)
	return agi.execute(cmd)
}

// SendText Sends the given text on a channel. Most channels do not support
// the transmission of text.
func (agi *AGI) SendText(text string) (Response, error) {
	cmd := fmt.Sprintf("SEND TEXT %q\n", text)
	return agi.execute(cmd)
}

// SetAutoHangup Cause the channel to automatically hangup at time seconds in the future.
// Setting to 0 will cause the autohangup feature to be disabled on this channel.
func (agi *AGI) SetAutoHangup(seconds int) (Response, error) {
	cmd := fmt.Sprintf("SET AUTOHANGUP %d\n", seconds)
	return agi.execute(cmd)
}

// SetCallerid Changes the callerid of the current channel.
func (agi *AGI) SetCallerid(clid string) (Response, error) {
	cmd := fmt.Sprintf("SET CALLERID %q\n", clid)
	return agi.execute(cmd)
}

// SetContext Sets the context for continuation upon exiting the application.
func (agi *AGI) SetContext(ctx string) (Response, error) {
	cmd := fmt.Sprintf("SET CONTEXT %s\n", ctx)
	return agi.execute(cmd)
}

// SetExtension Changes the extension for continuation upon exiting the application.
func (agi *AGI) SetExtension(ext string) (Response, error) {
	cmd := fmt.Sprintf("SET EXTENSION %s\n", ext)
	return agi.execute(cmd)
}

// SetMusic Enables/Disables the music on hold generator. If class is not specified,
// then the default music on hold class will be used.
func (agi *AGI) SetMusic(enable bool, class string) (Response, error) {
	cmd := "SET MUSIC"

	if enable {
		cmd = fmt.Sprintf("%s on", cmd)
	} else {
		cmd = fmt.Sprintf("%s off", cmd)
	}

	cmd = fmt.Sprintf("%s %q\n", cmd, class)
	return agi.execute(cmd)
}

// SetPriority Changes the priority for continuation upon exiting the application.
// The priority must be a valid priority or label.
func (agi *AGI) SetPriority(priority string) (Response, error) {
	cmd := fmt.Sprintf("SET PRIORITY %s\n", priority)
	return agi.execute(cmd)
}

// SetVariable Sets a variable to the current channel.
func (agi *AGI) SetVariable(name, value string) (Response, error) {
	cmd := fmt.Sprintf("SET VARIABLE %s %q\n", name, value)
	return agi.execute(cmd)
}

// StreamFile Send the given file, allowing playback to be interrupted by the given
// digits, if any.
func (agi *AGI) StreamFile(file, escDigits string, offset int) (Response, error) {
	cmd := fmt.Sprintf("STREAM FILE %s %q %d\n", file, escDigits, offset)
	return agi.execute(cmd)
}

// TDDMode Enable/Disable TDD transmission/reception on a channel.
//	Modes: on, off, mate, tdd
func (agi *AGI) TDDMode(mode string) (Response, error) {
	cmd := "TDD MODE"
	switch mode {
	case "on", "off", "mate", "tdd":
		cmd = fmt.Sprintf("%s %s\n", cmd, mode)
	default:
		cmd = fmt.Sprintf("%s off\n", cmd)
	}
	return agi.execute(cmd)
}

// Verbose Sends message to the console via verbose message system.
// level is the verbose level (1-4)
func (agi *AGI) Verbose(msg string, level ...int) (Response, error) {
	cmd := fmt.Sprintf("VERBOSE %q", msg)
	lvl := 1
	if level != nil {
		if level[0] > 0 && level[0] < 5 {
			lvl = level[0]
		}
	}
	cmd = fmt.Sprintf("%s %d\n", cmd, lvl)
	return agi.execute(cmd)
}

// WaitForDigit Waits up to timeout *milliseconds* for channel to receive a DTMF digit.
// Use -1 for the timeout value if you desire the call to block indefinitely.
//	Return digit pressed as string or error
func (agi *AGI) WaitForDigit(timeout int) (Response, error) {
	cmd := fmt.Sprintf("WAIT FOR DIGIT %d\n", timeout)
	return agi.execute(cmd)
}
