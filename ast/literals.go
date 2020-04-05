package ast

import (
	"fmt"
	"monkey-interpreter/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (e Identifier) String() string { return e.Value }

type Integer struct {
	Token token.Token
	Value int64
}

func (e Integer) String() string { return e.Token.Value }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (e PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator, e.Right.String())
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (e InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left.String(), e.Operator, e.Right.String())
}
