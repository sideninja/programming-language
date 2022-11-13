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

	errors []error
}

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: l,
	}

	// we fill current token and peek token, so they are not empty
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

	for !p.isPeekType(tokens.EOF) {
		st := p.parseStatement()
		if st != nil {
			program.Statements = append(program.Statements, st)
		}

		p.nextToken()
	}

	return program, nil
}

func (p *Parser) parseStatement() ast.Statement {

	switch p.token.Type {
	case tokens.LET:
		st, err := p.parseLetStatement()
		if err != nil {
			p.addParseError(fmt.Errorf("parsing let statement failed: %w", err))
		}
		return st
	default:
		p.addParseError(fmt.Errorf("invalid parse statement: %v", p.token))
		return nil
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

	p.expectPeekType(tokens.IDENTIFIER)
	p.nextToken() // identifier

	st.Identifier = ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}

	p.expectPeekType(tokens.ASSIGN)
	p.nextToken() // assign =

	st.Right = p.parseExpression()

	for p.token.Type != tokens.SEMICOLON {
		p.nextToken()
	}

	return st, nil
}

func (p *Parser) isPeekType(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeekType(t tokens.TokenType) {
	if !p.isPeekType(t) {
		p.addParseError(fmt.Errorf("expected %s, got %s", t, p.peekToken.Type))
	}
}

func (p *Parser) addParseError(err error) {
	p.errors = append(p.errors, err)
}
