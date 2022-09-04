package lexer

import (
	"interrupter/token"
)

// parse string into token
type Lexer struct {
	input   string
	ch      byte
	pos     int
	readPos int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	// 先要读到第一个字符
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	ch := l.curChar()
	var tok token.Token
	switch ch {
	case '=':
		ch := l.curChar()
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.EQT, string([]byte{ch, l.ch}))
		} else {
			tok = newToken(token.ASSIGN, string(ch))
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.NOTEQT, string([]byte{ch, l.ch}))
		} else {
			tok = newToken(token.BANG, string(ch))
		}
	case '+':
		tok = newToken(token.PLUS, string(ch))
	case '-':
		tok = newToken(token.SUB, string(ch))
	case '*':
		tok = newToken(token.MULTI, string(ch))
	case '/':
		tok = newToken(token.DIV, string(ch))
	case '(':
		tok = newToken(token.LPARENT, string(ch))
	case ')':
		tok = newToken(token.RPARENT, string(ch))
	case '{':
		tok = newToken(token.LBRACE, string(ch))
	case '}':
		tok = newToken(token.RBRACE, string(ch))
	case '<':
		tok = newToken(token.LT, string(ch))
	case '>':
		tok = newToken(token.GT, string(ch))
	case ';':
		tok = newToken(token.SEMICOLON, string(ch))
	case ',':
		tok = newToken(token.COMMA, string(ch))
	case 0:
		return newToken(token.EOF, "")
	default:
		if isNumber(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else if isAlpha(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookIdent(tok.Literal)
			return tok
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = "ILLEGAL"
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isNumber(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isAlpha(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) curChar() byte {
	return l.ch
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == '\r' || l.ch == '\n' || l.ch == '\t' || l.ch == ' ' {
		l.readChar()
	}
}

func newToken(t token.TokenType, l string) token.Token {
	return token.Token{
		Type:    t,
		Literal: l,
	}
}
