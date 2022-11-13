package parser

import (
	"fmt"
	"language/ast"
	"language/lexer"
	"language/tokens"
)

type Parser struct {
	lexer *lexer.Lexer

	token     tokens.Token
	peekToken tokens.Token
}

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{lexer: l}

	parser.nextToken()
	parser.nextToken()

	return parser
}

func (p *Parser) nextToken() {
	p.token = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}

	for p.token != tokens.EOFToken {
		var st ast.Statement
		var err error

		switch p.token.Type {
		case tokens.LET:
			st, err = p.parseLetStatement()
			if err != nil {
				return nil, fmt.Errorf("parsing let statement failed: %w", err)
			}
		}

		if st != nil {
			program.Statements = append(program.Statements, st)
		}

		p.nextToken()
	}

	return program, nil
}

func (p *Parser) parseExpression() {
	// todo
}

func (p *Parser) parseLetStatement() (ast.Statement, error) {
	letToken := p.token

	if p.peekToken.Type != tokens.IDENTIFIER {
		return nil, fmt.Errorf("expected identifier, got %s", p.peekToken.Type)
	}

	p.nextToken() // identifier

	identifier := ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}

	if p.peekToken.Type != tokens.ASSIGN {
		return nil, fmt.Errorf("expected assignment, got %s", p.peekToken.Type)
	}

	p.nextToken() // assign =

	p.parseExpression()

	return &ast.LetStatement{
		Token:      letToken,
		Identifier: identifier,
		Right:      nil, // todo fil
	}, nil
}
