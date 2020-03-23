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

	!-/*5;
	5 < 10 > 5; 

	if (5 < 10) { 
		return true; 
	} else { 
		return false; 
	}

	10 == 10;
	10 != 9;
	`
	lexer := New(code)

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

		{LET, "let"},
		{IDENTIFIER, "result"},
		{ASSIGN, "="},
		{IDENTIFIER, "add"},
		{LPAREN, "("},
		{IDENTIFIER, "five"},
		{COMMA, ","},
		{IDENTIFIER, "ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},

		{BANG, "!"},
		{MINUS, "-"},
		{SLASH, "/"},
		{ASTERISK, "*"},
		{INT, "5"},
		{SEMICOLON, ";"},

		{INT, "5"},
		{LES, "<"},
		{INT, "10"},
		{GRT, ">"},
		{INT, "5"},
		{SEMICOLON, ";"},

		{IF, "if"},
		{LPAREN, "("},
		{INT, "5"},
		{LES, "<"},
		{INT, "10"},
		{RPAREN, ")"},
		{LBRACKET, "{"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{RBRACKET, "}"},
		{ELSE, "else"},
		{LBRACKET, "{"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{RBRACKET, "}"},

		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{SEMICOLON, ";"},

		{INT, "10"},
		{NOT_EQ, "!="},
		{INT, "9"},
		{SEMICOLON, ";"},

		{EOF, "EOF"},
	}

	for _, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.token || token.Value != test.value {
			t.Fatalf("wanted Token = '%v', Value = '%v'; got Token = '%v', Value = '%v'",
				test.token, test.value, token.Type, token.Value)
		}
	}
}
