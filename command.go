package goagi

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
	l := &lexer{input: resp.data}
	if l.lookForward("endpos=") {
		l.start = 7
		l.pos = len(l.input)
		return l.atoi(), nil
	}
	return -1, errorNew("Invalid response.")
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
	val := resp.data
	ln := len(val)
	if ln < 2 {
		return "", nil
	}
	return val[1 : ln-1], nil
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
func (agi *AGI) GetData(file string, args ...interface{}) (digits int, timeout bool, err error) {
	cmd := "GET DATA " + file
	resp, err := agi.execute(cmd, args...)
	if err != nil {
		return -1, false, err
	}
	if resp.result == -1 {
		return -1, false, errorNew("Failed get data.")
	}
	timeout = len(resp.data) > 8 && resp.data[0:9] == "(timeout)"
	digits = int(resp.result)
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

	l := &lexer{input: resp.data}
	return l.extractResposeValue(), nil
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
	l := &lexer{input: resp.data}
	endpos := l.extractEndpos()

	if resp.result == -1 {
		return -1, 0, errorNew("Command failure")
	}

	if resp.result == 0 && endpos == 0 {
		return -1, 0, errorNew("Failure on open")
	}

	return int(resp.result), endpos, nil
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

	l := &lexer{input: resp.data}
	return l.extractResposeValue(), nil
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
