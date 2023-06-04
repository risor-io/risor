package json

import (
	"context"
	"encoding/json"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/object"
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
		return object.Errorf("type error: expected a string (got %v)", args[0].Type())
	}
	var obj interface{}
	if err := json.Unmarshal([]byte(s.Value()), &obj); err != nil {
		return object.Errorf("value error: json.unmarshal failed with: %s", object.NewError(err))
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return object.Errorf("type error: json.unmarshal failed")
	}
	return scriptObj
}

func Marshal(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("json.marshal", 1, args); err != nil {
		return err
	}
	obj := args[0].Interface()
	if err, ok := obj.(error); ok {
		return object.NewError(err)
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return object.Errorf("value error: json.marshal failed with: %s", object.NewError(err))
	}
	return object.NewString(string(b))
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
	a := args[0].Interface()
	if err, ok := a.(error); ok {
		return object.NewError(err)
	}
	b := args[1].Interface()
	if err, ok := b.(error); ok {
		return object.NewError(err)
	}
	aBytes, err := json.Marshal(a)
	if err != nil {
		return object.NewError(err)
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return object.NewError(err)
	}
	patch, err := jsondiff.CompareJSON(aBytes, bBytes)
	if err != nil {
		return object.NewError(err)
	}
	patchJSON, err := json.Marshal(patch)
	if err != nil {
		return object.NewError(err)
	}
	unmarshalArgs := []object.Object{object.NewString(string(patchJSON))}
	return Unmarshal(ctx, unmarshalArgs...)
}

func Module() *object.Module {
	m := object.NewBuiltinsModule(Name, map[string]object.Object{
		"unmarshal": object.NewBuiltin("unmarshal", Unmarshal),
		"marshal":   object.NewBuiltin("marshal", Marshal),
		"valid":     object.NewBuiltin("valid", Valid),
		"diff":      object.NewBuiltin("diff", Diff),
	})
	return m
}
