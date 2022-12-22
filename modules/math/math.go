package math

import (
	"context"
	"fmt"
	"math"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

// Name of this module
const Name = "math"

func Abs(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.abs", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Int:
		v := arg.Value()
		if v < 0 {
			v = v * -1
		}
		return object.NewInt(v)
	case *object.Float:
		v := arg.Value()
		if v < 0 {
			v = v * -1
		}
		return object.NewFloat(v)
	default:
		return object.Errorf("type error: argument to math.abs not supported, got=%s", args[0].Type())
	}
}

func Sqrt(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.sqrt", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Int:
		v := arg.Value()
		return object.NewFloat(math.Sqrt(float64(v)))
	case *object.Float:
		v := arg.Value()
		return object.NewFloat(math.Sqrt(v))
	default:
		return object.Errorf("type error: argument to math.sqrt not supported, got=%s", args[0].Type())
	}
}

func Max(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.max", 1, args); err != nil {
		return err
	}
	arg := args[0]
	var array []object.Object
	switch arg := arg.(type) {
	case *object.List:
		array = arg.Value()
	case *object.Set:
		array = arg.List().Value()
	default:
		return object.Errorf("type error: %s object is not iterable", args[0].Type())
	}
	if len(array) == 0 {
		return object.Errorf("value error: math.max argument is an empty sequence")
	}
	var maxFlt float64
	var maxInt int64
	var hasFlt, hasInt bool
	for _, value := range array {
		switch val := value.(type) {
		case *object.Int:
			v := val.Value()
			if !hasInt || maxInt < v {
				maxInt = v
				hasInt = true
			}
		case *object.Float:
			v := val.Value()
			if !hasFlt || maxFlt < v {
				maxFlt = v
				hasFlt = true
			}
		default:
			return object.Errorf("invalid array item for math.max: %s", val.Type())
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
	case *object.List:
		array = arg.Value()
	case *object.Set:
		array = arg.List().Value()
	default:
		return object.Errorf("type error: %s object is not iterable", args[0].Type())
	}
	if len(array) == 0 {
		return object.Errorf("value error: math.min argument is an empty sequence")
	}
	var minFlt float64
	var minInt int64
	var hasFlt, hasInt bool
	for _, value := range array {
		switch val := value.(type) {
		case *object.Int:
			v := val.Value()
			if !hasInt || minInt > v {
				minInt = v
				hasInt = true
			}
		case *object.Float:
			v := val.Value()
			if !hasFlt || minFlt > v {
				minFlt = v
				hasFlt = true
			}
		default:
			return object.Errorf("type error: invalid array item for math.min: %s", val.Type())
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

func Sum(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.sum", 1, args); err != nil {
		return err
	}
	arg := args[0]
	var array []object.Object
	switch arg := arg.(type) {
	case *object.List:
		array = arg.Value()
	case *object.Set:
		array = arg.List().Value()
	default:
		return object.Errorf("type error: %s object is not iterable", arg.Type())
	}
	if len(array) == 0 {
		return object.NewFloat(0)
	}
	var sum float64
	for _, value := range array {
		switch val := value.(type) {
		case *object.Int:
			sum += float64(val.Value())
		case *object.Float:
			sum += val.Value()
		default:
			return object.Errorf("value error: invalid input for math.sum: %s", val.Type())
		}
	}
	return object.NewFloat(sum)
}

func Ceil(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.ceil", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Int:
		return arg
	case *object.Float:
		return object.NewFloat(math.Ceil(arg.Value()))
	default:
		return object.Errorf("type error: argument to math.ceil not supported, got=%s", args[0].Type())
	}
}

func Floor(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.floor", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Int:
		return arg
	case *object.Float:
		return object.NewFloat(math.Floor(arg.Value()))
	default:
		return object.Errorf("type error: argument to math.floor not supported, got=%s", args[0].Type())
	}
}

func Sin(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.sin", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Int:
		return object.NewFloat(math.Sin(float64(arg.Value())))
	case *object.Float:
		return object.NewFloat(math.Sin(arg.Value()))
	default:
		return object.Errorf("type error: argument to math.sin not supported, got=%s", args[0].Type())
	}
}

func Cos(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("math.cos", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Int:
		return object.NewFloat(math.Cos(float64(arg.Value())))
	case *object.Float:
		return object.NewFloat(math.Cos(arg.Value()))
	default:
		return object.Errorf("type error: argument to math.cos not supported, got=%s", args[0].Type())
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
	return object.NewBool(math.IsInf(x, 0))
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

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := object.NewModule(Name, s)

	if err := s.AddBuiltins([]*object.Builtin{
		object.NewBuiltin("abs", Abs, m),
		object.NewBuiltin("sqrt", Sqrt, m),
		object.NewBuiltin("min", Min, m),
		object.NewBuiltin("max", Max, m),
		object.NewBuiltin("floor", Floor, m),
		object.NewBuiltin("ceil", Ceil, m),
		object.NewBuiltin("sin", Sin, m),
		object.NewBuiltin("cos", Cos, m),
		object.NewBuiltin("tan", Tan, m),
		object.NewBuiltin("mod", Mod, m),
		object.NewBuiltin("log", Log, m),
		object.NewBuiltin("log10", Log10, m),
		object.NewBuiltin("log2", Log2, m),
		object.NewBuiltin("pow", Pow, m),
		object.NewBuiltin("pow10", Pow10, m),
		object.NewBuiltin("is_inf", IsInf, m),
		object.NewBuiltin("round", Round, m),
		object.NewBuiltin("sum", Sum, m),
	}); err != nil {
		return nil, err
	}
	s.Declare("PI", object.NewFloat(math.Pi), true)
	s.Declare("E", object.NewFloat(math.E), true)
	return m, nil
}
