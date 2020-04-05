package ast

import "strings"

type Node interface {
	String() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type AST struct {
	Statements []Statement
}

func (ast *AST) String() string {
	var program strings.Builder

	for _, s := range ast.Statements {
		program.WriteString(s.String())
	}

	return program.String()
}
