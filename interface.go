package goagi

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// AGI interface structure
type AGI struct {
	env map[string]string
	arg []string
}

func newInterface(in io.Reader, out io.Writer) (*AGI, error) {
	agi := &AGI{make(map[string]string), make([]string, 0)}
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
		return errors.New("Invalid AGI env variable: " + line)
	}
	idx := strings.Index(line, ": ")
	if idx == -1 {
		return errors.New("Invalid AGI env variable: " + line)
	}
	if strings.HasPrefix(line, "agi_arg_") {
		agi.arg = append(agi.arg, line[idx+2:len(line)])
	} else {
		agi.env[line[4:idx]] = line[idx+2 : len(line)]
	}
	return nil
}
