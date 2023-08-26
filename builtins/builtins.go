// Package builtins defines a default set of built-in functions.
package builtins

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"unicode"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/limits"
	"github.com/risor-io/risor/object"
)

func Len(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("len", 1, args); err != nil {
		return err
	}
	container, ok := args[0].(object.Container)
	if !ok {
		return object.Errorf("type error: len() unsupported argument (%s given)", args[0].Type())
	}
	return container.Len()
}

func Sprintf(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("sprintf", 1, 64, args); err != nil {
		return err
	}
	fs, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	fmtArgs := make([]interface{}, len(args)-1)
	for i, v := range args[1:] {
		fmtArgs[i] = v.Interface()
	}
	result := object.NewString(fmt.Sprintf(fs, fmtArgs...))
	if err := limits.TrackCost(ctx, result.Cost()); err != nil {
		return object.NewError(err)
	}
	return result
}

func Delete(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("delete", 2, args); err != nil {
		return err
	}
	container, ok := args[0].(object.Container)
	if !ok {
		return object.Errorf("type error: delete() unsupported argument (%s given)", args[0].Type())
	}
	if err := container.DelItem(args[1]); err != nil {
		return err
	}
	return object.Nil
}

func Set(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("set", 0, 1, args); err != nil {
		return err
	}
	set := object.NewSetWithSize(0)
	if len(args) == 0 {
		return set
	}
	arg := args[0]
	if err := limits.TrackCost(ctx, arg.Cost()); err != nil {
		return object.NewError(err)
	}
	iter, err := object.AsIterator(arg)
	if err != nil {
		return err
	}
	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		if res := set.Add(val); object.IsError(res) {
			return res
		}
	}
	return set
}

func List(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("list", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.NewList(nil)
	}
	if intObj, ok := args[0].(*object.Int); ok {
		count := intObj.Value()
		if count < 0 {
			return object.Errorf("value error: list() argument must be >= 0 (%d given)", count)
		}
		if err := limits.TrackCost(ctx, int(count)*8); err != nil {
			return object.NewError(err)
		}
		arr := make([]object.Object, count)
		for i := 0; i < int(count); i++ {
			arr[i] = object.Nil
		}
		return object.NewList(arr)
	}
	arg := args[0]
	if err := limits.TrackCost(ctx, arg.Cost()); err != nil {
		return object.NewError(err)
	}
	iter, err := object.AsIterator(arg)
	if err != nil {
		return err
	}
	var items []object.Object
	for {
		val, ok := iter.Next()
		if !ok {
			break
		}
		items = append(items, val)
	}
	return object.NewList(items)
}

func Map(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("map", 0, 1, args); err != nil {
		return err
	}
	result := object.NewMap(nil)
	if len(args) == 0 {
		return result
	}
	arg := args[0]
	if err := limits.TrackCost(ctx, arg.Cost()); err != nil {
		return object.NewError(err)
	}
	list, ok := arg.(*object.List)
	if ok {
		for _, obj := range list.Value() {
			subListObj, ok := obj.(*object.List)
			if !ok || len(subListObj.Value()) != 2 {
				return object.Errorf("type error: map() received a list with an unsupported structure")
			}
			subList := subListObj.Value()
			key, ok := subList[0].(*object.String)
			if !ok {
				return object.Errorf("type error: map() received a list with an unsupported structure")
			}
			result.Set(key.Value(), subList[1])
		}
		return result
	}
	iter, err := object.AsIterator(arg)
	if err != nil {
		return err
	}
	for {
		if _, ok := iter.Next(); !ok {
			break
		}
		entry, _ := iter.Entry()
		k, v := entry.Key(), entry.Value()
		switch k := k.(type) {
		case *object.String:
			result.Set(k.Value(), v)
		default:
			result.Set(k.Inspect(), v)
		}
	}
	return result
}

