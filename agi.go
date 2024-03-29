package goagi

import (
	"bufio"
	"strings"
)

// ErrAGI goagi error
var ErrAGI = newError("AGI session")

// Reader interface for AGI object. Can be net.Conn, os.File or crafted
type Reader interface {
	Read(b []byte) (int, error)
}

// Writer interface for AGI object. Can be net.Conn, os.File or crafted
type Writer interface {
	Write(b []byte) (int, error)
}

// Debugger for AGI instance. Any interface that provides Printf method.
// It should be used only for debugging as it give lots of output.
type Debugger interface {
	Printf(format string, v ...interface{})
}

// AGI object
type AGI struct {
	env      map[string]string
	arg      []string
	reader   Reader
	writer   Writer
	isHUP    bool
	debugger Debugger
}

const (
	codeUnknown int = 0
	codeEarly       = 100
	codeSucc        = 200
	codeE503        = 503
	codeE510        = 510
	codeE511        = 511
	codeE520        = 520
)

var codeMap = map[string]int{
	"100 ": codeEarly,
	"200 ": codeSucc,
	"503 ": codeE503,
	"510 ": codeE510,
	"511 ": codeE511,
	"520 ": codeE520,
}

/*
New creates and returns AGI object.
Can be used to create agi and fastagi sessions.

Parameters:

- Reader that implements Read method

- Writer that implements Write method

- Debugger that allows to deep library debugging. Nil for production.
*/
func New(r Reader, w Writer, dbg Debugger) (*AGI, error) {
	agi := &AGI{
		reader:   r,
		writer:   w,
		debugger: dbg,
	}
	agi.dbg("[>] New AGI")
	sessData, err := agi.sessionInit()
	if err != nil {
		return nil, ErrAGI.Msg("Failed to read setup: %s", err)
	}
	agi.sessionSetup(sessData)
	return agi, nil
}

func (agi *AGI) Close() {
	agi.env = nil
	agi.arg = nil
}

// Env returns AGI environment variable by key
func (agi *AGI) Env(key string) string {
	agi.dbg("[>] Env for %q", key)
	val, ok := agi.env[key]
	if ok {
		return val
	}
	return ""
}

// EnvArgs returns list of environment arguments
func (agi *AGI) EnvArgs() []string {
	agi.dbg("[>] EnvArgs")
	return agi.arg
}

// IsHungup returns true if AGI channel received HANGUP signal
func (agi *AGI) IsHungup() bool {
	return agi.isHUP
}

func (agi *AGI) sessionInit() ([]string, error) {
	agi.dbg("[>] sessionInit")
	buf := bufio.NewReader(agi.reader)
	data := make([]string, 0)

	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\n" {
			break
		}
		agi.dbg(" [v] read line: %q", line)
		data = append(data, line[:len(line)-1])
	}
	return data, nil
}

func (agi *AGI) dbg(pattern string, vargs ...interface{}) {
	if agi.debugger != nil {
		pattern += "\n"
		agi.debugger.Printf(pattern, vargs...)
	}
}

// low level response from device and response as string, matched response
// code, true if channel reported as hangup and error.
func (agi *AGI) read() (resp string, code int, err error) {
	agi.dbg("[>] readResponse")
	buf := bufio.NewReader(agi.reader)
	var builder strings.Builder
	moreInputExpected := false

	for {
		line, fail := buf.ReadString('\n')
		if fail != nil {
			err = fail
			return
		}

		agi.dbg(" [v] got line: %q", line)

		builder.WriteString(line)
		resp = builder.String()
		if codeMatch, ok := matchCode(line); ok {
			code = codeMatch
			builder.Reset()
			return
		}

		if matchPrefix(line, "520-") {
			moreInputExpected = true
		}

		if matchPrefix(line, "HANGUP") {
			agi.isHUP = true
			builder.Reset()
			continue
		}

		if !moreInputExpected {
			err = ErrAGI.Msg("Invalid input while reading response: %q", resp)
			return
		}
	}
}

func matchPrefix(line, pattern string) bool {
	if len(line) < len(pattern) {
		return false
	}
	return line[:len(pattern)] == pattern
}

func matchCode(data string) (int, bool) {
	if len(data) < 4 {
		return 0, false
	}
	if codeMatch, ok := codeMap[data[:4]]; ok {
		return codeMatch, true
	}
	return 0, false
}

func (agi *AGI) write(command []byte) error {
	agi.dbg("[>] readResponse")

	agi.dbg(" [v] writing command: %q", string(command))

	_, err := agi.writer.Write(command)
	if err != nil {
		return err
	}
	agi.dbg(" [v] write successfully done")
	return nil
}

// write command, read and parse response
func (agi *AGI) execute(cmd string) (Response, error) {
	agi.dbg("[>] execute cmd: %q", cmd)
	if err := agi.write([]byte(cmd)); err != nil {
		return nil, err
	}

	resp, code, err := agi.read()
	if err != nil {
		return nil, err
	}

	return agi.parseResponse(resp, code)
}
