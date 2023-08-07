package math

import (
	"context"
	"math"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

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
			return object.Errorf("type error: invalid array item for math.max: %s", val.Type())
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

func Inf(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("math.inf", 0, 1, args); err != nil {
		return err
	}
	sign := 1
	if len(args) == 1 {
		arg, err := object.AsInt(args[0])
		if err != nil {
			return err
		}
		sign = int(arg)
	}
	return object.NewFloat(math.Inf(sign))
}

func Module() *object.Module {
	return object.NewBuiltinsModule("math", map[string]object.Object{
		"abs":    object.NewBuiltin("abs", Abs),
		"ceil":   object.NewBuiltin("ceil", Ceil),
		"cos":    object.NewBuiltin("cos", Cos),
		"E":      object.NewFloat(math.E),
		"floor":  object.NewBuiltin("floor", Floor),
		"inf":    object.NewBuiltin("inf", Inf),
		"is_inf": object.NewBuiltin("is_inf", IsInf),
		"log":    object.NewBuiltin("log", Log),
		"log10":  object.NewBuiltin("log10", Log10),
		"log2":   object.NewBuiltin("log2", Log2),
		"max":    object.NewBuiltin("max", Max),
		"min":    object.NewBuiltin("min", Min),
		"mod":    object.NewBuiltin("mod", Mod),
		"PI":     object.NewFloat(math.Pi),
		"pow":    object.NewBuiltin("pow", Pow),
		"pow10":  object.NewBuiltin("pow10", Pow10),
		"round":  object.NewBuiltin("round", Round),
		"sin":    object.NewBuiltin("sin", Sin),
		"sqrt":   object.NewBuiltin("sqrt", Sqrt),
		"sum":    object.NewBuiltin("sum", Sum),
		"tan":    object.NewBuiltin("tan", Tan),
	})
}
