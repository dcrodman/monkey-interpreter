package lexer

import (
	token2 "monkey-interpreter/token"
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
		token token2.TokenType
		value string
	}{
		{token2.LET, "let"},
		{token2.IDENTIFIER, "five"},
		{token2.ASSIGN, "="},
		{token2.INT, "5"},
		{token2.SEMICOLON, ";"},

		{token2.LET, "let"},
		{token2.IDENTIFIER, "ten"},
		{token2.ASSIGN, "="},
		{token2.INT, "10"},
		{token2.SEMICOLON, ";"},

		{token2.LET, "let"},
		{token2.IDENTIFIER, "add"},
		{token2.ASSIGN, "="},
		{token2.FUNCTION, "fn"},
		{token2.LPAREN, "("},
		{token2.IDENTIFIER, "x"},
		{token2.COMMA, ","},
		{token2.IDENTIFIER, "y"},
		{token2.RPAREN, ")"},
		{token2.LBRACKET, "{"},
		{token2.IDENTIFIER, "x"},
		{token2.PLUS, "+"},
		{token2.IDENTIFIER, "y"},
		{token2.SEMICOLON, ";"},
		{token2.RBRACKET, "}"},
		{token2.SEMICOLON, ";"},

		{token2.LET, "let"},
		{token2.IDENTIFIER, "result"},
		{token2.ASSIGN, "="},
		{token2.IDENTIFIER, "add"},
		{token2.LPAREN, "("},
		{token2.IDENTIFIER, "five"},
		{token2.COMMA, ","},
		{token2.IDENTIFIER, "ten"},
		{token2.RPAREN, ")"},
		{token2.SEMICOLON, ";"},

		{token2.BANG, "!"},
		{token2.MINUS, "-"},
		{token2.SLASH, "/"},
		{token2.ASTERISK, "*"},
		{token2.INT, "5"},
		{token2.SEMICOLON, ";"},

		{token2.INT, "5"},
		{token2.LES, "<"},
		{token2.INT, "10"},
		{token2.GRT, ">"},
		{token2.INT, "5"},
		{token2.SEMICOLON, ";"},

		{token2.IF, "if"},
		{token2.LPAREN, "("},
		{token2.INT, "5"},
		{token2.LES, "<"},
		{token2.INT, "10"},
		{token2.RPAREN, ")"},
		{token2.LBRACKET, "{"},
		{token2.RETURN, "return"},
		{token2.TRUE, "true"},
		{token2.SEMICOLON, ";"},
		{token2.RBRACKET, "}"},
		{token2.ELSE, "else"},
		{token2.LBRACKET, "{"},
		{token2.RETURN, "return"},
		{token2.FALSE, "false"},
		{token2.SEMICOLON, ";"},
		{token2.RBRACKET, "}"},

		{token2.INT, "10"},
		{token2.EQ, "=="},
		{token2.INT, "10"},
		{token2.SEMICOLON, ";"},

		{token2.INT, "10"},
		{token2.NOT_EQ, "!="},
		{token2.INT, "9"},
		{token2.SEMICOLON, ";"},

		{token2.EOF, "EOF"},
	}

	for _, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.token || token.Value != test.value {
			t.Fatalf("wanted Token = '%v', Value = '%v'; got Token = '%v', Value = '%v'",
				test.token, test.value, token.Type, token.Value)
		}
	}
}
