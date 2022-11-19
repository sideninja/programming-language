package parser

import (
	"fmt"
	"language/ast"
	"language/lexer"
	"language/tokens"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X !X
	CALL        // foo()
)

type (
	prefixParse func() ast.Expression
	infixParse  func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer *lexer.Lexer

	token     tokens.Token
	peekToken tokens.Token

	errors []error

	prefixParsers map[tokens.TokenType]prefixParse
	infixParsers  map[tokens.TokenType]infixParse
}

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: l,

		prefixParsers: make(map[tokens.TokenType]prefixParse),
		infixParsers:  make(map[tokens.TokenType]infixParse),
	}

	parser.registerPrefix(tokens.IDENTIFIER, parser.parseIdentifier)

	// we fill current token and peek token, so they are not empty
	parser.nextToken()
	parser.nextToken()

	return parser
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
	var st ast.Statement
	var err error

	switch p.token.Type {
	case tokens.LET:
		st, err = p.parseLetStatement()
		if err != nil {
			p.addParseError(fmt.Errorf("parsing let statement failed: %w", err))
		}
		return st
	case tokens.RETURN:
		st, err = p.parseReturnStatement()
		if err != nil {
			p.addParseError(fmt.Errorf("parsing return statement failed: %w", err))
		}
	default:
		st, err = p.parseExpressionStatement()
		if err != nil {
			p.addParseError(fmt.Errorf("parsing expression failed: %w", err))
		}
	}

	return st
}

func (p *Parser) nextToken() {
	p.token = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) isPeekType(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeekType(t tokens.TokenType) error {
	if !p.isPeekType(t) {
		return fmt.Errorf("expected %s, got %s", t, p.peekToken.Type)
	}

	return nil
}

func (p *Parser) registerPrefix(token tokens.TokenType, parser prefixParse) {
	p.prefixParsers[token] = parser
}

func (p *Parser) registerInfix(token tokens.TokenType, parser infixParse) {
	p.infixParsers[token] = parser
}

func (p *Parser) addParseError(err error) {
	p.errors = append(p.errors, err)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	st := &ast.ExpressionStatement{
		Token: p.token,
	}

	st.Expression = p.parseExpression(LOWEST)

	if p.isPeekType(tokens.SEMICOLON) {
		p.nextToken()
	}

	return st, nil
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parser, registered := p.prefixParsers[p.token.Type]
	if !registered {
		return nil
	}
	leftExpr := parser()

	return leftExpr
}

func (p *Parser) parseLetStatement() (ast.Statement, error) {
	st := &ast.LetStatement{
		Token: p.token, // let
	}

	if err := p.expectPeekType(tokens.IDENTIFIER); err != nil {
		return nil, err
	}
	p.nextToken() // identifier

	st.Identifier = ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}

	if err := p.expectPeekType(tokens.ASSIGN); err != nil {
		return nil, err
	}

	p.nextToken() // assign =

	st.Value = p.parseExpression(LOWEST)

	for p.token.Type != tokens.SEMICOLON {
		p.nextToken()
	}

	return st, nil
}

func (p *Parser) parseReturnStatement() (ast.Statement, error) {
	st := &ast.ReturnStatement{
		Token: p.token,
	}

	for p.token.Type != tokens.SEMICOLON {
		p.nextToken()
	}

	return st, nil
}
