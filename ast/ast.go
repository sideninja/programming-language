package ast

import (
	"bytes"
	"fmt"
	"language/tokens"
)

type Node interface {
	fmt.Stringer
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	var buf bytes.Buffer
	for _, s := range p.Statements {
		buf.WriteString(s.TokenLiteral())
	}

	return buf.String()
}

type LetStatement struct {
	Token      tokens.Token
	Identifier Identifier
	Value      Expression
}

func (a *LetStatement) statementNode() {}

func (a *LetStatement) TokenLiteral() string {
	return a.Token.Literal
}

func (a *LetStatement) String() string {
	var val string
	if a.Value != nil {
		val = a.Value.String()
	}

	return fmt.Sprintf(
		"%s %s = %s;",
		a.Token.Literal,
		a.Identifier.String(),
		val,
	)
}

type ReturnStatement struct {
	Token tokens.Token
	Value Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var val string
	if r.Value != nil {
		val = r.Value.String()
	}

	return fmt.Sprintf("%s %s;", r.Token.Literal, val)
}

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token tokens.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}

func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.TokenLiteral()
}

type BooleanLiteral struct {
	Token tokens.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode() {}

func (b *BooleanLiteral) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanLiteral) String() string {
	return b.TokenLiteral()
}

type ExpressionStatement struct {
	Token      tokens.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}

func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

type PrefixExpression struct {
	Token    tokens.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}

func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right)
}

type InfixExpression struct {
	Token    tokens.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}

func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left, i.Operator, i.Right)
}
