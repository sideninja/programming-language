package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"language/ast"
	"language/lexer"
	"testing"
)

func parseStatementsWithLen(t *testing.T, input string, statementsLen int) (*Parser, []ast.Statement) {
	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()
	require.NoError(t, err)
	require.NotNil(t, program)
	assert.Len(t, program.Statements, statementsLen)

	return p, program.Statements
}

func Test_LetStatement(t *testing.T) {
	t.Run("successfully parse let statements", func(t *testing.T) {
		input := `
			let foo = 1337;
			let boo = 1000000;
			let x = 2;
		`
		_, statements := parseStatementsWithLen(t, input, 3)
		identifiers := []string{"foo", "boo", "x"}

		for i, st := range statements {
			assert.Equal(t, "let", st.TokenLiteral())
			assert.IsType(t, &ast.LetStatement{}, st)

			letSt, ok := st.(*ast.LetStatement)
			assert.True(t, ok)
			assert.Equal(t, identifiers[i], letSt.Identifier.Value)
		}
	})

	t.Run("parse let statements with errors", func(t *testing.T) {
		input := `
			let = 1337;
			let boo 1000000;
			let 200;
			let x;
		`
		p, _ := parseStatementsWithLen(t, input, 0)
		require.Len(t, p.errors, 4)

		errors := []string{
			"parsing let statement failed: expected IDENTIFIER, got =",
			"parsing let statement failed: expected =, got INT",
			"parsing let statement failed: expected IDENTIFIER, got INT",
			"parsing let statement failed: expected =, got ;",
		}

		for i, err := range p.errors {
			assert.Equal(t, errors[i], err.Error(), fmt.Sprintf("test case %d failed", i))
		}
	})

}
