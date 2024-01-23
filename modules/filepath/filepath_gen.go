// Code generated by risor-modgen. DO NOT EDIT.

package filepath

import (
	"context"
	"github.com/risor-io/risor/object"
)

// Abs is a wrapper function around [abs]
// that implements [object.BuiltinFunction].
func Abs(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("filepath.abs", 1, len(args))
	}
	pathParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result, resultErr := abs(ctx, pathParam)
	if resultErr != nil {
		return object.NewError(resultErr)
	}
	return object.NewString(result)
}

// Base is a wrapper function around [base]
// that implements [object.BuiltinFunction].
func Base(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("filepath.base", 1, len(args))
	}
	pathParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := base(pathParam)
	return object.NewString(result)
}

// Clean is a wrapper function around [clean]
// that implements [object.BuiltinFunction].
func Clean(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("filepath.clean", 1, len(args))
	}
	pathParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := clean(pathParam)
	return object.NewString(result)
}

// Dir is a wrapper function around [dir]
// that implements [object.BuiltinFunction].
func Dir(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("filepath.dir", 1, len(args))
	}
	pathParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := dir(pathParam)
	return object.NewString(result)
}

// Ext is a wrapper function around [ext]
// that implements [object.BuiltinFunction].
func Ext(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("filepath.ext", 1, len(args))
	}
	pathParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := ext(pathParam)
	return object.NewString(result)
}

// IsAbs is a wrapper function around [isAbs]
// that implements [object.BuiltinFunction].
func IsAbs(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("filepath.is_abs", 1, len(args))
	}
	pathParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := isAbs(pathParam)
	return object.NewBool(result)
}

// Match is a wrapper function around [match]
// that implements [object.BuiltinFunction].
func Match(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("filepath.match", 2, len(args))
	}
	patternParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	nameParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result, resultErr := match(patternParam, nameParam)
	if resultErr != nil {
		return object.NewError(resultErr)
	}
	return object.NewBool(result)
}

// Rel is a wrapper function around [rel]
// that implements [object.BuiltinFunction].
func Rel(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("filepath.rel", 2, len(args))
	}
	basepathParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	targpathParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result, resultErr := rel(basepathParam, targpathParam)
	if resultErr != nil {
		return object.NewError(resultErr)
	}
	return object.NewString(result)
}

// Split is a wrapper function around [split]
// that implements [object.BuiltinFunction].
func Split(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("filepath.split", 1, len(args))
	}
	pathParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := split(pathParam)
	return object.NewStringList(result)
}

// SplitList is a wrapper function around [splitList]
// that implements [object.BuiltinFunction].
func SplitList(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("filepath.split_list", 1, len(args))
	}
	pathListParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := splitList(pathListParam)
	return object.NewStringList(result)
}

// addGeneratedBuiltins adds the generated builtin wrappers to the given map.
//
// Useful if you want to write your own "Module()" function.
func addGeneratedBuiltins(builtins map[string]object.Object) map[string]object.Object {
	builtins["abs"] = object.NewBuiltin("filepath.abs", Abs)
	builtins["base"] = object.NewBuiltin("filepath.base", Base)
	builtins["clean"] = object.NewBuiltin("filepath.clean", Clean)
	builtins["dir"] = object.NewBuiltin("filepath.dir", Dir)
	builtins["ext"] = object.NewBuiltin("filepath.ext", Ext)
	builtins["is_abs"] = object.NewBuiltin("filepath.is_abs", IsAbs)
	builtins["match"] = object.NewBuiltin("filepath.match", Match)
	builtins["rel"] = object.NewBuiltin("filepath.rel", Rel)
	builtins["split"] = object.NewBuiltin("filepath.split", Split)
	builtins["split_list"] = object.NewBuiltin("filepath.split_list", SplitList)
	return builtins
}


