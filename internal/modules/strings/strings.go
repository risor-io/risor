package strings

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

// Name of this module
const Name = "strings"

func Contains(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.contains", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	substr, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewBoolean(strings.Contains(s, substr))
}

func HasPrefix(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.has_prefix", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	prefix, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewBoolean(strings.HasPrefix(s, prefix))
}

func HasSuffix(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.has_suffix", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	suffix, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewBoolean(strings.HasSuffix(s, suffix))
}

func Count(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.count", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	substr, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return &object.Integer{Value: int64(strings.Count(s, substr))}
}

func Compare(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.compare", 2, args); err != nil {
		return err
	}
	s1, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	s2, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return &object.Integer{Value: int64(strings.Compare(s1, s2))}
}

func Join(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.join", 2, args); err != nil {
		return err
	}
	array, err := object.AsArray(args[0])
	if err != nil {
		return err
	}
	separator, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	var stringArray []string
	for _, item := range array.Elements {
		itemStr, err := object.AsString(item)
		if err != nil {
			return err
		}
		stringArray = append(stringArray, itemStr)
	}
	return object.NewString(strings.Join(stringArray, separator))
}

func Split(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.split", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	sep, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return NewStringArray(strings.Split(s, sep))
}

func Fields(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.fields", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return NewStringArray(strings.Fields(s))
}

func Index(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.index", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	substr, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return &object.Integer{Value: int64(strings.Index(s, substr))}
}

func LastIndex(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.last_index", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	substr, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return &object.Integer{Value: int64(strings.LastIndex(s, substr))}
}

func Replace(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.replace", 3, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	old, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	new, err := object.AsString(args[2])
	if err != nil {
		return err
	}
	return &object.String{Value: strings.ReplaceAll(s, old, new)}
}

func ToLower(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.to_lower", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(strings.ToLower(s))
}

func ToUpper(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.to_upper", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(strings.ToUpper(s))
}

func Trim(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.trim", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	cutset, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewString(strings.Trim(s, cutset))
}

func TrimPrefix(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.trim_prefix", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	prefix, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewString(strings.TrimPrefix(s, prefix))
}

func TrimSuffix(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.trim_suffix", 2, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	suffix, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewString(strings.TrimSuffix(s, suffix))
}

func TrimSpace(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("strings.trim_space", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(strings.TrimSpace(s))
}

func RequireArgs(funcName string, count int, args []object.Object) *object.Error {
	nArgs := len(args)
	if nArgs != count {
		return object.NewError(
			fmt.Sprintf("type error: %s() takes exactly one argument (%d given)", funcName, nArgs))
	}
	return nil
}

func NewStringArray(s []string) *object.Array {
	array := &object.Array{}
	for _, item := range s {
		array.Elements = append(array.Elements, &object.String{Value: item})
	}
	return array
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})
	if err := s.AddBuiltins([]scope.Builtin{
		{Name: "contains", Func: Contains},
		{Name: "count", Func: Count},
		{Name: "has_prefix", Func: HasPrefix},
		{Name: "has_suffix", Func: HasSuffix},
		{Name: "compare", Func: Compare},
		{Name: "join", Func: Join},
		{Name: "split", Func: Split},
		{Name: "fields", Func: Fields},
		{Name: "index", Func: Index},
		{Name: "last_index", Func: LastIndex},
		{Name: "replace", Func: Replace},
		{Name: "to_lower", Func: ToLower},
		{Name: "to_upper", Func: ToUpper},
		{Name: "trim", Func: Trim},
		{Name: "trim_prefix", Func: TrimPrefix},
		{Name: "trim_suffix", Func: TrimSuffix},
		{Name: "trim_space", Func: TrimSpace},
	}); err != nil {
		return nil, err
	}
	return &object.Module{Name: Name, Scope: s}, nil
}
