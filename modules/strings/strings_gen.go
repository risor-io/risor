// Code generated by risor-modgen. DO NOT EDIT.

package strings

import (
	"context"
	"github.com/risor-io/risor/object"
	"math"
)

// Contains is a wrapper function around [contains]
// that implements [object.BuiltinFunction].
func Contains(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.contains", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	substrParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := contains(sParam, substrParam)
	return object.NewBool(result)
}

// HasPrefix is a wrapper function around [hasPrefix]
// that implements [object.BuiltinFunction].
func HasPrefix(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.has_prefix", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	prefixParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := hasPrefix(sParam, prefixParam)
	return object.NewBool(result)
}

// HasSuffix is a wrapper function around [hasSuffix]
// that implements [object.BuiltinFunction].
func HasSuffix(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.has_suffix", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	suffixParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := hasSuffix(sParam, suffixParam)
	return object.NewBool(result)
}

// Count is a wrapper function around [count]
// that implements [object.BuiltinFunction].
func Count(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.count", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	substrParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := count(sParam, substrParam)
	return object.NewInt(int64(result))
}

// Compare is a wrapper function around [compare]
// that implements [object.BuiltinFunction].
func Compare(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.compare", 2, len(args))
	}
	aParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	bParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := compare(aParam, bParam)
	return object.NewInt(int64(result))
}

// Repeat is a wrapper function around [repeat]
// that implements [object.BuiltinFunction].
func Repeat(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.repeat", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	countParamRaw, err := object.AsInt(args[1])
	if err != nil {
		return err
	}
	if countParamRaw > math.MaxInt {
		return object.Errorf("type error: strings.repeat argument 'count' (index 1) cannot be > %v", math.MaxInt)
	}
	if countParamRaw < math.MinInt {
		return object.Errorf("type error: strings.repeat argument 'count' (index 1) cannot be < %v", math.MinInt)
	}
	countParam := int(countParamRaw)
	result := repeat(sParam, countParam)
	return object.NewString(result)
}

// Join is a wrapper function around [join]
// that implements [object.BuiltinFunction].
func Join(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.join", 2, len(args))
	}
	listParam, err := object.AsStringSlice(args[0])
	if err != nil {
		return err
	}
	sepParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := join(listParam, sepParam)
	return object.NewString(result)
}

// Split is a wrapper function around [split]
// that implements [object.BuiltinFunction].
func Split(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.split", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	sepParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := split(sParam, sepParam)
	return object.NewStringList(result)
}

// Fields is a wrapper function around [fields]
// that implements [object.BuiltinFunction].
func Fields(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("strings.fields", 1, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := fields(sParam)
	return object.NewStringList(result)
}

// Index is a wrapper function around [index]
// that implements [object.BuiltinFunction].
func Index(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.index", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	substrParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := index(sParam, substrParam)
	return object.NewInt(int64(result))
}

// LastIndex is a wrapper function around [lastIndex]
// that implements [object.BuiltinFunction].
func LastIndex(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.last_index", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	substrParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := lastIndex(sParam, substrParam)
	return object.NewInt(int64(result))
}

// ReplaceAll is a wrapper function around [replaceAll]
// that implements [object.BuiltinFunction].
func ReplaceAll(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.NewArgsError("strings.replace_all", 3, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	oldParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	newParam, err := object.AsString(args[2])
	if err != nil {
		return err
	}
	result := replaceAll(sParam, oldParam, newParam)
	return object.NewString(result)
}

// ToLower is a wrapper function around [toLower]
// that implements [object.BuiltinFunction].
func ToLower(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("strings.to_lower", 1, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := toLower(sParam)
	return object.NewString(result)
}

// ToUpper is a wrapper function around [toUpper]
// that implements [object.BuiltinFunction].
func ToUpper(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("strings.to_upper", 1, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := toUpper(sParam)
	return object.NewString(result)
}

// Trim is a wrapper function around [trim]
// that implements [object.BuiltinFunction].
func Trim(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.trim", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	cutsetParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := trim(sParam, cutsetParam)
	return object.NewString(result)
}

// TrimPrefix is a wrapper function around [trimPrefix]
// that implements [object.BuiltinFunction].
func TrimPrefix(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.trim_prefix", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	prefixParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := trimPrefix(sParam, prefixParam)
	return object.NewString(result)
}

// TrimSuffix is a wrapper function around [trimSuffix]
// that implements [object.BuiltinFunction].
func TrimSuffix(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.NewArgsError("strings.trim_suffix", 2, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	prefixParam, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	result := trimSuffix(sParam, prefixParam)
	return object.NewString(result)
}

// TrimSpace is a wrapper function around [trimSpace]
// that implements [object.BuiltinFunction].
func TrimSpace(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewArgsError("strings.trim_space", 1, len(args))
	}
	sParam, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	result := trimSpace(sParam)
	return object.NewString(result)
}

// addGeneratedBuiltins adds the generated builtin wrappers to the given map.
//
// Useful if you want to write your own "Module()" function.
func addGeneratedBuiltins(builtins map[string]object.Object) map[string]object.Object {
	builtins["contains"] = object.NewBuiltin("strings.contains", Contains)
	builtins["has_prefix"] = object.NewBuiltin("strings.has_prefix", HasPrefix)
	builtins["has_suffix"] = object.NewBuiltin("strings.has_suffix", HasSuffix)
	builtins["count"] = object.NewBuiltin("strings.count", Count)
	builtins["compare"] = object.NewBuiltin("strings.compare", Compare)
	builtins["repeat"] = object.NewBuiltin("strings.repeat", Repeat)
	builtins["join"] = object.NewBuiltin("strings.join", Join)
	builtins["split"] = object.NewBuiltin("strings.split", Split)
	builtins["fields"] = object.NewBuiltin("strings.fields", Fields)
	builtins["index"] = object.NewBuiltin("strings.index", Index)
	builtins["last_index"] = object.NewBuiltin("strings.last_index", LastIndex)
	builtins["replace_all"] = object.NewBuiltin("strings.replace_all", ReplaceAll)
	builtins["to_lower"] = object.NewBuiltin("strings.to_lower", ToLower)
	builtins["to_upper"] = object.NewBuiltin("strings.to_upper", ToUpper)
	builtins["trim"] = object.NewBuiltin("strings.trim", Trim)
	builtins["trim_prefix"] = object.NewBuiltin("strings.trim_prefix", TrimPrefix)
	builtins["trim_suffix"] = object.NewBuiltin("strings.trim_suffix", TrimSuffix)
	builtins["trim_space"] = object.NewBuiltin("strings.trim_space", TrimSpace)
	return builtins
}

// The "Module()" function can be disabled with "//risor:generate no-module-func"

// Module returns the Risor module object with all the associated builtin
// functions.
func Module() *object.Module {
	return object.NewBuiltinsModule("strings", addGeneratedBuiltins(map[string]object.Object{}))
}