func String(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("string", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.NewString("")
	}
	arg := args[0]
	argCost := arg.Cost()
	lim, ok := limits.GetLimits(ctx)
	if !ok {
		return object.NewError(limits.LimitsNotFound)
	}
	switch arg := arg.(type) {
	case *object.Buffer:
		if err := lim.TrackCost(argCost); err != nil {
			return object.NewError(err)
		}
		return object.NewString(string(arg.Value().Bytes()))
	case *object.ByteSlice:
		if err := lim.TrackCost(argCost); err != nil {
			return object.NewError(err)
		}
		return object.NewString(string(arg.Value()))
	case *object.String:
		if err := lim.TrackCost(argCost); err != nil {
			return object.NewError(err)
		}
		return object.NewString(arg.Value())
	case io.Reader:
		bytes, err := lim.ReadAll(arg)
		if err != nil {
			return object.NewError(err)
		}
		return object.NewString(string(bytes))
	default:
		if s, ok := arg.(fmt.Stringer); ok {
			return object.NewString(s.String())
		}
		return object.NewString(args[0].Inspect())
	}
}

func FloatSlice(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("float_slice", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.NewFloatSlice(nil)
	}
	arg := args[0]
	argCost := arg.Cost()
	if err := limits.TrackCost(ctx, argCost); err != nil {
		return object.NewError(err)
	}
	switch arg := arg.(type) {
	case *object.FloatSlice:
		return arg.Clone()
	case *object.Int:
		val := arg.Value()
		if err := limits.TrackCost(ctx, int(val)-argCost); err != nil {
			return object.NewError(err)
		}
		return object.NewFloatSlice(make([]float64, val))
	case *object.List:
		items := arg.Value()
		floats := make([]float64, len(items))
		for i, item := range items {
			switch item := item.(type) {
			case *object.Byte:
				floats[i] = float64(item.Value())
			case *object.Int:
				floats[i] = float64(item.Value())
			case *object.Float:
				floats[i] = item.Value()
			default:
				return object.Errorf("type error: float_slice() list item unsupported (%s given)", item.Type())
			}
		}
		return object.NewFloatSlice(floats)
	default:
		return object.Errorf("type error: float_slice() unsupported argument (%s given)", args[0].Type())
	}
}

func ByteSlice(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("byte_slice", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.NewByteSlice(nil)
	}
	arg := args[0]
	argCost := arg.Cost()
	if err := limits.TrackCost(ctx, argCost); err != nil {
		return object.NewError(err)
	}
	switch arg := arg.(type) {
	case *object.Buffer:
		return object.NewByteSlice(arg.Value().Bytes())
	case *object.ByteSlice:
		return arg.Clone()
	case *object.String:
		return object.NewByteSlice([]byte(arg.Value()))
	case *object.Int:
		val := arg.Value()
		if err := limits.TrackCost(ctx, int(val)-argCost); err != nil {
			return object.NewError(err)
		}
		return object.NewByteSlice(make([]byte, val))
	case *object.List:
		items := arg.Value()
		bytes := make([]byte, len(items))
		for i, item := range items {
			switch item := item.(type) {
			case *object.Int:
				bytes[i] = byte(item.Value())
			case *object.Byte:
				bytes[i] = item.Value()
			default:
				return object.Errorf("type error: byte_slice() list item unsupported (%s given)", item.Type())
			}
		}
		return object.NewByteSlice(bytes)
	default:
		return object.Errorf("type error: byte_slice() unsupported argument (%s given)", args[0].Type())
	}
}

func Buffer(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("buffer", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.NewBuffer(nil)
	}
	arg := args[0]
	lim, ok := limits.GetLimits(ctx)
	if !ok {
		return object.NewError(limits.LimitsNotFound)
	}
	switch arg := arg.(type) {
	case *object.Buffer:
		if err := lim.TrackCost(arg.Cost()); err != nil {
			return object.NewError(err)
		}
		return object.NewBufferFromBytes(arg.Value().Bytes())
	case *object.ByteSlice:
		if err := lim.TrackCost(arg.Cost()); err != nil {
			return object.NewError(err)
		}
		return object.NewBufferFromBytes(arg.Value())
	case *object.String:
		if err := lim.TrackCost(arg.Cost()); err != nil {
			return object.NewError(err)
		}
		return object.NewBufferFromBytes([]byte(arg.Value()))
	case *object.Int:
		// Special case: treat the value as the size to allocate
		val := arg.Value()
		if err := lim.TrackCost(int(val)); err != nil {
			return object.NewError(err)
		}
		return object.NewBufferFromBytes(make([]byte, val))
	case io.Reader:
		bytes, err := lim.ReadAll(arg)
		if err != nil {
			return object.NewError(err)
		}
		return object.NewBufferFromBytes(bytes)
	default:
		return object.Errorf("type error: buffer() unsupported argument (%s given)", args[0].Type())
	}
}

