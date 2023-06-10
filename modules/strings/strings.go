package strings

import (
	"context"

	"github.com/cloudcmds/tamarin/v2/internal/arg"
	"github.com/cloudcmds/tamarin/v2/object"
)

// Name of this module
const Name = "strings"

func asString(obj object.Object) (*object.String, *object.Error) {
	s, ok := obj.(*object.String)
	if !ok {
		return nil, object.Errorf("type error: expected a string (got %v)", obj.Type())
	}
	return s, nil
}

func Contains(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.contains", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.Contains(args[1])
}

func HasPrefix(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.has_prefix", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.HasPrefix(args[1])
}

func HasSuffix(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.has_suffix", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.HasSuffix(args[1])
}

func Count(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.count", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.Count(args[1])
}

func Compare(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.compare", 2, args); err != nil {
		return err
	}
	s, errObj := asString(args[0])
	if errObj != nil {
		return errObj
	}
	value, err := s.Compare(args[1])
	if err != nil {
		return object.NewError(err)
	}
	return object.NewInt(int64(value))
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

func Index(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.index", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.Index(args[1])
}

func LastIndex(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.last_index", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.LastIndex(args[1])
}

func ReplaceAll(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.replace", 3, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.ReplaceAll(args[1], args[2])
}

func ToLower(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.to_lower", 1, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.ToLower()
}

func ToUpper(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.to_upper", 1, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.ToUpper()
}

func Trim(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.trim", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.Trim(args[1])
}

func TrimPrefix(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.trim_prefix", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.TrimPrefix(args[1])
}

func TrimSuffix(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.trim_suffix", 2, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.TrimSuffix(args[1])
}

func TrimSpace(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("strings.trim_space", 1, args); err != nil {
		return err
	}
	s, err := asString(args[0])
	if err != nil {
		return err
	}
	return s.TrimSpace()
}

func Module() *object.Module {
	m := object.NewBuiltinsModule(Name, map[string]object.Object{
		"contains":    object.NewBuiltin("contains", Contains),
		"count":       object.NewBuiltin("count", Count),
		"has_prefix":  object.NewBuiltin("has_prefix", HasPrefix),
		"has_suffix":  object.NewBuiltin("has_suffix", HasSuffix),
		"compare":     object.NewBuiltin("compare", Compare),
		"join":        object.NewBuiltin("join", Join),
		"split":       object.NewBuiltin("split", Split),
		"fields":      object.NewBuiltin("fields", Fields),
		"index":       object.NewBuiltin("index", Index),
		"last_index":  object.NewBuiltin("last_index", LastIndex),
		"replace_all": object.NewBuiltin("replace_all", ReplaceAll),
		"to_lower":    object.NewBuiltin("to_lower", ToLower),
		"to_upper":    object.NewBuiltin("to_upper", ToUpper),
		"trim":        object.NewBuiltin("trim", Trim),
		"trim_prefix": object.NewBuiltin("trim_prefix", TrimPrefix),
		"trim_suffix": object.NewBuiltin("trim_suffix", TrimSuffix),
		"trim_space":  object.NewBuiltin("trim_space", TrimSpace),
	})
	return m
}
