package tokens

type TokenType string

/*
let x = 1;
fun foo() {}
if x == 2 {}
return

*/

const (
	IDENTIFIER = "IDENT"
	INT        = "INT"
	SEMICOLON  = ";"
	ASSIGN     = "="
	PLUS       = "+"
	MINUS      = "-"
	MULTIPLY   = "*"
	DIVIDE     = "/"
	EQUAL      = "=="
	NOTEQUAL   = "!="
	LPAREN     = "("
	RPAREN     = ")"
	NEGATIVE   = "!"
	LBRACE     = "{"
	RBRACE     = "}"
	COMMA      = ","
	SPACE      = " "
	EOF        = ""
	INVALID    = "INVALID"

	// Keywords
	LET    = "LET"
	FUN    = "FUN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	RETURN = "RETURN"
)

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