func Type(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("type", 1, args); err != nil {
		return err
	}
	return object.NewString(string(args[0].Type()))
}

func Assert(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("assert", 1, 2, args); err != nil {
		return err
	}
	if !args[0].IsTruthy() {
		if len(args) == 2 {
			switch arg := args[1].(type) {
			case *object.String:
				return object.Errorf(arg.Value())
			default:
				return object.Errorf(args[1].Inspect())
			}
		}
		return object.Errorf("assertion failed")
	}
	return object.Nil
}

func Any(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("any", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.List:
		for _, val := range arg.Value() {
			if val.IsTruthy() {
				return object.True
			}
		}
	case *object.Set:
		for _, val := range arg.Value() {
			if val.IsTruthy() {
				return object.True
			}
		}
	case *object.ByteSlice:
		for _, val := range arg.Value() {
			if val != 0 {
				return object.True
			}
		}
	case *object.Buffer:
		for _, val := range arg.Value().Bytes() {
			if val != 0 {
				return object.True
			}
		}
	case object.Iterable:
		iter := arg.Iter()
		for {
			val, ok := iter.Next()
			if !ok {
				break
			}
			if val.IsTruthy() {
				return object.True
			}
		}
	case object.Iterator:
		for {
			val, ok := arg.Next()
			if !ok {
				break
			}
			if val.IsTruthy() {
				return object.True
			}
		}
	default:
		return object.Errorf("type error: any() argument must be a container (%s given)", args[0].Type())
	}
	return object.False
}

func All(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("all", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.List:
		for _, val := range arg.Value() {
			if !val.IsTruthy() {
				return object.False
			}
		}
	case *object.Set:
		for _, val := range arg.Value() {
			if !val.IsTruthy() {
				return object.False
			}
		}
	case *object.ByteSlice:
		for _, val := range arg.Value() {
			if val == 0 {
				return object.False
			}
		}
	case *object.Buffer:
		for _, val := range arg.Value().Bytes() {
			if val == 0 {
				return object.False
			}
		}
	case object.Iterable:
		iter := arg.Iter()
		for {
			val, ok := iter.Next()
			if !ok {
				break
			}
			if !val.IsTruthy() {
				return object.False
			}
		}
	case object.Iterator:
		for {
			val, ok := arg.Next()
			if !ok {
				break
			}
			if !val.IsTruthy() {
				return object.False
			}
		}
	default:
		return object.Errorf("type error: all() argument must be a container (%s given)", args[0].Type())
	}
	return object.True
}

func Bool(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("bool", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.False
	}
	if args[0].IsTruthy() {
		return object.True
	}
	return object.False
}

func Sorted(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("sorted", 1, args); err != nil {
		return err
	}
	arg := args[0]
	if err := limits.TrackCost(ctx, arg.Cost()); err != nil {
		return object.NewError(err)
	}
	var items []object.Object
	switch arg := arg.(type) {
	case *object.List:
		items = arg.Value()
	case *object.Map:
		items = arg.Keys().Value()
	case *object.Set:
		items = arg.List().Value()
	case *object.String:
		items = arg.Runes()
	case *object.ByteSlice:
		items = arg.Integers()
	default:
		return object.Errorf("type error: sorted() unsupported argument (%s given)", arg.Type())
	}
	resultItems := make([]object.Object, len(items))
	copy(resultItems, items)
	if err := object.Sort(resultItems); err != nil {
		return err
	}
	return object.NewList(resultItems)
}

func Reversed(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("reversed", 1, args); err != nil {
		return err
	}
	arg := args[0]
	if err := limits.TrackCost(ctx, arg.Cost()); err != nil {
		return object.NewError(err)
	}
	switch arg := arg.(type) {
	case *object.List:
		return arg.Reversed()
	case *object.String:
		return arg.Reversed()
	case *object.ByteSlice:
		return arg.Reversed()
	default:
		return object.Errorf("type error: reversed() unsupported argument (%s given)", arg.Type())
	}
}

