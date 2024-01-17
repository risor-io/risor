package json

import (
	"context"
	"encoding/json"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func Unmarshal(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("json.unmarshal", 1, args); err != nil {
		return err
	}
	data, err := object.AsBytes(args[0])
	if err != nil {
		return err
	}
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return object.Errorf("value error: json.unmarshal failed with: %s", err.Error())
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return object.Errorf("type error: json.unmarshal failed")
	}
	return scriptObj
}

func Marshal(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("json.marshal", 1, 2, args); err != nil {
		return err
	}
	if len(args) == 2 {
		indent, objErr := object.AsString(args[1])
		if objErr != nil {
			return objErr
		}
		b, err := json.MarshalIndent(args[0], "", indent)
		if err != nil {
			return object.Errorf("value error: json.marshal failed: %s", object.NewError(err))
		}
		return object.NewString(string(b))
	}
	b, err := json.Marshal(args[0])
	if err != nil {
		return object.Errorf("value error: json.marshal failed: %s", object.NewError(err))
	}
	return object.NewString(string(b))
}

func Valid(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("json.valid", 1, args); err != nil {
		return err
	}
	data, err := object.AsBytes(args[0])
	if err != nil {
		return err
	}
	return object.NewBool(json.Valid(data))
}

func Module() *object.Module {
	return object.NewBuiltinsModule("json", map[string]object.Object{
		"unmarshal": object.NewBuiltin("unmarshal", Unmarshal),
		"marshal":   object.NewBuiltin("marshal", Marshal),
		"valid":     object.NewBuiltin("valid", Valid),
	})
}
