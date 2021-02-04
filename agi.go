package goagi

import (
	"bufio"
	"strings"
	"time"
)

// ErrAGI goagi error
var ErrAGI = newError("AGI session")

const rwDefaultTimeout = time.Second * 1

// Reader interface for AGI object. Can be net.Conn, os.File or crafted
type Reader interface {
	Read(b []byte) (int, error)
	SetReadDeadline(t time.Time) error
}

// Writer interface for AGI object. Can be net.Conn, os.File or crafted
type Writer interface {
	SetWriteDeadline(t time.Time) error
	Write(b []byte) (int, error)
}

/*
Debugger for AGI instance. Any interface that provides Printf method.
Usage example:
```go
	dbg := logger.New(os.Stdout, "myagi:", log.Lmicroseconds)
	r, w := net.Pipe()
	agi, err := goagi.New(r, w, dbg)
```
It should be used only for debugging as it give lots of output.
*/
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
	rwtout   time.Duration
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
Example for agi:
```go
	import (
		"github.com/staskobzar/goagi"
		"os"
	)

	agi, err := goagi.New(os.Stdin, os.Stdout, nil)
	if err != nil {
		panic(err)
	}
	agi.Verbose("Hello World!")
```

Fast agi example:
```go
	ln, err := net.Listen("tcp", "127.0.0.1:4573")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go func(conn net.Conn) {
			agi, err := goagi.New(conn, conn, nil)
			if err != nil {
				panic(err)
			}
			agi.Verbose("Hello World!")
		}(conn)
	}
```
*/
func New(r Reader, w Writer, dbg Debugger) (*AGI, error) {
	agi := &AGI{
		reader:   r,
		writer:   w,
		debugger: dbg,
		rwtout:   rwDefaultTimeout,
	}
	agi.dbg("[>] New AGI")
	sessData, err := agi.sessionInit()
	if err != nil {
		return nil, ErrAGI.Msg("Failed to read setup: %s", err)
	}
	agi.sessionSetup(sessData)
	return agi, nil
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

func (agi *AGI) sessionInit() ([]string, error) {
	agi.dbg("[>] sessionInit")
	buf := bufio.NewReader(agi.reader)
	data := make([]string, 0)

	for {
		tout := time.Now().Add(agi.rwtout)
		if err := agi.reader.SetReadDeadline(tout); err != nil {
			return nil, err
		}
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
// if timeout > 0 then will read with timeout
func (agi *AGI) read(timeout time.Duration) (resp string, code int, err error) {
	agi.dbg("[>] readResponse")
	buf := bufio.NewReader(agi.reader)
	var builder strings.Builder
	moreInputExpected := false

	for {
		if timeout > 0 {
			tout := time.Now().Add(timeout)
			if fail := agi.reader.SetReadDeadline(tout); fail != nil {
				err = fail
				return
			}
		}

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
	if agi.rwtout > 0 {
		tout := time.Now().Add(agi.rwtout)
		agi.dbg(" [v] set write timeout: %dns", tout)

		if err := agi.writer.SetWriteDeadline(tout); err != nil {
			return err
		}
	}

	agi.dbg(" [v] write command: %q\n", string(command))

	_, err := agi.writer.Write(command)
	if err != nil {
		return err
	}
	return nil
}

// write command, read and parse response
func (agi *AGI) execute(cmd string, timeout bool) (Response, error) {
	agi.dbg("[>] execute cmd: %q", cmd)
	if err := agi.write([]byte(cmd)); err != nil {
		return nil, err
	}

	var tout time.Duration
	if timeout {
		tout = agi.rwtout
	}
	agi.dbg(" [v] read timeout=%d", tout)

	resp, code, err := agi.read(tout)
	if err != nil {
		return nil, err
	}

	return agi.parseResponse(resp, code)
}
