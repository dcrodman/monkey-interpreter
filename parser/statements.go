package parser

import "monkey-interpreter/token"

type Identifier struct {
	Token token.Token
	Value string
}

type Expression interface{}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}
