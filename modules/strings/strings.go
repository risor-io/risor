package strings

import (
	"context"
	"strings"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

//risor:generate no-module-func

func asString(obj object.Object) (*object.String, *object.Error) {
	s, ok := obj.(*object.String)
	if !ok {
		return nil, object.Errorf("type error: expected a string (got %v)", obj.Type())
	}
	return s, nil
}

//risor:export
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

//risor:export has_prefix
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

//risor:export has_prefix
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

//risor:export
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

//risor:export
func Compare(a, b string) int {
	return strings.Compare(a, b)
}

func Join(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.join", 2, args); err != nil {
		return err
	}
	list, err := object.AsList(args[0])
	if err != nil {
		return err
	}
	sep, err := asString(args[1])
	if err != nil {
		return err
	}
	return sep.Join(list)
}

func Split(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.split", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.Split(args[1])
}

func Fields(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.fields", 1, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.Fields()
}

//risor:export
func Index(s, substr string) int {
	return strings.Index(s, substr)
}

//risor:export last_index
func LastIndex(s, substr string) int {
	return strings.LastIndex(s, substr)
}

//risor:export replace_all
func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

//risor:export to_lower
func ToLower(s string) string {
	return strings.ToLower(s)
}

//risor:export to_upper
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

//risor:export
func Trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

//risor:export trim_prefix
func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

//risor:export trim_suffix
func TrimSuffix(s, prefix string) string {
	return strings.TrimSuffix(s, prefix)
}

//risor:export trim_space
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("strings", addGeneratedBuiltins(map[string]object.Object{
		"fields": object.NewBuiltin("fields", Fields),
		"join":   object.NewBuiltin("join", Join),
		"split":  object.NewBuiltin("split", Split),
	}))
}
