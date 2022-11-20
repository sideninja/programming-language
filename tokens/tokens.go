package tokens

type TokenType string

const (
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	SEMICOLON  = ";"
	ASSIGN     = "="
	PLUS       = "+"
	MINUS      = "-"
	MULTIPLY   = "*"
	DIVIDE     = "/"
	EQUAL      = "=="
	NOTEQUAL   = "!="
	LESS       = "<"
	GREATER    = ">"
	LPAREN     = "("
	RPAREN     = ")"
	BANG       = "!"
	LBRACE     = "{"
	RBRACE     = "}"
	COMMA      = ","
	SPACE      = " "
	EOF        = ""
	INVALID    = "INVALID"

	// keywords
	LET    = "LET"
	FUN    = "FUN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	RETURN = "RETURN"
)

var EOFToken = Token{
	Literal: EOF,
	Type:    EOF,
}

var keywords = map[string]TokenType{
	"fun":    FUN,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

type Token struct {
	Literal string
	Type    TokenType
}

func New(literal string, t TokenType) Token {
	return Token{
		Literal: literal,
		Type:    t,
	}
}

func LookupIdentifier(ident string) TokenType {
	if token, ok := keywords[ident]; ok {
		return token
	}
	return IDENTIFIER
}
