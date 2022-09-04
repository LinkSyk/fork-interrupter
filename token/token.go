package token

type TokenType string

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	INT = "INT"

	ASSIGN = "="
	EQT    = "=="
	NOTEQT = "!="

	GT = ">"
	LT = "<"

	BANG  = "!"
	PLUS  = "+"
	SUB   = "-"
	MULTI = "*"
	DIV   = "/"

	LPARENT = "("
	RPARENT = ")"
	LBRACE  = "{"
	RBRACE  = "}"

	SEMICOLON = ";"
	COMMA     = ","

	// keywords
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	LET    = "LET"
	IDENT  = "IDENT"
	IF     = "IF"
	ELSE   = "ELSE"
	FN     = "FN"
	RETURN = "RETURN"
)

var keywords = map[string]TokenType{
	"fn":     FN,
	"let":    LET,
	"if":     IF,
	"true":   TRUE,
	"false":  FALSE,
	"else":   ELSE,
	"return": RETURN,
}

type Token struct {
	Type    TokenType
	Literal string
}

func LookIdent(key string) TokenType {
	t, ok := keywords[key]
	if ok {
		return t
	}
	return IDENT
}
