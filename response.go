package goagi

import (
	"strconv"
	"strings"
)

type Response interface {
	Code() int
	RawResponse() string
	Result() int
	Value() string
	Data() string
	EndPos() int64
	Digit() string
	SResults() int
}

type response struct {
	code   int
	result int
	raw    string
	data   string
}

func (r *response) Code() int           { return r.code }
func (r *response) RawResponse() string { return r.raw }
func (r *response) Result() int         { return r.result }
func (r *response) Value() string       { return r.data }
func (r *response) Data() string        { return r.data }
func (r *response) EndPos() int64       { return 0 }
func (r *response) Digit() string       { return "" }
func (r *response) SResults() int       { return 0 }

type responseSuccess struct {
	response
	value    string
	endpos   int64
	digit    string
	sresults int
}

func (r *responseSuccess) Value() string { return r.value }
func (r *responseSuccess) EndPos() int64 { return r.endpos }
func (r *responseSuccess) Digit() string { return r.digit }
func (r *responseSuccess) SResults() int { return r.sresults }

func (agi *AGI) sessionSetup(data []string) {
	agi.dbg("[>] sessionSetup")
	agi.env = make(map[string]string)
	agi.arg = make([]string, 0)

	for _, line := range data {
		idx := strings.Index(line, ": ")
		if idx == -1 || line[:4] != "agi_" {
			agi.dbg(" [!] ignore invalid line: %q", line)
			continue
		}
		if line[:8] == "agi_arg_" {
			arg := line[idx+2:]
			agi.arg = append(agi.arg, arg)
			agi.dbg(" [v] add arg: %q", arg)
			continue
		}
		key, val := line[4:idx], line[idx+2:]
		agi.dbg(" [v] add env: %s => %s", key, val)
		agi.env[key] = val
	}
}

// read and parse response
func (agi *AGI) parseResponse(data string, code int) (Response, error) {
	agi.dbg("[>] parseResponse")

	if len(data) < 4 {
		agi.dbg(" [!] response is invalid: %q", data)
		return nil, ErrAGI.Msg("Recieved response is invalid: %q", data)
	}

	if code == codeEarly {
		return agi.parseEarlyResponse(data, code)
	}

	if code == codeSucc {
		return agi.parseSuccessResponse(data)
	}

	if code > 500 && code < 600 {
		return agi.parseErrorResponse(data, code)
	}
	return nil, ErrAGI.Msg("Can not recognize the response: %q", data)
}

func (agi *AGI) parseSuccessResponse(data string) (Response, error) {
	agi.dbg("[>] parseSuccessResponse: %q", data)
	resp := &responseSuccess{}
	resp.code = codeSucc
	resp.raw = data

	data = data[4:]
	data = trimLastNL(data)
	data, result := scanResult(data)
	resp.result = result

	if len(data) == 0 {
		return resp, nil
	}

	// parse value
	data, value := parseValue(data)
	resp.value = value
	tokens := strings.Fields(data)
	resp.endpos = parseEndpos(tokens)
	resp.digit = parseDigit(tokens)
	resp.sresults = parseSResults(tokens)
	return resp, nil
}

func parseSResults(tokens []string) int {
	for _, tok := range tokens {
		if len(tok) < 9 || tok[:8] != "results=" {
			continue
		}
		if num, err := strconv.Atoi(tok[8:]); err == nil {
			return num
		}
		return 0
	}
	return 0
}

func parseDigit(tokens []string) string {
	for _, tok := range tokens {
		if len(tok) < 7 || tok[:6] != "digit=" {
			continue
		}
		return tok[6:]
	}
	return ""
}

func parseEndpos(tokens []string) int64 {
	for _, tok := range tokens {
		if len(tok) < 8 || tok[:7] != "endpos=" {
			continue
		}
		if pos, err := strconv.ParseInt(tok[7:], 10, 64); err == nil {
			return pos
		}
		return 0
	}
	return 0
}

func parseValue(data string) (string, string) {
	if data[0] != '(' {
		return data, ""
	}
	if idx := strings.IndexByte(data, ')'); idx > 0 {
		return data[idx+1:], data[1:idx]
	}
	return data, ""
}

func (agi *AGI) parseEarlyResponse(data string, code int) (Response, error) {
	agi.dbg("[>] parseEarlyResponse %d: %q", code, data)
	resp := &response{}
	resp.code = code
	resp.raw = data

	data = data[4:]
	data = trimLastNL(data)

	data, result := scanResult(data)
	resp.result = result
	resp.data = data

	return resp, nil
}

func (agi *AGI) parseErrorResponse(data string, code int) (Response, error) {
	agi.dbg("[>] parseErrorResponse %d: %q", code, data)
	resp := &response{}
	resp.code = code
	resp.raw = data

	if code == codeE520 && data[:4] == "520-" {
		resp.data = scanE520Usage(data)
		return resp, nil
	}

	data = data[4:]
	data = trimLastNL(data)

	data, result := scanResult(data)
	resp.result = result
	resp.data = data

	return resp, nil
}

func trimLastNL(data string) string {
	l := len(data)
	if l > 0 && data[l-1:] == "\n" {
		return data[:l-1]
	}
	return data
}

func scanResult(data string) (string, int) {
	if len(data) < 7 || data[:7] != "result=" {
		return data, 0
	}
	token := strings.SplitN(data, " ", 2)
	resultTok := strings.SplitN(token[0], "=", 2)

	if len(resultTok) > 1 {
		if num, err := strconv.Atoi(resultTok[1]); err == nil {
			return strings.Join(token[1:], " "), num
		}
	}
	return strings.Join(token[1:], " "), 0
}

func scanResultStrFromRaw(data string) string {
	raw := data[4:]
	raw = trimLastNL(raw)
	tokens := strings.Fields(raw)
	if strings.Compare(tokens[0], "result=") <= 0 {
		return ""
	}
	return tokens[0][7:]
}

func scanE520Usage(data string) string {
	token := strings.Split(data, "\n")

	return strings.Join(token[1:len(token)-2], "\n")
}
