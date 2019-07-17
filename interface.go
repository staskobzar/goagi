package goagi

import (
	"bufio"
	"io"
	"strings"
)

// AGI interface structure
type AGI struct {
	env    map[string]string
	arg    []string
	input  io.Reader
	output io.Writer
}

var (
	// EInvalEnv error returned when AGI environment header is not valid
	EInvalEnv = errorNew("Invalid AGI env variable")
)

func newInterface(in io.Reader, out io.Writer) (*AGI, error) {
	agi := &AGI{make(map[string]string), make([]string, 0), in, out}
	scanner := bufio.NewScanner(in)
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
	chResp := make(chan string)
	chErr := make(chan error)
	go func() {
		defer close(chResp)
		defer close(chErr)
		respStr, err := agi.read()
		if err != nil {
			chErr <- err
		}
		chResp <- respStr
	}()

	agi.output.Write([]byte(cmd + "\n"))

	select {
	case err := <-chErr:
		return nil, err
	case respStr := <-chResp:
		resp, err := parseResponse(respStr)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func (agi *AGI) read() (string, error) {
	reader := bufio.NewReader(agi.input)
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(str, "520-") {
		return str, nil
	}
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(s, "520 End") {
			return str, nil
		}
		str += s
	}
}
