package lexer

import (
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	code := `let five = 5;
	let ten = 10; 

	let add = fn(x, y) { 
		x + y; 
	}; 
	
	let result = add( five, ten );
	let bigNumber = 130;
	`
	lexer := NewLexer(code)

	tests := []struct {
		token TokenType
		value string
	}{
		{LET, "let"},
		{IDENTIFIER, "five"},
		{ASSIGN, "="},
		{INT, "5"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENTIFIER, "ten"},
		{ASSIGN, "="},
		{INT, "10"},
		{SEMICOLON, ";"},

		{LET, "let"},
		{IDENTIFIER, "add"},
		{ASSIGN, "="},
		{FUNCTION, "fn"},
		{LPAREN, "("},
		{IDENTIFIER, "x"},
		{COMMA, ","},
		{IDENTIFIER, "y"},
		{RPAREN, ")"},
		{LBRACKET, "{"},
		{IDENTIFIER, "x"},
		{PLUS, "+"},
		{IDENTIFIER, "y"},
		{SEMICOLON, ";"},
		{RBRACKET, "}"},
		{SEMICOLON, ";"},
	}

	for _, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.token || token.Value != test.value {
			t.Fatalf("wanted Token = '%v', Value = '%v'; got Token = '%v', Value = '%v'",
				test.token, test.value, token.Type, token.Value)
		}
	}
}
