package ast

import (
	"fmt"
	"monkey-interpreter/token"
	"strings"
)

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s LetStatement) String() string {
	return fmt.Sprintf("%s %s = %s;", s.Token.Value, s.Name, s.Value.String())
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (s ReturnStatement) String() string {
	return fmt.Sprintf("%s %s;", s.Token.Value, s.Value.String())
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (s ExpressionStatement) String() string {
	return s.Expression.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (s BlockStatement) String() string {
	var str strings.Builder

	for _, s := range s.Statements {
		str.WriteString(s.String())
	}

	return str.String()
}
