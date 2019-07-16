package goagi

import (
	"unicode"
	"unicode/utf8"
)

// EInvalResp error returns when AGI response does not match pattern
var EInvalResp = errorNew("Invalid AGI response")

type agiResp struct {
	code   int
	result string
	data   string
}

func parseResponse(text string) (*agiResp, error) {
	resp := &agiResp{}
	lex := &lexer{input: text}

	// state machine
	for scanner, err := scanCode(lex, resp); ; scanner, err = scanner(lex, resp) {
		if err != nil {
			return nil, err
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
		char -= '0'
		resp.code = resp.code*10 + int(char)
	}
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

	if !l.hasPrefix(pattern) {
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
		return nil, EInvalResp.withInfo("scanResult:empty result:" + l.input)
	}
	resp.result = l.input[l.start:l.pos]
	return scanData, nil
}

func scanData(l *lexer, resp *agiResp) (scanFunc, error) {
	// skip spaces in the begining
	for {
		char := l.next()
		if char == eof {
			return nil, nil
		}
		if !unicode.IsSpace(char) {
			break
		}
		l.ignore()
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
	resp.result = "-1"
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
		if char == eof || (char == '\n' && l.hasPrefix("520 End")) {
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

func (l *lexer) hasPrefix(pattern string) bool {
	pos := l.pos + len(pattern)
	return pos <= len(l.input) && l.input[l.pos:pos] == pattern
}
