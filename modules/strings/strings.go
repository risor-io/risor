package strings

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

// Name of this module
const Name = "strings"

func Contains(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.contains", 2, args); err != nil {
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
	if err := arg.Require("strings.has_prefix", 2, args); err != nil {
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
	if err := arg.Require("strings.has_suffix", 2, args); err != nil {
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
	if err := arg.Require("strings.count", 2, args); err != nil {
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
	if err := arg.Require("strings.compare", 2, args); err != nil {
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
	if err := arg.Require("strings.join", 2, args); err != nil {
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
	if err := arg.Require("strings.split", 2, args); err != nil {
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
	return object.NewStringArray(strings.Split(s, sep))
}

func Fields(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.fields", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewStringArray(strings.Fields(s))
}

func Index(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.index", 2, args); err != nil {
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
	if err := arg.Require("strings.last_index", 2, args); err != nil {
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
	if err := arg.Require("strings.replace", 3, args); err != nil {
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
	if err := arg.Require("strings.to_lower", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(strings.ToLower(s))
}

func ToUpper(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.to_upper", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(strings.ToUpper(s))
}

func Trim(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.trim", 2, args); err != nil {
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
	if err := arg.Require("strings.trim_prefix", 2, args); err != nil {
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
	if err := arg.Require("strings.trim_suffix", 2, args); err != nil {
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
	if err := arg.Require("strings.trim_space", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	return object.NewString(strings.TrimSpace(s))
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})

	m := &object.Module{Name: Name, Scope: s}

	if err := s.AddBuiltins([]*object.Builtin{
		{Module: m, Name: "contains", Fn: Contains},
		{Module: m, Name: "count", Fn: Count},
		{Module: m, Name: "has_prefix", Fn: HasPrefix},
		{Module: m, Name: "has_suffix", Fn: HasSuffix},
		{Module: m, Name: "compare", Fn: Compare},
		{Module: m, Name: "join", Fn: Join},
		{Module: m, Name: "split", Fn: Split},
		{Module: m, Name: "fields", Fn: Fields},
		{Module: m, Name: "index", Fn: Index},
		{Module: m, Name: "last_index", Fn: LastIndex},
		{Module: m, Name: "replace", Fn: Replace},
		{Module: m, Name: "to_lower", Fn: ToLower},
		{Module: m, Name: "to_upper", Fn: ToUpper},
		{Module: m, Name: "trim", Fn: Trim},
		{Module: m, Name: "trim_prefix", Fn: TrimPrefix},
		{Module: m, Name: "trim_suffix", Fn: TrimSuffix},
		{Module: m, Name: "trim_space", Fn: TrimSpace},
	}); err != nil {
		return nil, err
	}
	return m, nil
}
