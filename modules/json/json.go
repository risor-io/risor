package json

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
	"github.com/wI2L/jsondiff"
)

// Name of this module
const Name = "json"

func Unmarshal(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("json.unmarshal", 1, args); err != nil {
		return err
	}
	s, ok := args[0].(*object.String)
	if !ok {
		return object.NewError("type error: expected a string (got %v)", args[0].Type())
	}
	var obj interface{}
	if err := json.Unmarshal([]byte(s.Value), &obj); err != nil {
		return object.NewErrorResult("value error: json.unmarshal failed with: %s", err.Error())
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return object.NewErrorResult("type error: json.unmarshal failed")
	}
	return object.NewOkResult(scriptObj)
}

func Marshal(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("json.marshal", 1, args); err != nil {
		return err
	}
	obj := object.ToGoType(args[0])
	if err, ok := obj.(error); ok {
		return object.NewErrorResult(err.Error())
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return object.NewErrorResult("value error: json.marshal failed with: %s", err.Error())
	}
	return object.NewOkResult(object.NewString(string(b)))
}

func Valid(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("json.valid", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewBool(json.Valid([]byte(s)))
}

func Diff(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("json.diff", 2, args); err != nil {
		return err
	}
	a := object.ToGoType(args[0])
	if err, ok := a.(error); ok {
		return object.NewErrorResult(err.Error())
	}
	b := object.ToGoType(args[1])
	if err, ok := b.(error); ok {
		return object.NewErrorResult(err.Error())
	}
	aBytes, err := json.Marshal(a)
	if err != nil {
		return object.NewErrorResult(err.Error())
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return object.NewErrorResult(err.Error())
	}
	patch, err := jsondiff.CompareJSON(aBytes, bBytes)
	if err != nil {
		return object.NewErrorResult(err.Error())
	}
	patchJSON, err := json.Marshal(patch)
	if err != nil {
		return object.NewErrorResult(err.Error())
	}
	unmarshalArgs := []object.Object{object.NewString(string(patchJSON))}
	return Unmarshal(ctx, unmarshalArgs...)
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := &object.Module{Name: Name, Scope: s}

	if err := s.AddBuiltins([]*object.Builtin{
		{Module: m, Name: "unmarshal", Fn: Unmarshal},
		{Module: m, Name: "marshal", Fn: Marshal},
		{Module: m, Name: "valid", Fn: Valid},
		{Module: m, Name: "diff", Fn: Diff},
	}); err != nil {
		return nil, err
	}
	return m, nil
}
