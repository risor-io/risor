package strconv

import (
	"context"
	"strconv"

	"github.com/cloudcmds/tamarin/v2/internal/arg"
	"github.com/cloudcmds/tamarin/v2/object"
)

// Name of this module
const Name = "strconv"

func Atoi(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strconv.atoi", 1, args); err != nil {
		return err
	}
	s, typeErr := object.AsString(args[0])
	if typeErr != nil {
		return typeErr
	}
	i, err := strconv.Atoi(s)
	if err == nil {
		return object.NewInt(int64(i))
	}
	return object.NewError(err)
}

func ParseBool(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strconv.parse_bool", 1, args); err != nil {
		return err
	}
	s, typeErr := object.AsString(args[0])
	if typeErr != nil {
		return typeErr
	}
	b, err := strconv.ParseBool(s)
	if err == nil {
		return object.NewBool(b)
	}
	return object.NewError(err)
}

func ParseFloat(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strconv.parse_float", 1, args); err != nil {
		return err
	}
	s, typeErr := object.AsString(args[0])
	if typeErr != nil {
		return typeErr
	}
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return object.NewFloat(f)
	}
	return object.NewError(err)
}

func ParseInt(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strconv.parse_int", 3, args); err != nil {
		return err
	}
	s, typeErr := object.AsString(args[0])
	if typeErr != nil {
		return typeErr
	}
	base, typeErr := object.AsInt(args[1])
	if typeErr != nil {
		return typeErr
	}
	bitSize, typeErr := object.AsInt(args[2])
	if typeErr != nil {
		return typeErr
	}
	i, err := strconv.ParseInt(s, int(base), int(bitSize))
	if err == nil {
		return object.NewInt(i)
	}
	return object.NewError(err)
}

func Module() *object.Module {
	m := object.NewBuiltinsModule(Name, map[string]object.Object{
		"atoi":        object.NewBuiltin("atoi", Atoi),
		"parse_bool":  object.NewBuiltin("parse_bool", ParseBool),
		"parse_float": object.NewBuiltin("parse_float", ParseFloat),
		"parse_int":   object.NewBuiltin("parse_int", ParseInt),
	})
	return m
}
