package ast

import (
	"monkey-interpreter/token"
	"testing"
)

func TestASTString(t *testing.T) {
	ast := AST{
		Statements: []Statement{
			LetStatement{
				Token: token.Token{
					Type:  token.LET,
					Value: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:  token.IDENTIFIER,
						Value: "var",
					},
					Value: "var",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:  token.IDENTIFIER,
						Value: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}

	expected := "let var = anotherVar;"
	if ast.String() != expected {
		t.Errorf("AST.String() incorrect. expected %s, got %s", expected, ast.String())
	}
}
