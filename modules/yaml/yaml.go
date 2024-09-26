package yaml

import (
	"context"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"gopkg.in/yaml.v3"
)

func Unmarshal(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("yaml.unmarshal", 1, args); err != nil {
		return err
	}
	data, err := object.AsBytes(args[0])
	if err != nil {
		return err
	}
	var obj interface{}
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return object.Errorf("value error: yaml.unmarshal failed with: %s", err.Error())
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return object.TypeErrorf("type error: yaml.unmarshal failed")
	}
	return scriptObj
}

func Marshal(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("yaml.marshal", 1, args); err != nil {
		return err
	}
	b, err := yaml.Marshal(args[0].Interface())
	if err != nil {
		return object.Errorf("value error: yaml.marshal failed: %s", object.NewError(err))
	}
	return object.NewString(string(b))
}

func Valid(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("yaml.valid", 1, args); err != nil {
		return err
	}
	data, err := object.AsBytes(args[0])
	if err != nil {
		return err
	}
	var v any
	return object.NewBool(yaml.Unmarshal(data, &v) == nil)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("yaml", map[string]object.Object{
		"unmarshal": object.NewBuiltin("unmarshal", Unmarshal),
		"marshal":   object.NewBuiltin("marshal", Marshal),
		"valid":     object.NewBuiltin("valid", Valid),
	})
}
