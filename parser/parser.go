package parser

import (
	"fmt"
	"language/ast"
	"language/lexer"
	"language/tokens"
	"strconv"
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

	tokenPrecedences map[tokens.TokenType]int
}

func New(l *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: l,

		prefixParsers: make(map[tokens.TokenType]prefixParse),
		infixParsers:  make(map[tokens.TokenType]infixParse),

		tokenPrecedences: map[tokens.TokenType]int{
			tokens.MINUS:    SUM,
			tokens.PLUS:     SUM,
			tokens.LESS:     LESSGREATER,
			tokens.GREATER:  LESSGREATER,
			tokens.EQUAL:    EQUALS,
			tokens.MULTIPLY: PRODUCT,
			tokens.DIVIDE:   PRODUCT,
		},
	}

	parser.registerPrefix(tokens.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefix(tokens.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(tokens.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(tokens.MINUS, parser.parsePrefixExpression)

	parser.registerInfix(tokens.PLUS, parser.parseInfixExpression)
	parser.registerInfix(tokens.MINUS, parser.parseInfixExpression)
	parser.registerInfix(tokens.DIVIDE, parser.parseInfixExpression)
	parser.registerInfix(tokens.MULTIPLY, parser.parseInfixExpression)
	parser.registerInfix(tokens.EQUAL, parser.parseInfixExpression)
	parser.registerInfix(tokens.NOTEQUAL, parser.parseInfixExpression)
	parser.registerInfix(tokens.LESS, parser.parseInfixExpression)
	parser.registerInfix(tokens.GREATER, parser.parseInfixExpression)

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

	switch p.token.Type {
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		st = p.parseReturnStatement()
	default:
		st = p.parseExpressionStatement()
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

func (p *Parser) peekPrecedence() int {
	return p.tokenPrecedences[p.peekToken.Type]
}

func (p *Parser) lookupPrecedence(token tokens.Token) int {
	return p.tokenPrecedences[token.Type]
}

func (p *Parser) expectPeekType(t tokens.TokenType) bool {
	if !p.isPeekType(t) {
		p.addParseError(fmt.Errorf("expected %s, got %s", t, p.peekToken.Type))
		return false
	}

	p.nextToken()
	return true
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

/*
	Expression parsing
*/

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	st := &ast.ExpressionStatement{
		Token:      p.token,
		Expression: p.parseExpression(LOWEST),
	}

	if p.isPeekType(tokens.SEMICOLON) {
		p.nextToken()
	}

	return st
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	parser, exists := p.prefixParsers[p.token.Type]
	if !exists {
		p.addParseError(fmt.Errorf("no prefix parser found for token %s", p.token))
		return nil
	}
	leftExpr := parser()

	for precedence < p.peekPrecedence() && !p.isPeekType(tokens.SEMICOLON) {
		infixParser, exists := p.infixParsers[p.peekToken.Type]
		if !exists {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infixParser(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefix := &ast.PrefixExpression{
		Token:    p.token,
		Operator: p.token.Literal,
	}

	p.nextToken()
	prefix.Right = p.parseExpression(PREFIX)

	return prefix
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infix := &ast.InfixExpression{
		Token:    p.token,
		Left:     left,
		Operator: p.token.Literal,
	}

	precedence := p.lookupPrecedence(p.token)

	p.nextToken()
	infix.Right = p.parseExpression(precedence)

	return infix
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	integer := &ast.IntegerLiteral{
		Token: p.token,
	}

	val, err := strconv.ParseInt(p.token.Literal, 10, 64)
	if err != nil {
		return nil
	}

	integer.Value = val
	return integer
}

/*
	Statement parsing
*/

func (p *Parser) parseLetStatement() ast.Statement {
	st := &ast.LetStatement{
		Token: p.token, // let
	}

	if !p.expectPeekType(tokens.IDENTIFIER) {
		return nil
	}

	st.Identifier = ast.Identifier{
		Token: p.token,
		Value: p.token.Literal,
	}

	if !p.expectPeekType(tokens.ASSIGN) {
		return nil
	}

	p.nextToken() // assign =

	st.Value = p.parseExpression(LOWEST)

	for p.token.Type != tokens.SEMICOLON {
		p.nextToken()
	}

	return st
}

func (p *Parser) parseReturnStatement() ast.Statement {
	st := &ast.ReturnStatement{
		Token: p.token,
	}

	for p.token.Type != tokens.SEMICOLON {
		p.nextToken()
	}

	return st
}
