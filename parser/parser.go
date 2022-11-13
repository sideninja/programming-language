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
		st, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		program.Statements = append(program.Statements, st)
		p.nextToken()
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.token.Type {
	case tokens.LET:
		st, err := p.parseLetStatement()
		if err != nil {
			return nil, fmt.Errorf("parsing let statement failed: %w", err)
		}
		return st, nil
	default:
		return nil, fmt.Errorf("invalid statement")
	}
}

func (p *Parser) parseExpression() ast.Expression {
	// todo
	return nil
}

func (p *Parser) parseLetStatement() (ast.Statement, error) {
	st := &ast.LetStatement{
		Token: p.token, // let
	}

	if err := p.mustPeekType(tokens.IDENTIFIER); err != nil {
		return nil, err
	}

	p.nextToken() // identifier

	st.Identifier = ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}

	if err := p.mustPeekType(tokens.ASSIGN); err != nil {
		return nil, err
	}

	p.nextToken() // assign =

	st.Right = p.parseExpression()

	for p.token.Type != tokens.SEMICOLON {
		p.nextToken()
	}

	return st, nil
}

func (p *Parser) mustPeekType(t tokens.TokenType) error {
	if p.peekToken.Type != t {
		return fmt.Errorf("expected %s, got %s", t, p.peekToken.Type)
	}
	return nil
}
