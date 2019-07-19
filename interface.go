package goagi

import (
	"bufio"
	"strings"
)

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
		agi.arg = append(agi.arg, line[idx+2:len(line)])
	} else {
		agi.env[line[4:idx]] = line[idx+2 : len(line)]
	}
	return nil
}

func (agi *AGI) execute(cmd string) (*agiResp, error) {

	agi.io.Write([]byte(cmd + "\n"))

	respStr, err := agi.read()
	if err != nil {
		return nil, err
	}
	return parseResponse(respStr)
}

func (agi *AGI) read() (string, error) {
	str, err := agi.io.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(str, "520-") {
		return str, nil
	}
	for {
		s, err := agi.io.ReadString('\n')
		if err != nil {
			return "", err
		}
		str += s
		if strings.HasPrefix(s, "520 End") {
			return str, nil
		}
	}
}
