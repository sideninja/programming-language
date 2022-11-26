package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"language/ast"
	"language/lexer"
	"language/tokens"
	"testing"
)

func parseStatementsWithLen(t *testing.T, input string, statementsLen int) (*Parser, []ast.Statement) {
	l := lexer.New(input)
	p := New(l)
	program, err := p.Parse()
	require.NoError(t, err, fmt.Sprintf("parsing statement %s failed", input))
	require.NotNil(t, program, fmt.Sprintf("parsing statement %s produced nil output", input))
	require.Len(
		t,
		program.Statements,
		statementsLen,
		fmt.Sprintf("parsing input (%s) didn't produce required length of statements: %v", input, program.Statements),
	)

	return p, program.Statements
}

func assertIntegerLiteral(t *testing.T, expression ast.Expression, value int64) {
	lit, ok := expression.(*ast.IntegerLiteral)
	require.True(t, ok)
	assert.Equal(t, value, lit.Value)
}

func Test_LetStatement(t *testing.T) {
	t.Run("successfully parse let statements", func(t *testing.T) {
		input := `
			let foo = 1337;
			let boo = 1000000;
			let x = 2;
		`
		p, statements := parseStatementsWithLen(t, input, 3)
		require.Len(t, p.errors, 0)

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

func Test_ReturnStatement(t *testing.T) {
	t.Run("successfully parse return statements", func(t *testing.T) {
		input := `
			return 1;
			return x;
		`

		p, _ := parseStatementsWithLen(t, input, 2)
		require.Len(t, p.errors, 0)
	})
}

func Test_ProgramStringer(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token:      tokens.Token{Type: tokens.LET, Literal: "let"},
				Identifier: ast.Identifier{Token: tokens.Token{Literal: "x", Type: tokens.IDENTIFIER}, Value: "x"},
				Value:      &ast.Identifier{Token: tokens.Token{Literal: "100", Type: tokens.INT}, Value: "100"},
			},
		},
	}

	assert.Equal(t, "let x = 100;", program.String())
}

func Test_Identifier(t *testing.T) {
	input := `foo;`
	p, statements := parseStatementsWithLen(t, input, 1)
	require.Len(t, p.errors, 0)

	stm, ok := statements[0].(*ast.ExpressionStatement)
	assert.True(t, ok)

	ident, ok := stm.Expression.(*ast.Identifier)
	assert.True(t, ok)

	assert.Equal(t, "foo", ident.Value)
	assert.Equal(t, tokens.IDENTIFIER, string(ident.Token.Type))
}

func Test_Integer(t *testing.T) {
	input := `1337;`
	p, statements := parseStatementsWithLen(t, input, 1)
	require.Len(t, p.errors, 0)

	stm, ok := statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	exp, ok := stm.Expression.(*ast.IntegerLiteral)
	require.True(t, ok)

	assert.Equal(t, 1337, exp.Value)
	assert.Equal(t, "1337", exp.TokenLiteral())
}

func Test_PrefixExpressions(t *testing.T) {
	input := "!foo; -5;"

	testOut := []struct {
		operator string
		literal  string
		out      string
	}{
		{operator: "!", literal: "foo", out: "(!foo)"},
		{operator: "-", literal: "5", out: "(-5)"},
	}

	p, statements := parseStatementsWithLen(t, input, 2)
	require.Len(t, p.errors, 0)

	for i, stm := range statements {
		st, ok := stm.(*ast.ExpressionStatement)
		require.True(t, ok)

		prefix, ok := st.Expression.(*ast.PrefixExpression)
		require.True(t, ok)

		assert.Equal(t, testOut[i].operator, prefix.Operator)
		assert.Equal(t, testOut[i].literal, prefix.Right.TokenLiteral())
		assert.Equal(t, testOut[i].out, prefix.String())
	}
}

func Test_InfixExpressionSimple(t *testing.T) {
	tests := []struct {
		in       string
		left     int64
		operator string
		right    int64
	}{
		{in: "5 + 4", left: 5, operator: "+", right: 4},
		{in: "3 - 4", left: 3, operator: "-", right: 4},
		{in: "2 * 3", left: 2, operator: "*", right: 3},
		{in: "10 / 2", left: 10, operator: "/", right: 2},
		{in: "2 == 3", left: 2, operator: "==", right: 3},
		{in: "3 < 2", left: 3, operator: "<", right: 2},
		{in: "10 > 2", left: 10, operator: ">", right: 2},
		{in: "2 != 2", left: 2, operator: "!=", right: 2},
	}

	for _, test := range tests {
		p, statements := parseStatementsWithLen(t, test.in, 1)
		require.Len(t, p.errors, 0)

		st, ok := statements[0].(*ast.ExpressionStatement)
		require.True(t, ok)

		exp, ok := st.Expression.(*ast.InfixExpression)
		require.True(t, ok)

		assert.Equal(t, test.operator, exp.Operator)
		assertIntegerLiteral(t, exp.Left, test.left)
		assertIntegerLiteral(t, exp.Right, test.right)
	}
}

func Test_InfixSum(t *testing.T) {
	p, statements := parseStatementsWithLen(t, "1 + 2 + 3", 1)
	require.Len(t, p.errors, 0)

	st, ok := statements[0].(*ast.ExpressionStatement)
	require.True(t, ok)

	assert.Equal(t, "((1 + 2) + 3)", st.String())

	exp1, ok := st.Expression.(*ast.InfixExpression)
	require.True(t, ok)

	assert.Equal(t, "+", exp1.Operator)
	assertIntegerLiteral(t, exp1.Right, 3)

	exp2, ok := exp1.Left.(*ast.InfixExpression)
	require.True(t, ok)

	assert.Equal(t, "+", exp2.Operator)
	assertIntegerLiteral(t, exp2.Right, 2)
	assertIntegerLiteral(t, exp2.Left, 1)
}

func Test_InfixExpressionMultiple(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{in: "1 + 2 + 3", out: "((1 + 2) + 3)"},
		{in: "1 - 2 + 3", out: "((1 - 2) + 3)"},
		{in: "1 + 2 * 3", out: "(1 + (2 * 3))"},
		{in: "1 + 2 * 3", out: "(1 + (2 * 3))"},
		{in: "1 < 2 * 3", out: "(1 < (2 * 3))"},
		{in: "1 / 2 * 3", out: "((1 / 2) * 3)"},
	}

	for _, test := range tests {
		p, statements := parseStatementsWithLen(t, test.in, 1)
		require.Len(t, p.errors, 0)

		assert.Equal(t, test.out, statements[0].String())
	}
}

func Test_GroupExpression(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{in: "1 + (2 + 3)", out: "(1 + (2 + 3))"},
		{in: "(1 + 2) * 3", out: "((1 + 2) * 3)"},
		{in: "1 + (2 + 3) + 4", out: "((1 + (2 + 3)) + 4)"},
		{in: "-(5 + 4)", out: "(-(5 + 4))"},
		{in: "!(true == false)", out: "(!(true == false))"},
	}

	for _, test := range tests {
		p, statements := parseStatementsWithLen(t, test.in, 1)
		require.Len(t, p.errors, 0)
		assert.Equal(t, test.out, statements[0].String())
	}
}