func GetAttr(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("getattr", 2, 3, args); err != nil {
		return err
	}
	attrName, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	if attr, found := args[0].GetAttr(attrName); found {
		return attr
	}
	if len(args) == 3 {
		return args[2]
	}
	return object.Errorf("type error: getattr() %s object has no attribute %q",
		args[0].Type(), attrName)
}

func Call(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("call", 1, 64, args); err != nil {
		return err
	}
	switch fn := args[0].(type) {
	case *object.Function:
		callFunc, found := object.GetCallFunc(ctx)
		if !found {
			return object.Errorf("eval error: context did not contain a call function")
		}
		result, err := callFunc(ctx, fn, args[1:])
		if err != nil {
			return object.Errorf(err.Error())
		}
		return result
	case object.Callable:
		return fn.Call(ctx, args[1:]...)
	default:
		return object.Errorf("type error: call() unsupported argument (%s given)", args[0].Type())
	}
}

func Keys(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("keys", 1, args); err != nil {
		return err
	}
	arg := args[0]
	if err := limits.TrackCost(ctx, arg.Cost()); err != nil {
		return object.NewError(err)
	}
	switch arg := arg.(type) {
	case *object.Map:
		return arg.Keys()
	case *object.List:
		return arg.Keys()
	case *object.Set:
		return arg.List()
	case object.Iterable:
		return iterKeys(arg.Iter())
	case object.Iterator:
		return iterKeys(arg)
	default:
		return object.Errorf("type error: keys() unsupported argument (%s given)", arg.Type())
	}
}

func iterKeys(iter object.Iterator) object.Object {
	var keys []object.Object
	for {
		if _, ok := iter.Next(); !ok {
			break
		}
		entry, _ := iter.Entry()
		keys = append(keys, entry.Key())
	}
	return object.NewList(keys)
}

func Byte(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("byte", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.NewByte(0)
	}
	switch obj := args[0].(type) {
	case *object.Int:
		return object.NewByte(byte(obj.Value()))
	case *object.Byte:
		return object.NewByte(obj.Value())
	case *object.Float:
		return object.NewByte(byte(obj.Value()))
	case *object.String:
		if i, err := strconv.ParseInt(obj.Value(), 0, 8); err == nil {
			return object.NewByte(byte(i))
		}
		return object.Errorf("value error: invalid literal for byte(): %q", obj.Value())
	default:
		return object.Errorf("type error: byte() unsupported argument (%s given)", args[0].Type())
	}
}

func Int(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("int", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.NewInt(0)
	}
	switch obj := args[0].(type) {
	case *object.Int:
		return obj
	case *object.Byte:
		return object.NewInt(int64(obj.Value()))
	case *object.Float:
		return object.NewInt(int64(obj.Value()))
	case *object.Duration:
		return object.NewInt(int64(obj.Value()))
	case *object.String:
		if i, err := strconv.ParseInt(obj.Value(), 0, 64); err == nil {
			return object.NewInt(i)
		}
		return object.Errorf("value error: invalid literal for int(): %q", obj.Value())
	default:
		return object.Errorf("type error: int() unsupported argument (%s given)", args[0].Type())
	}
}

func Float(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("float", 0, 1, args); err != nil {
		return err
	}
	if len(args) == 0 {
		return object.NewFloat(0)
	}
	switch obj := args[0].(type) {
	case *object.Int:
		return object.NewFloat(float64(obj.Value()))
	case *object.Byte:
		return object.NewFloat(float64(obj.Value()))
	case *object.Float:
		return obj
	case *object.String:
		if f, err := strconv.ParseFloat(obj.Value(), 64); err == nil {
			return object.NewFloat(f)
		}
		return object.Errorf("value error: invalid literal for float(): %q", obj.Value())
	default:
		return object.Errorf("type error: float() unsupported argument (%s given)", args[0].Type())
	}
}

