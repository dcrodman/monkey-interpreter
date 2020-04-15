package ast

import (
	"fmt"
	"monkey-interpreter/token"
	"strings"
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

type Boolean struct {
	Token token.Token
	Value bool
}

func (e Boolean) String() string { return e.Token.Value }

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

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (e IfExpression) String() string {
	var str strings.Builder

	str.WriteString("if ")
	str.WriteString(e.Condition.String())
	str.WriteString(" ")
	str.WriteString(e.Consequence.String())

	if e.Alternative != nil {
		str.WriteString("else")
		str.WriteString(e.Alternative.String())
	}

	return str.String()
}
