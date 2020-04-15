package parser

import (
	"fmt"
	"monkey-interpreter/ast"
	"monkey-interpreter/lexer"
	"monkey-interpreter/token"
	"strconv"
)

// Parser is an implementation of a recursive decent parser that is capable
// of turning tokenized Monkey source code into an abstract syntax tree.
type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	errors []error

	infixParseFns  map[token.TokenType]infixParseFn
	prefixParseFns map[token.TokenType]prefixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: make([]error, 0)}
	p.registerParseFns()
	// Advance the counter so that it's in a usable state immediately.
	p.advanceToken()
	p.advanceToken()
	return p
}

// mapping of all prefix and infix operators to the functions that can parse them.
func (p *Parser) registerParseFns() {
	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENTIFIER: p.parseIdentifier,
		token.INT:        p.parseInteger,
		token.BANG:       p.parsePrefixExpression,
		token.MINUS:      p.parsePrefixExpression,
		token.TRUE:       p.parseBoolean,
		token.FALSE:      p.parseBoolean,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NOT_EQ:   p.parseInfixExpression,
		token.LES:      p.parseInfixExpression,
		token.GRT:      p.parseInfixExpression,
	}
}

// ParseProgram uses the Parser's Lexer to advance through the tokenized program
// code and construct an abstract syntax tree.
func (p *Parser) ParseProgram() *ast.AST {
	var statements []ast.Statement

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			statements = append(statements, statement)
		}

		p.advanceToken()
	}

	return &ast.AST{Statements: statements}
}

// Errors returns a slice of errors encountered by the parser during the execution
// of ParseProgram().
func (p *Parser) Errors() []error { return p.errors }

// called for every line in the program (since Monkey is a series of statements) and
// starts the parse tree corresponding to the statement type.
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// checks if the next token matches the specified token and if so advances the parser
// to the next set of tokens. If the token does not match then an error is reported.
func (p *Parser) expectAndAdvance(tokenType token.TokenType) bool {
	if p.nextTokenMatches(tokenType) {
		p.advanceToken()
		return true
	} else {
		p.addExpectedTokenError(tokenType)
		return false
	}
}

// moves the current token and next token each forward by one.
func (p *Parser) advanceToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) nextTokenMatches(tokenType token.TokenType) bool {
	return p.nextToken.Type == tokenType
}

func (p *Parser) addExpectedTokenError(expected token.TokenType) {
	msg := fmt.Errorf("expected token %s, got %s", expected, p.nextToken)
	p.errors = append(p.errors, msg)
}

// parses `let <identifier> = <expression>;` statements.
func (p *Parser) parseLetStatement() ast.Statement {
	letToken := p.currentToken

	if !p.expectAndAdvance(token.IDENTIFIER) {
		return nil
	}

	identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	if !p.expectAndAdvance(token.ASSIGN) {
		return nil
	}

	//expressionValue := p.parseExpression()

	for p.currentToken.Type != token.SEMICOLON {
		p.advanceToken()
	}

	return ast.LetStatement{
		Token: letToken,
		Name:  identifier,
		//Value: expressionValue,
	}
}

// parses `return <expression>;` statements.
func (p *Parser) parseReturnStatement() ast.Statement {
	returnToken := p.currentToken
	p.advanceToken()
	//expressionValue := p.parseExpression()

	for p.currentToken.Type != token.SEMICOLON {
		p.advanceToken()
	}

	return ast.ReturnStatement{
		Token: returnToken,
		//Value: expressionValue,
	}
}

// parses `<expression>;` statements.
func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(LOWEST),
	}

	for p.nextTokenMatches(token.SEMICOLON) {
		p.advanceToken()
	}

	return statement
}

// parse tree for all of Monkey's expressions.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]

	if prefix == nil {
		msg := fmt.Errorf("no prefix parse function for %s found", p.currentToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	leftExpr := prefix()

	for !p.nextTokenMatches(token.SEMICOLON) && precedence < p.nextTokenPrecedence() {
		infix := p.infixParseFns[p.nextToken.Type]
		if infix == nil {
			return leftExpr
		}

		p.advanceToken()
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

// returns the operator precedence for the current token.
func (p *Parser) currentTokenPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

// returns the operator precedence for the next token.
func (p *Parser) nextTokenPrecedence() int {
	if p, ok := precedences[p.nextToken.Type]; ok {
		return p
	}
	return LOWEST
}

// parse a prefix expression (of the type <operator><expression>).
func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
	}

	p.advanceToken()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

// parse an infix expression (of the type <expression> <operator> <expression>).
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Value,
	}

	precedence := p.currentTokenPrecedence()
	p.advanceToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Value,
	}
}

func (p *Parser) parseInteger() ast.Expression {
	intVal, err := strconv.ParseInt(p.currentToken.Value, 0, 64)

	if err != nil {
		msg := fmt.Errorf("unable to parse %q as integer", p.currentToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.Integer{Token: p.currentToken, Value: intVal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentToken.Type == token.TRUE,
	}
}
