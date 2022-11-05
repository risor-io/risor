package math

import (
	"context"
	"fmt"
	"math"

	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/internal/scope"
)

// Name of this module
const Name = "math"

func Abs(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("type error: math.abs() takes exactly one argument (%d given)", len(args))
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		v := arg.Value
		if v < 0 {
			v = v * -1
		}
		return &object.Integer{Value: v}
	case *object.Float:
		v := arg.Value
		if v < 0 {
			v = v * -1
		}
		return &object.Float{Value: v}
	default:
		return object.NewError("type error: argument to math.abs not supported, got=%s", args[0].Type())
	}
}

func Sqrt(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("type error: math.sqrt() takes exactly one argument (%d given)", len(args))
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		v := arg.Value
		return &object.Float{Value: math.Sqrt(float64(v))}
	case *object.Float:
		v := arg.Value
		return &object.Float{Value: math.Sqrt(v)}
	default:
		return object.NewError("type error: argument to math.sqrt not supported, got=%s", args[0].Type())
	}
}

func Max(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("type error: math.max() takes exactly one argument (%d given)", len(args))
	}
	arg := args[0]
	var array []object.Object
	switch arg := arg.(type) {
	case *object.Array:
		array = arg.Elements
	case *object.Set:
		array = arg.Array()
	default:
		return object.NewError("type error: %s object is not iterable", args[0].Type())
	}
	if len(array) == 0 {
		return object.NewError("value error: math.max argument is an empty sequence")
	}
	var maxFlt float64
	var maxInt int64
	var hasFlt, hasInt bool
	for _, value := range array {
		switch val := value.(type) {
		case *object.Integer:
			if !hasInt || maxInt < val.Value {
				maxInt = val.Value
				hasInt = true
			}
		case *object.Float:
			if !hasFlt || maxFlt < val.Value {
				maxFlt = val.Value
				hasFlt = true
			}
		default:
			return object.NewError("invalid array item for math.max: %s", val.Type())
		}
	}
	if hasFlt {
		if hasInt && float64(maxInt) > maxFlt {
			return &object.Float{Value: float64(maxInt)}
		}
		return &object.Float{Value: maxFlt}
	}
	return &object.Integer{Value: maxInt}
}

func Min(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("type error: math.min() takes exactly one argument (%d given)", len(args))
	}
	arg := args[0]
	var array []object.Object
	switch arg := arg.(type) {
	case *object.Array:
		array = arg.Elements
	case *object.Set:
		array = arg.Array()
	default:
		return object.NewError("type error: %s object is not iterable", args[0].Type())
	}
	if len(array) == 0 {
		return object.NewError("value error: math.min argument is an empty sequence")
	}
	var minFlt float64
	var minInt int64
	var hasFlt, hasInt bool
	for _, value := range array {
		switch val := value.(type) {
		case *object.Integer:
			if !hasInt || minInt > val.Value {
				minInt = val.Value
				hasInt = true
			}
		case *object.Float:
			if !hasFlt || minFlt > val.Value {
				minFlt = val.Value
				hasFlt = true
			}
		default:
			return object.NewError("type error: invalid array item for math.min: %s", val.Type())
		}
	}
	if hasFlt {
		if hasInt && float64(minInt) < minFlt {
			return &object.Float{Value: float64(minInt)}
		}
		return &object.Float{Value: minFlt}
	}
	return &object.Integer{Value: minInt}
}

// Module returns the `math` module object
func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})
	if err := s.AddBuiltins([]scope.Builtin{
		{Name: "abs", Func: Abs},
		{Name: "sqrt", Func: Sqrt},
		{Name: "min", Func: Min},
		{Name: "max", Func: Max},
	}); err != nil {
		return nil, err
	}
	s.Declare("PI", &object.Float{Value: math.Pi}, true)
	s.Declare("E", &object.Float{Value: math.E}, true)
	return &object.Module{Name: Name, Scope: s}, nil
}
