package ast

type Statement interface{}

type AST struct {
	Statements []Statement
}
