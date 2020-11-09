package goagi

import (
	"bufio"
	"fmt"
	"strings"
)

// AGI interface structure
type AGI struct {
	env    map[string]string
	arg    []string
	hangup bool
	io     *bufio.ReadWriter
}

var (
	// error returned when AGI environment header is not valid
	EInvalEnv = errorNew("Invalid AGI env variable")
)

func newInterface(iodev *bufio.ReadWriter) (*AGI, error) {
	agi := &AGI{
		env:    make(map[string]string),
		arg:    make([]string, 0),
		hangup: false,
		io:     iodev,
	}
	scanner := bufio.NewScanner(iodev)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		if err := agi.setEnv(line); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return agi, nil
}

// Env returns AGI environment variable by key
func (agi *AGI) Env(key string) string {
	val, ok := agi.env[key]
	if ok {
		return val
	}
	return ""
}

// EnvArgs returns list of environment arguments
func (agi *AGI) EnvArgs() []string {
	return agi.arg
}

func (agi *AGI) setEnv(line string) error {
	if !strings.HasPrefix(line, "agi_") {
		return EInvalEnv.withInfo(line)
	}
	idx := strings.Index(line, ": ")
	if idx == -1 {
		return EInvalEnv.withInfo(line)
	}
	if strings.HasPrefix(line, "agi_arg_") {
		agi.arg = append(agi.arg, line[idx+2:])
	} else {
		agi.env[line[4:idx]] = line[idx+2:]
	}
	return nil
}

func (agi *AGI) execute(cmd string, args ...interface{}) (*agiResp, error) {
	_, err := agi.io.WriteString(compileCmd(cmd, args...))
	if err != nil {
		return nil, err
	}
	agi.io.Flush()

	chStr, chErr := agi.read()
	select {
	case str := <-chStr:
		return parseResponse(str)
	case err := <-chErr:
		return nil, err
	}
}

func (agi *AGI) read() (chan string, chan error) {
	chStr := make(chan string)
	chErr := make(chan error)
	go func() {
		defer close(chStr)
		defer close(chErr)

		str, err := agi.io.ReadString('\n')
		if err != nil {
			chErr <- err
			return
		}
		// TODO: read ignore HANGUP or set agi var
		if strings.HasPrefix(str, "HANGUP") {
			str, err := agi.io.ReadString('\n')
			agi.hangup = true // set AGI Hangup flag
			if err != nil {
				chErr <- err
				return
			}
			chStr <- str
		}
		if !strings.HasPrefix(str, "520-") {
			chStr <- str
			return
		}
		for {
			s, err := agi.io.ReadString('\n')
			if err != nil {
				chErr <- err
				return
			}
			str += s
			if strings.HasPrefix(s, "520 End") {
				chStr <- str
				return
			}
		}
	}()
	return chStr, chErr
}

func compileCmd(cmd string, args ...interface{}) string {
	for _, arg := range args {
		val := fmt.Sprintf("%v", arg)
		if len(val) > 0 {
			cmd = fmt.Sprintf("%s %s", cmd, val)
		} else {
			cmd = fmt.Sprintf("%s \"\"", cmd)
		}
	}

	return fmt.Sprintf("%s\n", cmd)
}
