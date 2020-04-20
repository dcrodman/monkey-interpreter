package evaluator

import (
	"monkey-interpreter/ast"
	"monkey-interpreter/object"
)

var (
	NULL = &object.Null{}

	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.AST:
		return evalStatements(node.Statements)
	case ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return evalBoolean(node)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func evalBoolean(b *ast.Boolean) *object.Boolean {
	if b.Value {
		return TRUE
	}
	return FALSE
}
