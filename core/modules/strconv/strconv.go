package strconv

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudcmds/tamarin/core/arg"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/scope"
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
		return object.NewOkResult(object.NewInt(int64(i)))
	}
	return object.NewErrResult(object.NewError(err))
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
		return object.NewOkResult(object.NewBool(b))
	}
	return object.NewErrResult(object.NewError(err))
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
		return object.NewOkResult(object.NewFloat(f))
	}
	return object.NewErrResult(object.NewError(err))
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
		return object.NewOkResult(object.NewInt(i))
	}
	return object.NewErrResult(object.NewError(err))
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := object.NewModule(Name, s)

	if err := s.AddBuiltins([]*object.Builtin{
		object.NewBuiltin("atoi", Atoi, m),
		object.NewBuiltin("parse_bool", ParseBool, m),
		object.NewBuiltin("parse_float", ParseFloat, m),
		object.NewBuiltin("parse_int", ParseInt, m),
	}); err != nil {
		return nil, err
	}
	return m, nil
}
