package lexer

import (
	"interrupter/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadChar(t *testing.T) {
	input := `let a = 1`
	tables := []struct {
		output byte
	}{
		{'l'},
		{'e'},
		{'t'},
		{' '},
		{'a'},
		{' '},
		{'='},
		{' '},
		{'1'},
		{0},
	}

	l := New(input)
	for _, table := range tables {
		l.readChar()
		assert.Equal(t, table.output, l.ch)
	}
}

func TestParseToken(t *testing.T) {
	input := `

	 +-/* ,;( )123{}let==abcd!=<>
	`
	tables := []token.Token{
		newToken(token.PLUS, "+"),
		newToken(token.SUB, "-"),
		newToken(token.DIV, "/"),
		newToken(token.MULTI, "*"),
		newToken(token.COMMA, ","),
		newToken(token.SEMICOLON, ";"),
		newToken(token.LPARENT, "("),
		newToken(token.RPARENT, ")"),
		newToken(token.INT, "123"),
		newToken(token.LBRACE, "{"),
		newToken(token.RBRACE, "}"),
		newToken(token.IDENT, "let"),
		newToken(token.EQT, "=="),
		newToken(token.IDENT, "abcd"),
		newToken(token.NOTEQT, "!="),
		newToken(token.LT, "<"),
		newToken(token.GT, ">"),
		newToken(token.EOF, ""),
	}
	l := New(input)
	for _, tb := range tables {
		tk := l.NextToken()
		assert.Equal(t, tk.Type, tb.Type)
		assert.Equal(t, tk.Literal, tb.Literal)
	}
}
