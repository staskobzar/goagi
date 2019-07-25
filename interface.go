package goagi

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"time"
)

const execTimeoutSec = 2

// AGI interface structure
type AGI struct {
	env map[string]string
	arg []string
	io  *bufio.ReadWriter
}

var (
	// EInvalEnv error returned when AGI environment header is not valid
	EInvalEnv = errorNew("Invalid AGI env variable")
)

func newInterface(iodev *bufio.ReadWriter) (*AGI, error) {
	agi := &AGI{make(map[string]string), make([]string, 0), iodev}
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

	ctx, cancel := context.WithTimeout(context.Background(), execTimeoutSec*time.Second)
	defer cancel()

	chStr, chErr := agi.read()
	select {
	case <-ctx.Done():
		return nil, errorNew("Execution read timeout.")
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
