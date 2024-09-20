package strings

import (
	"context"

	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
)

func asBytes(obj object.Object) (*object.ByteSlice, *object.Error) {
	b, ok := obj.(*object.ByteSlice)
	if !ok {
		return nil, object.Errorf("type error: expected a byte_slice (%s given)", obj.Type())
	}
	return b, nil
}

func Clone(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.clone", 1, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.Clone()
}

func Equals(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.equals", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.Equals(args[1])
}

func Contains(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.contains", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.Contains(args[1])
}

func ContainsAny(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.contains_any", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.ContainsAny(args[1])
}

func ContainsRune(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.contains_rune", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.ContainsRune(args[1])
}

func Count(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.count", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.Count(args[1])
}

func HasPrefix(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.has_prefix", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.HasPrefix(args[1])
}

func HasSuffix(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.hasSuffix", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.HasSuffix(args[1])
}

func Index(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.index", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.Index(args[1])
}

func IndexAny(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.index_any", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.IndexAny(args[1])
}

func IndexByte(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.index_byte", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.IndexByte(args[1])
}

func IndexRune(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.index_rune", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.IndexRune(args[1])
}

func Repeat(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.repeat", 2, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.Repeat(args[1])
}

func Replace(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.replace", 4, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.Replace(args[1], args[2], args[3])
}

func ReplaceAll(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bytes.replace_all", 3, args); err != nil {
		return err
	}
	b, err := asBytes(args[0])
	if err != nil {
		return err
	}
	return b.ReplaceAll(args[1], args[2])
}

func Module() *object.Module {
	return object.NewBuiltinsModule("bytes", map[string]object.Object{
		"clone":         object.NewBuiltin("clone", Clone),
		"contains_any":  object.NewBuiltin("contains_any", ContainsAny),
		"contains_rune": object.NewBuiltin("contains_rune", ContainsRune),
		"contains":      object.NewBuiltin("contains", Contains),
		"count":         object.NewBuiltin("count", Count),
		"equals":        object.NewBuiltin("equals", Equals),
		"has_prefix":    object.NewBuiltin("has_prefix", HasPrefix),
		"has_suffix":    object.NewBuiltin("has_suffix", HasSuffix),
		"index_any":     object.NewBuiltin("index_any", IndexAny),
		"index_byte":    object.NewBuiltin("index_byte", IndexByte),
		"index_rune":    object.NewBuiltin("index_rune", IndexRune),
		"index":         object.NewBuiltin("index", Index),
		"repeat":        object.NewBuiltin("repeat", Repeat),
		"replace_all":   object.NewBuiltin("replace_all", ReplaceAll),
		"replace":       object.NewBuiltin("replace", Replace),
	})
}
