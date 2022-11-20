package lexer

import (
	"fmt"
	"language/tokens"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readAllTokens(lexer *Lexer) []tokens.Token {
	all := make([]tokens.Token, 0)
	for token := lexer.NextToken(); token.Type != tokens.EOF; token = lexer.NextToken() {
		all = append(all, token)
	}
	return all
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		in  string
		out []tokens.Token
	}{{
		in: "1 + 2",
		out: []tokens.Token{
			{Literal: "1", Type: tokens.INT},
			{Literal: "+", Type: tokens.PLUS},
			{Literal: "2", Type: tokens.INT},
		},
	}, {
		in: "1 + 2 * 3 / (let foo = 0)",
		out: []tokens.Token{
			{Literal: "1", Type: tokens.INT},
			{Literal: "+", Type: tokens.PLUS},
			{Literal: "2", Type: tokens.INT},
			{Literal: "*", Type: tokens.MULTIPLY},
			{Literal: "3", Type: tokens.INT},
			{Literal: "/", Type: tokens.DIVIDE},
			{Literal: "(", Type: tokens.LPAREN},
			{Literal: "let", Type: tokens.LET},
			{Literal: "foo", Type: tokens.IDENTIFIER},
			{Literal: "=", Type: tokens.ASSIGN},
			{Literal: "0", Type: tokens.INT},
			{Literal: ")", Type: tokens.RPAREN},
		},
	}, {
		in: "== !=",
		out: []tokens.Token{
			{Literal: "==", Type: tokens.EQUAL},
			{Literal: "!=", Type: tokens.NOTEQUAL},
		},
	}, {
		in: "let add = fun(x, y) { x + y }",
		out: []tokens.Token{
			{Literal: "let", Type: tokens.LET},
			{Literal: "add", Type: tokens.IDENTIFIER},
			{Literal: "=", Type: tokens.ASSIGN},
			{Literal: "fun", Type: tokens.FUN},
			{Literal: "(", Type: tokens.LPAREN},
			{Literal: "x", Type: tokens.IDENTIFIER},
			{Literal: ",", Type: tokens.COMMA},
			{Literal: "y", Type: tokens.IDENTIFIER},
			{Literal: ")", Type: tokens.RPAREN},
			{Literal: "{", Type: tokens.LBRACE},
			{Literal: "x", Type: tokens.IDENTIFIER},
			{Literal: "+", Type: tokens.PLUS},
			{Literal: "y", Type: tokens.IDENTIFIER},
			{Literal: "}", Type: tokens.RBRACE},
		},
	}}

	for i, test := range tests {
		lexer := New(test.in)
		all := readAllTokens(lexer)
		assert.Equal(t, test.out, all, fmt.Sprintf("test number: %d failed", i))
	}
}
