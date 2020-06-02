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
		token.LPAREN:     p.parseGroupedExpression,
		token.IF:         p.parseIfExpression,
		token.FUNCTION:   p.parseFunction,
		token.STRING:     p.parseStringLiteral,
		token.LBRACKET:   p.parseArrayLiteral,
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
		token.LPAREN:   p.parseCallExpression,
		token.LBRACKET: p.parseIndexExpression,
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
	if p.nextTokenIs(tokenType) {
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

func (p *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) nextTokenIs(tokenType token.TokenType) bool {
	return p.nextToken.Type == tokenType
}

func (p *Parser) addExpectedTokenError(expected token.TokenType) {
	msg := fmt.Errorf("expected token %s, got %s", expected, p.nextToken)
	p.errors = append(p.errors, msg)
}

// parses `let <identifier> = <expression>;` statements.
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectAndAdvance(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	if !p.expectAndAdvance(token.ASSIGN) {
		return nil
	}

	p.advanceToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.nextTokenIs(token.SEMICOLON) {
		p.advanceToken()
	}

	return stmt
}

// parses `return <expression>;` statements.
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.advanceToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.nextTokenIs(token.SEMICOLON) {
		p.advanceToken()
	}

	return stmt
}

// parses `<expression>;` statements.
func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(LOWEST),
	}

	for p.nextTokenIs(token.SEMICOLON) {
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

	for !p.nextTokenIs(token.SEMICOLON) && precedence < p.nextTokenPrecedence() {
		infix := p.infixParseFns[p.nextToken.Type]
		if infix == nil {
			return leftExpr
		}

		p.advanceToken()
		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectAndAdvance(token.RPAREN) {
		return nil
	}

	return exp
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
		Value: p.currentTokenIs(token.TRUE),
	}
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.currentToken}

	if !p.expectAndAdvance(token.LPAREN) {
		return nil
	}

	p.advanceToken()
	exp.Condition = p.parseExpression(LOWEST)

	// Consume the expected tokens up until we start the conditional block.
	if !p.expectAndAdvance(token.RPAREN) {
		return nil
	} else if !p.expectAndAdvance(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.nextTokenIs(token.ELSE) {
		p.advanceToken()

		if !p.expectAndAdvance(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.currentToken,
		Statements: []ast.Statement{},
	}

	p.advanceToken()

	// Read statements until we hit the end of the block (or the file). Conceptually
	// this is how the top-level parser loop iterates over the code, just in the
	// context of a specific block rather than the entire program.
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.advanceToken()
	}

	return block
}

func (p *Parser) parseFunction() ast.Expression {
	f := &ast.Function{Token: p.currentToken}

	if !p.expectAndAdvance(token.LPAREN) {
		return nil
	}

	f.Parameters = p.parseFunctionParameters()

	if !p.expectAndAdvance(token.LBRACE) {
		return nil
	}

	f.Body = p.parseBlockStatement()

	return f
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.nextTokenIs(token.RPAREN) {
		p.advanceToken()
		return identifiers
	}

	p.advanceToken()

	identifiers = append(identifiers, &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Value,
	})

	// Continue reading the list of parameters until we hit the ).
	for p.nextTokenIs(token.COMMA) {
		p.advanceToken()
		p.advanceToken()

		identifiers = append(identifiers, &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Value,
		})
	}

	if !p.expectAndAdvance(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token:     p.currentToken,
		Function:  function,
		Arguments: p.parseExpressionList(token.RPAREN),
	}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.String{Token: p.currentToken, Value: p.currentToken.Value}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	elements := p.parseExpressionList(token.RBRACKET)
	return &ast.Array{Token: p.currentToken, Elements: elements}
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.nextTokenIs(end) {
		p.advanceToken()
		return list
	}

	p.advanceToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.nextTokenIs(token.COMMA) {
		p.advanceToken()
		p.advanceToken()

		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectAndAdvance(end) {
		return nil
	}

	return list
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.advanceToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectAndAdvance(token.RBRACKET) {
		return nil
	}

	return exp
}
