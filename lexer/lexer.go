package lexer

import (
	"bytes"
	"language/tokens"
)

const EOF = 0

type symbol byte

func (s symbol) String() string {
	return string(s)
}

type Lexer struct {
	input   string
	pos     int
	nextPos int
	symbol  symbol
}

func New(input string) *Lexer {
	lex := &Lexer{
		input: input,
	}
	lex.readChar()
	return lex
}

func (l *Lexer) NextToken() tokens.Token {
	var token tokens.Token
	l.skipWhitespace()

	switch l.symbol {
	case '+':
		token = tokens.New(l.symbol.String(), tokens.PLUS)
	case '-':
		token = tokens.New(l.symbol.String(), tokens.MINUS)
	case '*':
		token = tokens.New(l.symbol.String(), tokens.MULTIPLY)
	case '\\':
		token = tokens.New(l.symbol.String(), tokens.DIVIDE)
	case '(':
		token = tokens.New(l.symbol.String(), tokens.LPAREN)
	case ')':
		token = tokens.New(l.symbol.String(), tokens.RPAREN)
	case '{':
		token = tokens.New(l.symbol.String(), tokens.LBRACE)
	case '}':
		token = tokens.New(l.symbol.String(), tokens.RBRACE)
	case ',':
		token = tokens.New(l.symbol.String(), tokens.COMMA)
	case '!':
		if l.peakNext() == '=' {
			l.readChar()
			token = tokens.New("!=", tokens.NOTEQUAL)
		} else {
			token = tokens.New("!", tokens.NEGATIVE)
		}
	case '=':
		if l.peakNext() == '=' {
			l.readChar()
			token = tokens.New("==", tokens.EQUAL)
		} else {
			token = tokens.New(l.symbol.String(), tokens.ASSIGN)
		}
	case EOF:
		token = tokens.New(l.symbol.String(), tokens.EOF)
	default:
		if l.isChar() {
			ident := l.readIdentifier()
			return tokens.New(ident, tokens.LookupIdentifier(ident))
		}
		if l.isNumber() {
			return tokens.New(l.readInteger(), tokens.INT)
		}
	}

	l.readChar()
	return token
}

func (l *Lexer) skipWhitespace() {
	for l.symbol == ' ' || l.symbol == '\r' || l.symbol == '\t' || l.symbol == '\n' {
		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.symbol = EOF
	} else {
		l.symbol = symbol(l.input[l.nextPos])
	}

	l.pos = l.nextPos
	l.nextPos += 1
}

func (l *Lexer) peakNext() byte {
	if l.nextPos >= len(l.input) {
		return EOF
	}

	return l.input[l.nextPos]
}

func (l *Lexer) readIdentifier() string {
	var ident bytes.Buffer
	for l.isChar() {
		ident.WriteByte(byte(l.symbol))
		l.readChar()
	}
	return ident.String()
}

func (l *Lexer) readInteger() string {
	var integer bytes.Buffer
	for l.isNumber() {
		integer.WriteByte(byte(l.symbol))
		l.readChar()
	}

	return integer.String()
}

func (l *Lexer) isChar() bool {
	return l.symbol >= 'a' && l.symbol <= 'z' || l.symbol >= 'A' && l.symbol <= 'Z'
}

func (l *Lexer) isNumber() bool {
	return l.symbol >= '0' && l.symbol <= '9'
}
