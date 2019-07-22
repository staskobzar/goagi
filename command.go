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
func (agi *AGI) ControlStreamFile(filename, digits string, args ...string) (int32, error) {
	resp, err := agi.execute("CONTROL STREAM FILE", filename, digits, args)
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
	resp, err := agi.execute("GET DATA", file, args)
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
