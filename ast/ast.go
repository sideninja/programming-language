package ast

import (
	"bytes"
	"language/tokens"
)

type Node interface {
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
	Right      Expression
}

func (a *LetStatement) statementNode() {}

func (a *LetStatement) TokenLiteral() string {
	return a.Token.Literal
}

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
