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

func cmdParse(text string) (*agiResp, error) {
	resp := &agiResp{}
	lex := &lexer{input: text}

	scanner, err := scanCode(lex, resp)
	if err != nil {
		return nil, err
	}
	// state machine
	for {
		scanner, err := scanner(lex, resp)
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
		if char == eof {
			return nil, EInvalResp.withInfo("scanCode:eof:" + l.input)
		}
		if !unicode.IsDigit(char) {
			l.backup()
			break
		}
		char -= '0'
		resp.code = resp.code*10 + int(char)
	}
	l.start = l.pos
	return scanResult, nil
}

func scanResult(l *lexer, resp *agiResp) (scanFunc, error) {
	pattern := "result="
	char := l.next()
	if char != ' ' {
		return nil, EInvalResp.withInfo("scanResult:space expected:" + l.input)
	}
	l.ignore()

	if pattern != l.input[l.pos:l.pos+len(pattern)] {
		return nil, EInvalResp.withInfo("scanResult:result= expected:" + l.input)
	}
	l.pos += len(pattern)
	l.start = l.pos
	if l.pos > len(l.input) {
		return nil, EInvalResp.withInfo("scanResult:eof:" + l.input)
	}

	for {
		char = l.next()
		if char == eof {
			return nil, EInvalResp.withInfo("scanResult:eof:" + l.input)
		}
		if unicode.IsSpace(char) {
			l.backup()
			break
		}
	}
	resp.result = l.input[l.start:l.pos]
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
	if l.pos > len(l.input) {
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

/*

	pos := 1
	d := int(text[0] - '0')
	if d > 9 {
		return nil, EInvalResp.withInfo(":code:" + text)
	}
	resp.code = d

	// scan response code
	for _, ch := range []byte(text[1:]) {
		ch -= '0'
		if ch > 9 {
			break
		}
		resp.code = resp.code*10 + int(ch)
		pos++
	}

	if text[pos] != ' ' {
		return nil, EInvalResp.withInfo(":code space:" + text)
	}
	pos++

	// scan result=
	if text[pos:pos+7] != "result=" {
		return nil, EInvalResp.withInfo(":code space result=:" + text)
	}
	text = text[pos+7:]
	pos = 0
	for _, ch := range []byte(text[pos:]) {
		if ch == ' ' || ch == '\n' {
			break
		}
		pos++
	}
	resp.result = text[:pos]
	if text[pos] != '\n' {
		resp.data = text[pos:len(text)]
	}

	return resp, nil
}
*/
