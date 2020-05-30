package lexer

import (
	"monkey-interpreter/token"
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
	"foobar"
	"foo bar"
	`
	lexer := New(code)

	tests := []struct {
		token token.TokenType
		value string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACKET, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACKET, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LES, "<"},
		{token.INT, "10"},
		{token.GRT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LES, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACKET, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACKET, "}"},
		{token.ELSE, "else"},
		{token.LBRACKET, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACKET, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},

		{token.EOF, "EOF"},
	}

	for _, test := range tests {
		nextToken := lexer.NextToken()

		if nextToken.Type != test.token || nextToken.Value != test.value {
			t.Fatalf("wanted Token = '%v', Value = '%v'; got Token = '%v', Value = '%v'",
				test.token, test.value, nextToken.Type, nextToken.Value)
		}
	}
}
