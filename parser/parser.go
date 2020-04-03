package parser

import (
	"fmt"
	"monkey-interpreter/ast"
	"monkey-interpreter/lexer"
	"monkey-interpreter/token"
)

// Parser is an implementation of a recursive decent parser that is capable
// of turning tokenized Monkey source code into an abstract syntax tree.
type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	errors []error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: make([]error, 0)}
	// Advance the counter so that it's in a usable state immediately.
	p.advanceToken()
	p.advanceToken()
	return p
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// Checks if the next token matches the specified token and if so advances the parser
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

func (p *Parser) parseLetStatement() ast.Statement {
	letToken := p.currentToken

	if !p.expectAndAdvance(token.IDENTIFIER) {
		return nil
	}

	identifier := &Identifier{Token: p.currentToken, Value: p.currentToken.Value}

	if !p.expectAndAdvance(token.ASSIGN) {
		return nil
	}

	expressionValue := p.parseExpression()

	for p.currentToken.Type != token.SEMICOLON {
		p.advanceToken()
	}

	return LetStatement{
		Token: letToken,
		Name:  identifier,
		Value: expressionValue,
	}
}

func (p *Parser) parseExpression() Expression {
	return nil
}

func (p *Parser) parseReturnStatement() ast.Statement {
	returnToken := p.currentToken
	p.advanceToken()
	expressionValue := p.parseExpression()

	for p.currentToken.Type != token.SEMICOLON {
		p.advanceToken()
	}

	return ReturnStatement{
		Token: returnToken,
		Value: expressionValue,
	}
}
