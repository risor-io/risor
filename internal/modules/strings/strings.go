package strings

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudcmds/tamarin/internal/object"
	"github.com/cloudcmds/tamarin/internal/scope"
)

// Name of this module
const Name = "strings"

// Contains determines whether a substring is contained in another string
func Contains(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError("type error: strings.contains() takes exactly two arguments (%d given)", len(args))
	}
	s, ok := args[0].(*object.String)
	if !ok {
		return object.NewError("type error: expected a string (got %v)", args[0].Type())
	}
	substr, ok := args[1].(*object.String)
	if !ok {
		return object.NewError("type error: expected a string (got %v)", args[1].Type())
	}
	if strings.Contains(s.Value, substr.Value) {
		return object.TRUE
	}
	return object.FALSE
}

// Join an array of strings with a given separator
func Join(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewError("type error: strings.join() takes exactly two arguments (%d given)", len(args))
	}
	array, ok := args[0].(*object.Array)
	if !ok {
		return object.NewError("type error: expected an array (got %v)", args[0].Type())
	}
	separator, ok := args[1].(*object.String)
	if !ok {
		return object.NewError("type error: expected a string (got %v)", args[1].Type())
	}
	var stringArray []string
	for _, item := range array.Elements {
		if itemStr, ok := item.(*object.String); ok {
			stringArray = append(stringArray, itemStr.Value)
		} else {
			return object.NewError("type error: array contained a non-string item (type %v)", item.Type())
		}
	}
	return &object.String{Value: strings.Join(stringArray, separator.Value)}
}

// Module returns the `strings` module object
func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})
	if err := s.AddBuiltins([]scope.Builtin{
		{Name: "contains", Func: Contains},
		{Name: "join", Func: Join},
	}); err != nil {
		return nil, err
	}
	return &object.Module{Name: Name, Scope: s}, nil
}
