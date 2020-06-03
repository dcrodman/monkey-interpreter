package evaluator

import (
	"fmt"
	"monkey-interpreter/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			var size int

			switch arg := args[0].(type) {
			case *object.String:
				size = len(arg.Value)
			case *object.Array:
				size = len(arg.Elements)
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}

			return &object.Integer{Value: int64(size)}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"tail": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				tailArray := make([]object.Object, len(arg.Elements)-1)

				if len(arg.Elements) > 1 {
					copy(tailArray, arg.Elements[1:])
				}

				return &object.Array{Elements: tailArray}
			default:
				return newError("argument to `tail` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("wrong number of arguments. got=%d, wanted at least 2", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				newArray := make([]object.Object, len(arg.Elements)-1)
				if len(arg.Elements) > 1 {
					copy(newArray, arg.Elements[1:])
				}

				for _, element := range args[1:] {
					newArray = append(newArray, element)
				}

				return &object.Array{Elements: newArray}
			default:
				return newError("argument to `push` not supported, got %s", args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
}
