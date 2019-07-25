package goagi

import (
	"unicode"
	"unicode/utf8"
)

var (
	// EInvalResp error returns when AGI response does not match pattern
	EInvalResp = errorNew("Invalid AGI response")
	// EHangUp error when HANGUP signal received
	EHangUp = errorNew("HANGUP")
)

type agiResp struct {
	code   int
	result int32
	endpos int32
	value  string
	data   string
	raw    string
}

func (resp *agiResp) isOk() bool {
	return resp.code == 200 && resp.result == 0
}

func parseResponse(text string) (*agiResp, error) {
	resp := &agiResp{result: -1, endpos: -1, raw: text}
	lex := &lexer{input: text}

	if lex.lookForward("HANGUP\n") {
		return resp, EHangUp
	}

	// state machine
	for scanner, err := scanCode(lex, resp); ; scanner, err = scanner(lex, resp) {
		if err != nil {
			return resp, err
		}
		if scanner == nil {
			break
		}
	}
	return resp, nil
}

// scanner routins
type scanFunc func(*lexer, *agiResp) (scanFunc, error)

func scanCode(l *lexer, resp *agiResp) (scanFunc, error) {
	char := l.peek()
	if !unicode.IsDigit(char) {
		return nil, EInvalResp.withInfo("scanCode:digit expected:" + l.input)
	}
	for {
		char = l.next()
		if !unicode.IsDigit(char) {
			l.backup()
			break
		}
	}
	resp.code = int(l.atoi())
	l.ignore()

	if resp.code >= 500 {
		return scanError, nil
	}
	return scanResult, nil
}

func scanResult(l *lexer, resp *agiResp) (scanFunc, error) {
	pattern := "result="
	char := l.next()
	if char != ' ' {
		return nil, EInvalResp.withInfo("scanResult:space expected:" + l.input)
	}
	l.ignore()

	if !l.lookForward(pattern) {
		return nil, EInvalResp.withInfo("scanResult:result= expected:" + l.input)
	}
	l.pos += len(pattern)
	l.ignore()

	for {
		char = l.next()
		if unicode.IsSpace(char) || char == eof {
			l.backup()
			break
		}
	}
	if l.start == l.pos {
		resp.result = -3
	} else {
		resp.result = l.atoi()
	}
	return scanValue, nil
}

func scanValue(l *lexer, resp *agiResp) (scanFunc, error) {
	if pos := l.skipWhitespace(); pos == eof {
		return nil, nil
	}
	if chr := l.peek(); chr != '(' {
		return scanEndpos, nil
	}
	l.next() // skip "("
	l.ignore()
	for {
		char := l.next()
		if char == eof || char == '\n' || char == ')' {
			l.backup()
			break
		}
	}
	resp.value = l.input[l.start:l.pos]
	l.next()
	l.ignore()
	return scanEndpos, nil
}

func scanEndpos(l *lexer, resp *agiResp) (scanFunc, error) {
	if pos := l.skipWhitespace(); pos == eof {
		return nil, nil
	}
	pattern := "endpos="
	if !l.lookForward(pattern) {
		return scanData, nil
	}
	l.pos += len(pattern)
	l.start = l.pos
	for {
		if chr := l.peek(); unicode.IsSpace(chr) || chr == eof {
			break
		}
		l.next()
	}
	resp.endpos = l.atoi()
	return scanData, nil
}

func scanData(l *lexer, resp *agiResp) (scanFunc, error) {
	if pos := l.skipWhitespace(); pos == eof {
		return nil, nil
	}
	for {
		char := l.next()
		if char == eof || char == '\n' {
			l.backup()
			break
		}
	}
	resp.data = l.input[l.start:l.pos]
	return nil, nil
}

func scanError(l *lexer, resp *agiResp) (scanFunc, error) {
	resp.result = -1
	if resp.code == 520 && l.peek() == '-' {
		return scanErrorUsage, nil
	}
	return scanData, nil
}

func scanErrorUsage(l *lexer, resp *agiResp) (scanFunc, error) {
	// skip hyphen
	l.next()
	l.ignore()

	for {
		char := l.next()
		if char == eof || (char == '\n' && l.lookForward("520 End")) {
			break
		}
	}

	resp.data = l.input[l.start:l.pos]
	return nil, nil
}

// lexer routins
type lexer struct {
	input string
	start int
	pos   int
	width int
}

const eof = -1

func (l *lexer) peek() rune {
	b := l.next()
	l.backup()
	return b
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) lookForward(pattern string) bool {
	pos := l.pos + len(pattern)
	return pos <= len(l.input) && l.input[l.pos:pos] == pattern
}

func (l *lexer) skipWhitespace() int {
	for {
		char := l.next()
		if char == eof {
			return eof
		}
		if !unicode.IsSpace(char) {
			l.backup()
			break
		}
		l.ignore()
	}
	return l.pos
}

func (l *lexer) atoi() int32 {
	s := l.input[l.start:l.pos]
	sign := 1
	if s[0] == '-' {
		sign = -1
		s = s[1:]
	}
	n := 0
	for _, ch := range []byte(s) {
		ch -= '0'
		n = n*10 + int(ch)
	}
	return int32(n * sign)
}
