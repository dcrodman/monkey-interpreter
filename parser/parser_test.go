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
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		parsedProgram := p.ParseProgram()
		checkParserHasNoErrors(t, p)

		if len(parsedProgram.Statements) != 1 {
			t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
		}

		letStmt, ok := parsedProgram.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("expected LetStatement, got %T", parsedProgram.Statements[0])
		}

		if letStmt.Token.Type != token.LET {
			t.Fatalf("let statement token does not have LET token type, got: %s", letStmt.Token.Value)
		}

		if !testIdentifierLiteral(t, letStmt.Name, tt.expectedIdentifier) {
			return
		}

		if !testLiteralExpression(t, letStmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return x;", "x"},
		//{"return add(x, y);", "add(x, y)"},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		parsedProgram := p.ParseProgram()
		checkParserHasNoErrors(t, p)

		if len(parsedProgram.Statements) != 1 {
			t.Fatalf("AST contained %d statements, expected 1", len(parsedProgram.Statements))
		}

		returnStmt, ok := parsedProgram.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("expected ReturnStatement, got %T", parsedProgram.Statements[0])
		}

		if returnStmt.Token.Type != token.RETURN {
			t.Fatalf("let statement token does not have RETURN token type, got: %s", returnStmt.Token.Type)
		}

		if !testLiteralExpression(t, returnStmt.Value, tt.expectedValue) {
			return
		}
	}
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

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	p := New(lexer.New(input))
	parsedProgram := p.ParseProgram()
	checkParserHasNoErrors(t, p)

	if len(parsedProgram.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(parsedProgram.Statements))
	}

	stmt, ok := parsedProgram.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", parsedProgram.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifierLiteral(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "add();", expectedParams: []string{}},
		{input: "add(x, y, z);", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		parsedProgram := p.ParseProgram()
		checkParserHasNoErrors(t, p)

		stmt := parsedProgram.Statements[0].(ast.ExpressionStatement)
		callExp := stmt.Expression.(*ast.CallExpression)

		if len(callExp.Arguments) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(callExp.Arguments))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, callExp.Arguments[i], ident)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

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

	function, ok := stmt.Expression.(*ast.Function)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		parsedProgram := p.ParseProgram()
		checkParserHasNoErrors(t, p)

		stmt := parsedProgram.Statements[0].(ast.ExpressionStatement)
		function := stmt.Expression.(*ast.Function)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
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
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
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
