//go:build aws
// +build aws

package aws

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func NewMethod(name string, obj interface{}, method *GoMethod) *object.Builtin {
	return object.NewBuiltin(name, func(ctx context.Context, args ...object.Object) object.Object {

		if err := arg.RequireRange(name, 0, 1, args); err != nil {
			return err
		}
		// The AWS SDK v2 has a common signature for all service methods:
		// - client (usually c *Client)
		// - ctx context.Context
		// - params (e.g. *ListBucketsInput)
		// - optFns (e.g. ...func(*Options))
		if method.NumIn != 4 {
			return object.Errorf("unsupported aws method: %s (unexpected parameter count: %d)", name, method.NumIn)
		}

		// Build up the input arguments. If a map was supplied, convert it into
		// the params struct type (e.g. *ListBucketsInput).
		inputs := make([]reflect.Value, method.NumIn)
		inputs[0] = reflect.ValueOf(obj)
		inputs[1] = reflect.ValueOf(ctx)
		inst := reflect.New(method.InTypes[2].Elem())
		if len(args) == 1 {
			m, err := object.AsMap(args[0])
			if err != nil {
				return err
			}
			paramData, jsErr := json.Marshal(m.Interface())
			if jsErr != nil {
				return object.NewError(jsErr)
			}
			if err := json.Unmarshal(paramData, inst.Interface()); err != nil {
				return object.NewError(err)
			}
		}
		inputs[2] = inst

		// Call the method
		outputs := method.Method.Func.Call(inputs[:3])

		// There should be two return values: the output and an error
		if len(outputs) != 2 {
			return object.Errorf("unsupported aws method: %s (unexpected return count: %d)", name, len(outputs))
		}

		// If there is a non-nil error in the output, return a Risor error
		if errOut := outputs[1]; !errOut.IsNil() {
			return object.NewError(errOut.Interface().(error))
		}

		// Convert the output to a Risor map
		outputData, err := json.Marshal(outputs[0].Interface())
		if err != nil {
			return object.NewError(err)
		}
		var resultMap map[string]interface{}
		if err := json.Unmarshal(outputData, &resultMap); err != nil {
			return object.NewError(err)
		}
		return object.FromGoType(resultMap)
	})
}
