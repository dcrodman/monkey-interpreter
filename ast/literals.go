package ast

import (
	"fmt"
	"monkey-interpreter/token"
	"strings"
)

type Integer struct {
	Token token.Token
	Value int64
}

func (e *Integer) String() string { return e.Token.Value }

type Boolean struct {
	Token token.Token
	Value bool
}

func (e *Boolean) String() string { return e.Token.Value }

type String struct {
	Token token.Token
	Value string
}

func (e *String) String() string { return e.Token.Value }

type Identifier struct {
	Token token.Token
	Value string
}

func (e *Identifier) String() string { return e.Value }

type Function struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (e *Function) String() string {
	var str strings.Builder

	str.WriteString("func ")
	str.WriteString("(")

	var params []string
	for _, p := range e.Parameters {
		params = append(params, p.String())
	}

	str.WriteString(strings.Join(params, ","))
	str.WriteString(")")
	str.WriteString(e.Body.String())

	return str.String()
}

type Array struct {
	Token    token.Token
	Elements []Expression
}

func (a *Array) String() string {
	var str strings.Builder

	str.WriteString("[")

	var params []string
	for _, e := range a.Elements {
		params = append(params, e.String())
	}

	str.WriteString(strings.Join(params, ", "))
	str.WriteString("]")

	return str.String()
}

type Hash struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *Hash) String() string {
	var str strings.Builder

	pairs := make([]string, 0)
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	str.WriteString("{")
	str.WriteString(strings.Join(pairs, ", "))
	str.WriteString("}")

	return str.String()
}

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

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (e CallExpression) String() string {
	var str strings.Builder

	str.WriteString(e.Function.String())
	str.WriteString("(")

	var params []string
	for _, p := range e.Arguments {
		params = append(params, p.String())
	}

	str.WriteString(strings.Join(params, ", "))
	str.WriteString(")")

	return str.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) String() string {
	var str strings.Builder

	str.WriteString("(")
	str.WriteString(ie.Left.String())
	str.WriteString("[")
	str.WriteString(ie.Index.String())
	str.WriteString("]")
	str.WriteString(")")

	return str.String()
}
