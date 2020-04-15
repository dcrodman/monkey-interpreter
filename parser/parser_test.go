package parser

import (
	"monkey-interpreter/ast"
	"monkey-interpreter/lexer"
	"monkey-interpreter/token"
	"strconv"
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
	//if returnStmt.Value.String() != expectedValue {
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

	testLiteralExpression(t, statement.Expression, "foobar")
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

	testLiteralExpression(t, statement.Expression, 5)
}

func TestBooleanExpression(t *testing.T) {
	inputs := []struct {
		code     string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range inputs {
		p := New(lexer.New(tt.code))
		parsedProgram := p.ParseProgram()
		checkParserHasNoErrors(t, p)

		if len(parsedProgram.Statements) != 1 {
			t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
		}

		statement, ok := parsedProgram.Statements[0].(ast.ExpressionStatement)
		if !ok {
			t.Fatalf("parsedProgram.Statements[0] is not an ast.ExpressionStatement, is: %T", parsedProgram.Statements[0])
		}

		testLiteralExpression(t, statement.Expression, tt.expected)
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
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

		if !testLiteralExpression(t, expr.Right, tt.value) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	prefixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
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

		if !testInfixExpression(t, statement.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	p := New(lexer.New(input))
	parsedProgram := p.ParseProgram()
	checkParserHasNoErrors(t, p)

	if len(parsedProgram.Statements) != 1 {
		t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
	}

	stmt, ok := parsedProgram.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", parsedProgram.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifierLiteral(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	p := New(lexer.New(input))
	parsedProgram := p.ParseProgram()
	checkParserHasNoErrors(t, p)

	if len(parsedProgram.Statements) != 1 {
		t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
	}

	stmt, ok := parsedProgram.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", parsedProgram.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifierLiteral(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifierLiteral(t, alternative.Expression, "y") {
		return
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
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
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

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifierLiteral(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("unknown expected type %T", expected)
	return false
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

	if i.Token.Value != strconv.FormatInt(value, 10) {
		t.Errorf("expected i.TokenLiteral = %d, got = %s", value, i.Token.Value)
		return false
	}

	return true
}

func testIdentifierLiteral(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("expected ident to be *ast.Identifier, got: %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("expected ident.Value to = %s, got = %s", value, ident.Value)
		return false
	}

	if ident.Token.Value != value {
		t.Errorf("expected ident.Token.Value to = %s, got = %s", value, ident.Token.Value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	ident, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("expected ident to be *ast.Boolean, got: %T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("expected ident.Value to = %v, got = %v", value, ident.Value)
		return false
	}

	if ident.Token.Value != strconv.FormatBool(value) {
		t.Errorf("expected ident.Token.Value to = %v, got = %v", value, ident.Token.Value)
		return false
	}

	return true
}
