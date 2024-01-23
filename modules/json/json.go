package json

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

//risor:generate no-module-func

//risor:export
func unmarshal(data []byte) (object.Object, error) {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("value error: json.unmarshal failed with: %w", err)
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return nil, fmt.Errorf("type error: json.unmarshal failed")
	}
	return scriptObj, nil
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

//risor:export
func valid(data []byte) bool {
	return json.Valid(data)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("json", addGeneratedBuiltins(map[string]object.Object{
		"marshal": object.NewBuiltin("marshal", Marshal),
	}))
}
