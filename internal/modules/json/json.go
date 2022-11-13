package json

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

// Name of this module
const Name = "json"

func Unmarshal(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("json.unmarshal", 1, args); err != nil {
		return err
	}
	s, ok := args[0].(*object.String)
	if !ok {
		return object.NewError("type error: expected a string (got %v)", args[0].Type())
	}
	var obj interface{}
	if err := json.Unmarshal([]byte(s.Value), &obj); err != nil {
		return object.NewErrorResult(err.Error())
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return object.NewErrorResult("type error: json.unmarshal failed")
	}
	return &object.Result{Ok: scriptObj}
}

func Marshal(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("json.marshal", 1, args); err != nil {
		return err
	}
	obj := object.ToGoType(args[0])
	if err, ok := obj.(error); ok {
		return object.NewErrorResult(err.Error())
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return object.NewErrorResult(err.Error())
	}
	return &object.Result{Ok: &object.String{Value: string(b)}}
}

func Valid(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("json.valid", 1, args); err != nil {
		return err
	}
	s, err := AsString(args[0])
	if err != nil {
		return err
	}
	return ToBoolean(json.Valid([]byte(s)))
}

func AsString(obj object.Object) (result string, err *object.Error) {
	s, ok := obj.(*object.String)
	if !ok {
		return "", object.NewError("type error: expected a string (got %v)", obj.Type())
	}
	return s.Value, nil
}

func ToBoolean(b bool) *object.Boolean {
	if b {
		return object.TRUE
	}
	return object.FALSE
}

func RequireArgs(funcName string, count int, args []object.Object) *object.Error {
	nArgs := len(args)
	if nArgs != count {
		return object.NewError(
			fmt.Sprintf("type error: %s() takes exactly one argument (%d given)", funcName, nArgs))
	}
	return nil
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})
	if err := s.AddBuiltins([]scope.Builtin{
		{Name: "unmarshal", Func: Unmarshal},
		{Name: "marshal", Func: Marshal},
		{Name: "valid", Func: Valid},
	}); err != nil {
		return nil, err
	}
	return &object.Module{Name: Name, Scope: s}, nil
}
