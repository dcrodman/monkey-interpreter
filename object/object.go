package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "integer"
	BOOLEAN_OBJ = "boolean"
	NULL_OBJ    = "null"

	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR_OBJ"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType  { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (*Boolean) Type() ObjectType  { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (*Null) Type() ObjectType { return NULL_OBJ }
func (*Null) Inspect() string  { return "null" }

type Return struct {
	Value Object
}

func (*Return) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *Return) Inspect() string {
	return fmt.Sprintf("return %s", r.Value.Inspect())
}

type Error struct {
	Message string
}

func (e Error) Type() ObjectType { return ERROR_OBJ }
func (e Error) Inspect() string  { return "ERROR: " + e.Message }
