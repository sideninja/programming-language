package parser

import (
	"github.com/stretchr/testify/assert"
	"language/ast"
	"language/lexer"
	"language/tokens"
	"testing"
)

func Test_LetStatement(t *testing.T) {
	input := `
		let foo = 1337
		let boo = 1000000
		let x = 2
	`

	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()
	assert.NoError(t, err)

	identifiers := []string{"foo", "boo", "x"}

	for i, st := range program.Statements {
		assert.Equal(t, tokens.LET, st.TokenLiteral())
		assert.IsType(t, &ast.LetStatement{}, st)

		letSt, ok := st.(*ast.LetStatement)
		assert.True(t, ok)
		assert.Equal(t, identifiers[i], letSt.Identifier.Value)
	}
}
