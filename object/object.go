package object

import (
	"fmt"
	"hash/fnv"
	"monkey-interpreter/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ  = "integer"
	BOOLEAN_OBJ  = "boolean"
	NULL_OBJ     = "null"
	FUNCTION_OBJ = "function"

	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR_OBJ"

	STRING_OBJ = "STRING"
	ARRAY_OBJ  = "ARRAY"
	HASH_OBJ   = "HASH"

	BUILTIN_OBJ = "BUILTIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Null struct{}

func (*Null) Type() ObjectType { return NULL_OBJ }
func (*Null) Inspect() string  { return "null" }

type Integer struct {
	Value int64
}

func (*Integer) Type() ObjectType  { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (*Boolean) Type() ObjectType  { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

func (b *Boolean) HashKey() HashKey {
	var value uint64 = 0
	if b.Value {
		value = 1
	}

	return HashKey{Type: b.Type(), Value: value}
}

type String struct {
	Value string
}

func (*String) Type() ObjectType  { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Array struct {
	Elements []Object
}

func (*Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var str strings.Builder

	elements := make([]string, 0)
	for _, elem := range a.Elements {
		elements = append(elements, elem.Inspect())
	}

	str.WriteString("[")
	str.WriteString(strings.Join(elements, ","))
	str.WriteString("]")

	return str.String()
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var str strings.Builder

	pairs := make([]string, 0)
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	str.WriteString("{")
	str.WriteString(strings.Join(pairs, ", "))
	str.WriteString("}")

	return str.String()
}

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

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var str strings.Builder

	params := make([]string, 0)
	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	str.WriteString("fn (")
	str.WriteString(strings.Join(params, ","))
	str.WriteString(") {\n")
	str.WriteString(f.Body.String())
	str.WriteString("}")

	return str.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
