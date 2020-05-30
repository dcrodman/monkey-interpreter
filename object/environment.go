package object

type Environment struct {
	symbols map[string]Object

	outer *Environment
}

func NewEnvironment() *Environment {
	return NewEnclosingEnvironment(nil)
}

func NewEnclosingEnvironment(enclosing *Environment) *Environment {
	return &Environment{symbols: make(map[string]Object), outer: enclosing}
}

func (e *Environment) Get(identifier string) (Object, bool) {
	val, ok := e.symbols[identifier]

	if !ok && e.outer != nil {
		val, ok = e.outer.Get(identifier)
	}
	return val, ok
}

func (e *Environment) Set(identifier string, val Object) Object {
	e.symbols[identifier] = val
	return val
}
