package json

import (
	"github.com/risor-io/risor/object"
)

const Name = "jsondiff"

// func Diff(ctx context.Context, args ...object.Object) object.Object {
// 	if err := arg.Require("json.diff", 2, args); err != nil {
// 		return err
// 	}
// 	a := args[0].Interface()
// 	if err, ok := a.(error); ok {
// 		return object.NewError(err)
// 	}
// 	b := args[1].Interface()
// 	if err, ok := b.(error); ok {
// 		return object.NewError(err)
// 	}
// 	aBytes, err := json.Marshal(a)
// 	if err != nil {
// 		return object.NewError(err)
// 	}
// 	bBytes, err := json.Marshal(b)
// 	if err != nil {
// 		return object.NewError(err)
// 	}
// 	patch, err := jsondiff.CompareJSON(aBytes, bBytes)
// 	if err != nil {
// 		return object.NewError(err)
// 	}
// 	patchJSON, err := json.Marshal(patch)
// 	if err != nil {
// 		return object.NewError(err)
// 	}
// 	unmarshalArgs := []object.Object{object.NewString(string(patchJSON))}
// 	return Unmarshal(ctx, unmarshalArgs...)
// }

func Module() *object.Module {
	m := object.NewBuiltinsModule(Name, map[string]object.Object{
		// "diff": object.NewBuiltin("diff", Diff),
	})
	return m
}
