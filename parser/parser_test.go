package parser

import (
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
	letStmt, ok := s.(LetStatement)

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
	returnStmt, ok := s.(ReturnStatement)

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
