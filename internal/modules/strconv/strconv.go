package strconv

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudcmds/tamarin/internal/arg"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
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
		return object.NewOkResult(object.NewInteger(int64(i)))
	}
	return object.NewErrorResult("strconv.atoi: %s", err)
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
		return object.NewOkResult(object.NewBoolean(b))
	}
	return object.NewErrorResult("strconv.parse_bool: %s", err)
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
	return object.NewErrorResult("strconv.parse_float: %s", err)
}

func ParseInt(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strconv.parse_int", 3, args); err != nil {
		return err
	}
	s, typeErr := object.AsString(args[0])
	if typeErr != nil {
		return typeErr
	}
	base, typeErr := object.AsInteger(args[1])
	if typeErr != nil {
		return typeErr
	}
	bitSize, typeErr := object.AsInteger(args[2])
	if typeErr != nil {
		return typeErr
	}
	i, err := strconv.ParseInt(s, int(base), int(bitSize))
	if err == nil {
		return object.NewOkResult(object.NewInteger(i))
	}
	return object.NewErrorResult("strconv.parse_int: %s", err)
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})
	if err := s.AddBuiltins([]scope.Builtin{
		{},
		{Name: "atoi", Func: Atoi},
		{Name: "parse_bool", Func: ParseBool},
		{Name: "parse_float", Func: ParseFloat},
		{Name: "parse_int", Func: ParseInt},
	}); err != nil {
		return nil, err
	}
	return &object.Module{Name: Name, Scope: s}, nil
}
