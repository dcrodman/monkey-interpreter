package parser

import (
	"fmt"
	"monkey-interpreter/ast"
	"monkey-interpreter/lexer"
	"monkey-interpreter/token"
	"testing"
)

func checkParserHasNoErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) == 0 {
		return
	}

	for _, err := range p.Errors() {
		t.Error("parser error: ", err)
	}

	t.Fatalf("failed due to unexpected parser errors")
}

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
	`

	tests := []struct {
		name string
		//value string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	p := New(lexer.New(input))
	parsedProgram := p.ParseProgram()
	checkParserHasNoErrors(t, p)

	if len(parsedProgram.Statements) < 3 {
		t.Fatalf("AST contained %d statements, expected 3", len(parsedProgram.Statements))
	}

	for i, tt := range tests {
		statement := parsedProgram.Statements[i]

		if !testLetStatement(t, statement, tt.name) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, expectedIdent string) bool {
	letStmt, ok := s.(ast.LetStatement)

	if !ok {
		t.Errorf("expected LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != expectedIdent {
		t.Errorf("expected let statement to have identifier %s, got %s", expectedIdent, letStmt.Name.Value)
		return false
	}

	if letStmt.Token.Type != token.LET {
		t.Errorf("let statement token does not have LET token type, got: %s", letStmt.Token.Value)
		return false
	}

	// TODO: Value
	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 838383;
	`

	tests := []struct {
		value string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	p := New(lexer.New(input))
	parsedProgram := p.ParseProgram()
	checkParserHasNoErrors(t, p)

	if len(parsedProgram.Statements) < 3 {
		t.Fatalf("AST contained %d statements, expected 3", len(parsedProgram.Statements))
	}

	for i, tt := range tests {
		statement := parsedProgram.Statements[i]

		if !testReturnStatement(t, statement, tt.value) {
			return
		}
	}
}

func testReturnStatement(t *testing.T, s ast.Statement, expectedValue string) bool {
	returnStmt, ok := s.(ast.ReturnStatement)

	if !ok {
		t.Errorf("expected LetStatement, got %T", s)
		return false
	}

	if returnStmt.Token.Type != token.RETURN {
		t.Errorf("let statement token does not have RETURN token type, got: %s", returnStmt.Token.Type)
		return false
	}

	// TODO
	//if returnStmt.Value != expectedValue {
	//	t.Errorf("expected let statement to have identifier %s, got %s", expectedValue, returnStmt.Value)
	//	return false
	//}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	p := New(lexer.New(input))
	parsedProgram := p.ParseProgram()
	checkParserHasNoErrors(t, p)

	if len(parsedProgram.Statements) != 1 {
		t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
	}

	statement, ok := parsedProgram.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("parsedProgram.Statements[0] is not an ast.ExpressionStatement, is: %T", parsedProgram.Statements[0])
	}

	expr, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression in parsedProgram.Statements[0] is not an *ast.Identifier, is: %T", statement.Expression)
	}

	if expr.Token.Type != token.IDENTIFIER {
		t.Fatalf("expected identifier.Token to be an IDENTIFIER type, got %v", expr.Token)
	}

	if expr.Value != "foobar" {
		t.Fatalf("exxpected identifier.Value = foobar, got %v", expr.Token)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := `5;`

	p := New(lexer.New(input))
	parsedProgram := p.ParseProgram()
	checkParserHasNoErrors(t, p)

	if len(parsedProgram.Statements) != 1 {
		t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
	}

	statement, ok := parsedProgram.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("parsedProgram.Statements[0] is not an ast.ExpressionStatement, is: %T", parsedProgram.Statements[0])
	}

	expr, ok := statement.Expression.(*ast.Integer)
	if !ok {
		t.Fatalf("expression in parsedProgram.Statements[0] is not an *ast.Integer, is: %T", statement.Expression)
	}

	if expr.Token.Type != token.INT {
		t.Fatalf("expected identifier.Token to be an INT type, got %v", expr.Token)
	}

	if expr.Value != 5 {
		t.Fatalf("exxpected identifier.Value = 5, got %v", expr.Token)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		p := New(lexer.New(tt.input))
		parsedProgram := p.ParseProgram()
		checkParserHasNoErrors(t, p)

		if len(parsedProgram.Statements) != 1 {
			t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
		}

		statement, ok := parsedProgram.Statements[0].(ast.ExpressionStatement)
		if !ok {
			t.Fatalf("parsedProgram.Statements[0] is not an ast.ExpressionStatement, is: %T", statement)
		}

		expr, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expression in parsedProgram.Statements[0] is not an *ast.PrefixExpression, is: %T", statement.Expression)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("expected expr.Operator to = %s, got: %s", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, ie ast.Expression, value int64) bool {
	i, ok := ie.(*ast.Integer)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", ie)
		return false
	}

	if i.Value != value {
		t.Errorf("exptected i.Value = %d, got = %d", value, i.Value)
		return false
	}

	if i.Token.Value != fmt.Sprintf("%d", value) {
		t.Errorf("expected i.TokenLiteral = %d, got = %s", value, i.Token.Value)
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	prefixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range prefixTests {
		p := New(lexer.New(tt.input))
		parsedProgram := p.ParseProgram()
		checkParserHasNoErrors(t, p)

		if len(parsedProgram.Statements) != 1 {
			t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
		}

		statement, ok := parsedProgram.Statements[0].(ast.ExpressionStatement)
		if !ok {
			t.Fatalf("parsedProgram.Statements[0] is not an ast.ExpressionStatement, is: %T", statement)
		}

		expr, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expression in parsedProgram.Statements[0] is not an *ast.InfixExpression, is: %T", statement.Expression)
		}

		if !testIntegerLiteral(t, expr.Left, tt.leftValue) {
			return
		}

		if expr.Operator != tt.operator {
			t.Fatalf("expected expr.Operator to = %s, got: %s", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		parsedProgram := p.ParseProgram()
		checkParserHasNoErrors(t, p)

		actual := parsedProgram.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
