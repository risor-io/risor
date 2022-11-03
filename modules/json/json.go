package json

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/myzie/tamarin/object"
	"github.com/myzie/tamarin/scope"
)

// Name of this module
const Name = "json"

// Unmarshal a JSON string, creating the corresponding objects
func Unmarshal(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewError("type error: json.unmarshal() takes exactly one argument (%d given)", len(args))
	}
	s, ok := args[0].(*object.String)
	if !ok {
		return object.NewError("type error: expected a string (got %v)", args[0].Type())
	}
	var obj interface{}
	if err := json.Unmarshal([]byte(s.Value), &obj); err != nil {
		return &object.Result{Err: &object.Error{Message: err.Error()}}
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return object.NewError("type error: json.unmarshal failed")
	}
	return &object.Result{Ok: scriptObj}
}

// Module returns the `json` module object
func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})
	if err := s.AddBuiltins([]scope.Builtin{
		{Name: "unmarshal", Func: Unmarshal},
	}); err != nil {
		return nil, err
	}
	return &object.Module{Name: Name, Scope: s}, nil
}