func Ord(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("ord", 1, args); err != nil {
		return err
	}
	switch obj := args[0].(type) {
	case *object.String:
		runes := []rune(obj.Value())
		if len(runes) != 1 {
			return object.Errorf("value error: ord() expected a character, but string of length %d found", len(obj.Value()))
		}
		return object.NewInt(int64(runes[0]))
	}
	return object.Errorf("type error: ord() expected a string of length 1 (%s given)", args[0].Type())
}

func Chr(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("chr", 1, args); err != nil {
		return err
	}
	switch obj := args[0].(type) {
	case *object.Int:
		v := obj.Value()
		if v < 0 {
			return object.Errorf("value error: chr() argument out of range (%d given)", v)
		}
		if v > unicode.MaxRune {
			return object.Errorf("value error: chr() argument out of range (%d given)", v)
		}
		return object.NewString(string(rune(v)))
	}
	return object.Errorf("type error: chr() expected an int (%s given)", args[0].Type())
}

func Error(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("error", 1, 64, args); err != nil {
		return err
	}
	msg, ok := args[0].(*object.String)
	if !ok {
		return object.Errorf("type error: error() expected a string (%s given)", args[0].Type())
	}
	var msgArgs []interface{}
	for _, arg := range args[1:] {
		msgArgs = append(msgArgs, arg.Interface())
	}
	return object.Errorf(msg.Value(), msgArgs...)
}

func Try(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("try", 1, 64, args); err != nil {
		return err
	}
	try := func(arg object.Object) (object.Object, error) {
		switch obj := arg.(type) {
		case *object.Function:
			callFunc, found := object.GetCallFunc(ctx)
			if !found {
				return nil, fmt.Errorf("eval error: context did not contain a call function")
			}
			result, err := callFunc(ctx, obj, nil)
			if err != nil {
				return nil, err
			}
			switch result := result.(type) {
			case *object.Error:
				return nil, result.Value()
			default:
				return result, nil
			}
		case object.Callable:
			result := obj.Call(ctx)
			switch result := result.(type) {
			case *object.Error:
				return nil, result.Value()
			default:
				return result, nil
			}
		default:
			return obj, nil
		}
	}
	for _, arg := range args {
		result, err := try(arg)
		if err == nil {
			return result
		}
	}
	return object.Nil
}

func Iter(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("iter", 1, args); err != nil {
		return err
	}
	container, ok := args[0].(object.Iterable)
	if !ok {
		return object.Errorf("type error: iter() expected an iterable (%s given)", args[0].Type())
	}
	return container.Iter()
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"all":         object.NewBuiltin("all", All),
		"any":         object.NewBuiltin("any", Any),
		"assert":      object.NewBuiltin("assert", Assert),
		"bool":        object.NewBuiltin("bool", Bool),
		"buffer":      object.NewBuiltin("buffer", Buffer),
		"byte_slice":  object.NewBuiltin("byte_slice", ByteSlice),
		"byte":        object.NewBuiltin("byte", Byte),
		"call":        object.NewBuiltin("call", Call),
		"chr":         object.NewBuiltin("chr", Chr),
		"decode":      object.NewBuiltin("decode", Decode),
		"delete":      object.NewBuiltin("delete", Delete),
		"encode":      object.NewBuiltin("encode", Encode),
		"error":       object.NewBuiltin("error", Error),
		"float_slice": object.NewBuiltin("float_slice", FloatSlice),
		"float":       object.NewBuiltin("float", Float),
		"getattr":     object.NewBuiltin("getattr", GetAttr),
		"int":         object.NewBuiltin("int", Int),
		"iter":        object.NewBuiltin("iter", Iter),
		"keys":        object.NewBuiltin("keys", Keys),
		"len":         object.NewBuiltin("len", Len),
		"list":        object.NewBuiltin("list", List),
		"map":         object.NewBuiltin("map", Map),
		"ord":         object.NewBuiltin("ord", Ord),
		"reversed":    object.NewBuiltin("reversed", Reversed),
		"set":         object.NewBuiltin("set", Set),
		"sorted":      object.NewBuiltin("sorted", Sorted),
		"sprintf":     object.NewBuiltin("sprintf", Sprintf),
		"string":      object.NewBuiltin("string", String),
		"try":         object.NewBuiltin("try", Try),
		"type":        object.NewBuiltin("type", Type),
	}
}
