package math

import (
	"context"
	"fmt"
	"math"

	"github.com/cloudcmds/tamarin/internal/arg"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

// Name of this module
const Name = "math"

func Abs(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.abs", 1, args); err != nil {
		return err
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
		return object.NewFloat(v)
	default:
		return object.NewError("type error: argument to math.abs not supported, got=%s", args[0].Type())
	}
}

func Sqrt(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.sqrt", 1, args); err != nil {
		return err
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
	if err := arg.Require("math.max", 1, args); err != nil {
		return err
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
			return object.NewFloat(float64(maxInt))
		}
		return object.NewFloat(maxFlt)
	}
	return object.NewFloat(float64(maxInt))
}

func Min(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.min", 1, args); err != nil {
		return err
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
			return object.NewFloat(float64(minInt))
		}
		return object.NewFloat(minFlt)
	}
	return object.NewFloat(float64(minInt))
}

func Ceil(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.ceil", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return object.NewFloat(math.Ceil(arg.Value))
	default:
		return object.NewError("type error: argument to math.ceil not supported, got=%s", args[0].Type())
	}
}

func Floor(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.floor", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		return arg
	case *object.Float:
		return object.NewFloat(math.Floor(arg.Value))
	default:
		return object.NewError("type error: argument to math.floor not supported, got=%s", args[0].Type())
	}
}

func Sin(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.sin", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewFloat(math.Sin(float64(arg.Value)))
	case *object.Float:
		return object.NewFloat(math.Sin(arg.Value))
	default:
		return object.NewError("type error: argument to math.sin not supported, got=%s", args[0].Type())
	}
}

func Cos(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.cos", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Integer:
		return object.NewFloat(math.Cos(float64(arg.Value)))
	case *object.Float:
		return object.NewFloat(math.Cos(arg.Value))
	default:
		return object.NewError("type error: argument to math.cos not supported, got=%s", args[0].Type())
	}
}

func Tan(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.tan", 1, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	return object.NewFloat(math.Tan(x))
}

func Mod(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.mod", 2, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	y, err := object.AsFloat(args[1])
	if err != nil {
		return err
	}
	return object.NewFloat(math.Mod(x, y))
}

func Log(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.log", 1, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	return object.NewFloat(math.Log(x))
}

func Log10(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.log10", 1, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	return object.NewFloat(math.Log10(x))
}

func Log2(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.log2", 1, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	return object.NewFloat(math.Log2(x))
}

func Pow(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.pow", 2, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	y, err := object.AsFloat(args[1])
	if err != nil {
		return err
	}
	return object.NewFloat(math.Pow(x, y))
}

func Pow10(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.pow10", 1, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	return object.NewFloat(math.Pow10(int(x)))
}

func IsInf(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.is_inf", 1, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	return object.NewBoolean(math.IsInf(x, 0))
}

func Round(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.round", 1, args); err != nil {
		return err
	}
	x, err := object.AsFloat(args[0])
	if err != nil {
		return err
	}
	return object.NewFloat(math.Round(x))
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
		{Name: "floor", Func: Floor},
		{Name: "ceil", Func: Ceil},
		{Name: "sin", Func: Sin},
		{Name: "cos", Func: Cos},
		{Name: "tan", Func: Tan},
		{Name: "mod", Func: Mod},
		{Name: "log", Func: Log},
		{Name: "log10", Func: Log10},
		{Name: "log2", Func: Log2},
		{Name: "pow", Func: Pow},
		{Name: "pow10", Func: Pow10},
		{Name: "is_inf", Func: IsInf},
		{Name: "round", Func: Round},
	}); err != nil {
		return nil, err
	}
	s.Declare("PI", &object.Float{Value: math.Pi}, true)
	s.Declare("E", &object.Float{Value: math.E}, true)
	return &object.Module{Name: Name, Scope: s}, nil
}
