//go:build google
// +build google

package google

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"google.golang.org/api/iterator"
)

func NewMethod(name string, obj interface{}, method *GoMethod) *object.Builtin {
	return object.NewBuiltin(name, func(ctx context.Context, args ...object.Object) object.Object {

		if err := arg.RequireRange(name, 0, 1, args); err != nil {
			return err
		}

		// Build up the input arguments. If a map was supplied, convert it into
		// the params struct type (e.g. *ListBucketsInput).
		inputs := make([]reflect.Value, method.NumIn)
		inputs[0] = reflect.ValueOf(obj)
		inputs[1] = reflect.ValueOf(ctx)
		inst := reflect.New(method.InTypes[2].Elem())
		fmt.Println("ARG", inst.Type())
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

		if len(outputs) == 1 {
			typ := outputs[0].Type()
			fmt.Println("ONE OUTPUT", typ)
			for i := 0; i < typ.NumMethod(); i++ {
				m := typ.Method(i)
				if m.Name == "Next" {
					nextInputs := []reflect.Value{outputs[0]}
					nextResults := m.Func.Call(nextInputs)
					fmt.Println("NEXT RESULTS", nextResults)
					if len(nextResults) == 2 {
						nextVal := nextResults[0]
						nextErr := nextResults[1]
						if !nextErr.IsNil() {
							err := nextErr.Interface().(error)
							if err == iterator.Done {
								fmt.Println("DONE")
							} else {
								fmt.Println("NEXT ERR", err)
								return object.NewError(err)
							}
						}
						fmt.Println("NEXT OK", nextVal.Interface())
					}
				}
				fmt.Println("OUT METHOD", i, m.Name)
			}
		}

		for _, out := range outputs {
			fmt.Println("OUT:", out)
		}
		return object.NewInt(42)
	})
}
