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
