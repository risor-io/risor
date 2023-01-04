package evaluator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/internal/httputil"
	"github.com/cloudcmds/tamarin/object"
)

// Len returns the length of a string, list, set, or map
func Len(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("len", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.String:
		return object.NewInt(int64(utf8.RuneCountInString(arg.Value())))
	case *object.List:
		return object.NewInt(int64(len(arg.Value())))
	case *object.Set:
		return object.NewInt(int64(len(arg.Value())))
	case *object.Map:
		return object.NewInt(int64(len(arg.Value())))
	default:
		return object.Errorf("type error: len() argument is unsupported (%s given)", args[0].Type())
	}
}

func Sprintf(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.Errorf("type error: sprintf() takes 1 or more arguments (%d given)", len(args))
	}
	fs, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	fmtArgs := make([]interface{}, len(args)-1)
	for i, v := range args[1:] {
		fmtArgs[i] = v.Interface()
	}
	return object.NewString(fmt.Sprintf(fs, fmtArgs...))
}

func Delete(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("delete", 2, args); err != nil {
		return err
	}
	container, ok := args[0].(object.Container)
	if !ok {
		return object.Errorf("type error: delete() argument is unsupported (%s given)", args[0].Type())
	}
	if err := container.DelItem(args[1]); err != nil {
		return err
	}
	return object.Nil
}

func Set(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: set() expected at most 1 argument (%d given)", nArgs)
	}
	set := object.NewSetWithSize(0)
	if nArgs == 0 {
		return set
	}
	arg := args[0]
	switch arg := arg.(type) {
	case *object.String:
		for _, v := range arg.Value() {
			if res := set.Add(object.NewString(string(v))); object.IsError(res) {
				return res
			}
		}
	case *object.List:
		for _, obj := range arg.Value() {
			if res := set.Add(obj); object.IsError(res) {
				return res
			}
		}
	case *object.Set:
		for _, obj := range arg.Value() {
			if res := set.Add(obj); object.IsError(res) {
				return res
			}
		}
	case *object.Map:
		for k := range arg.Value() {
			if res := set.Add(object.NewString(k)); object.IsError(res) {
				return res
			}
		}
	default:
		return object.Errorf("type error: set() argument is unsupported (%s given)", args[0].Type())
	}
	return set
}

func List(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: list() expected at most 1 argument (%d given)", nArgs)
	}
	if nArgs == 0 {
		return object.NewList(nil)
	}
	switch obj := args[0].(type) {
	case *object.String:
		var items []object.Object
		for _, v := range obj.Value() {
			items = append(items, object.NewString(string(v)))
		}
		return object.NewList(items)
	case *object.List:
		return obj.Copy()
	case *object.Set:
		return object.NewList(obj.SortedItems())
	case *object.Map:
		return obj.Keys()
	default:
		return object.Errorf("type error: list() argument is unsupported (%s given)", args[0].Type())
	}
}

func Map(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: map() expected at most 1 argument (%d given)", nArgs)
	}
	result := object.NewMap(nil)
	if nArgs == 0 {
		return result
	}
	list, ok := args[0].(*object.List)
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
	container, ok := args[0].(object.Container)
	if !ok {
		return object.Errorf("type error: map() argument is unsupported (%s given)", args[0].Type())
	}
	iter := container.Iter()
	for {
		entry, ok := iter.Next()
		if !ok {
			break
		}
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
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: string() expected at most 1 argument (%d given)", nArgs)
	}
	if nArgs == 0 {
		return object.NewString("")
	}
	switch arg := args[0].(type) {
	case *object.String:
		return object.NewString(arg.Value())
	default:
		return object.NewString(args[0].Inspect())
	}
}

func Type(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("type", 1, args); err != nil {
		return err
	}
	return object.NewString(string(args[0].Type()))
}

func Ok(ctx context.Context, args ...object.Object) object.Object {
	switch len(args) {
	case 0:
		return object.NewOkResult(object.Nil)
	case 1:
		return object.NewOkResult(args[0])
	default:
		return object.Errorf("type error: ok() takes 0 or 1 arguments (%d given)", len(args))
	}
}

func Err(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 20 {
		return object.NewArgsRangeError("err", 1, 20, len(args))
	}
	switch obj := args[0].(type) {
	case *object.Error:
		return object.NewErrResult(obj)
	case *object.String:
		var extraArgs []interface{}
		for _, arg := range args[1:] {
			extraArgs = append(extraArgs, arg.Interface())
		}
		return object.NewErrResult(object.Errorf(obj.Value(), extraArgs...))
	default:
		return object.Errorf("type error: err() argument is unsupported (%s given)", args[0].Type())
	}
}

func Assert(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.Errorf("type error: assert() takes 1 or 2 arguments (%d given)", len(args))
	}
	if !args[0].IsTruthy() {
		if numArgs == 2 {
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
		for _, obj := range arg.Value() {
			if obj.IsTruthy() {
				return object.True
			}
		}
	case *object.Set:
		for _, obj := range arg.Value() {
			if obj.IsTruthy() {
				return object.True
			}
		}
	default:
		return object.Errorf("type error: any() argument must be an array, hash, or set (%s given)", args[0].Type())
	}
	return object.False
}

func All(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("all", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.List:
		for _, obj := range arg.Value() {
			if !obj.IsTruthy() {
				return object.False
			}
		}
	case *object.Set:
		for _, obj := range arg.Value() {
			if !obj.IsTruthy() {
				return object.False
			}
		}
	default:
		return object.Errorf("type error: all() argument must be an array, hash, or set (%s given)", args[0].Type())
	}
	return object.True
}

func Bool(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: bool() expected at most 1 argument (%d given)", nArgs)
	}
	if nArgs == 0 {
		return object.False
	}
	if args[0].IsTruthy() {
		return object.True
	}
	return object.False
}

func Fetch(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsRangeError("fetch", 1, 2, len(args))
	}
	urlArg, argErr := object.AsString(args[0])
	if argErr != nil {
		return argErr
	}
	var errObj *object.Error
	var params *object.Map
	if numArgs == 2 {
		params, errObj = object.AsMap(args[1])
		if errObj != nil {
			return errObj
		}
	}
	client := &http.Client{Timeout: 3 * time.Second}
	req, timeout, errObj := httputil.NewRequestFromParams(ctx, urlArg, params)
	if errObj != nil {
		return object.NewErrResult(errObj)
	}
	if timeout != 0 {
		client.Timeout = timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return object.NewErrResult(object.NewError(err))
	}
	return object.NewOkResult(object.NewHttpResponse(resp))
}

// output a string to stdout
func Print(ctx context.Context, args ...object.Object) object.Object {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		switch arg := arg.(type) {
		case *object.String:
			values[i] = arg.Value()
		default:
			values[i] = arg.Inspect()
		}
	}
	fmt.Println(values...)
	return object.Nil
}

func Printf(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 {
		return object.Errorf("type error: printf() takes 1 or more arguments (%d given)", len(args))
	}
	format, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var values []interface{}
	for _, arg := range args[1:] {
		switch arg := arg.(type) {
		case *object.String:
			values = append(values, arg.Value())
		default:
			values = append(values, arg.Interface())
		}
	}
	fmt.Printf(format, values...)
	return object.Nil
}

func Unwrap(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("unwrap", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Result:
		fn, ok := arg.GetAttr("unwrap")
		if !ok {
			return object.Errorf("type error: unwrap() method not found")
		}
		return fn.(*object.Builtin).Call(ctx, args[1:]...)
	default:
		return object.Errorf("type error: unwrap() argument must be a result (%s given)", arg.Type())
	}
}

func UnwrapOr(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("unwrap_or", 2, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Result:
		fn, ok := arg.GetAttr("unwrap_or")
		if !ok {
			return object.Errorf("type error: unwrap_or() method not found")
		}
		return fn.(*object.Builtin).Call(ctx, args[1:]...)
	default:
		return object.Errorf("type error: unwrap_or() argument must be a result (%s given)", arg.Type())
	}
}

func Sorted(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("sorted", 1, args); err != nil {
		return err
	}
	var items []object.Object
	switch arg := args[0].(type) {
	case *object.List:
		items = arg.Value()
	case *object.Map:
		items = arg.Keys().Value()
	case *object.Set:
		items = arg.List().Value()
	case *object.String:
		items = arg.Runes()
	default:
		return object.Errorf("type error: sorted() argument must be a list, map, or set (%s given)", arg.Type())
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
	switch arg := args[0].(type) {
	case *object.List:
		return arg.Reversed()
	case *object.String:
		return arg.Reversed()
	default:
		return object.Errorf("type error: reversed() argument must be an array or string (%s given)", arg.Type())
	}
}

func GetAttr(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 2 || numArgs > 3 {
		return object.Errorf("type error: getattr() takes 2 or 3 arguments (%d given)", len(args))
	}
	attrName, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	if attr, found := args[0].GetAttr(attrName); found {
		return attr
	}
	if numArgs == 3 {
		return args[2]
	}
	return object.Errorf("attribute error: %s object has no attribute %q", args[0].Type(), attrName)
}

// Call the given function with the provided arguments
func Call(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 {
		return object.Errorf("type error: call() takes 1 or more arguments (%d given)", len(args))
	}
	switch fn := args[0].(type) {
	case *object.Builtin:
		return fn.Call(ctx, args[1:]...)
	case *object.Function:
		callFunc, found := object.GetCallFunc(ctx)
		if !found {
			return object.Errorf("eval error: context did not contain a call function")
		}
		return callFunc(ctx, fn.Scope(), fn, args[1:])
	}
	return object.Errorf("type error: unable to call object (%s given)", args[0].Type())
}

func Keys(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("keys", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Map:
		return arg.Keys()
	case *object.List:
		return arg.Keys()
	default:
		return object.Errorf("type error: keys() argument must be a map or list (%s given)", arg.Type())
	}
}

func Int(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: int() expected at most 1 argument (%d given)", nArgs)
	}
	if nArgs == 0 {
		return object.NewInt(0)
	}
	switch obj := args[0].(type) {
	case *object.Int:
		return obj
	case *object.Float:
		return object.NewInt(int64(obj.Value()))
	case *object.String:
		if i, err := strconv.ParseInt(obj.Value(), 0, 64); err == nil {
			return object.NewInt(i)
		}
		return object.Errorf("value error: invalid literal for int(): %q", obj.Value())
	}
	return object.Errorf("type error: int() argument must be a string, float, or int (%s given)", args[0].Type())
}

func Float(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: float() expected at most 1 argument (%d given)", nArgs)
	}
	if nArgs == 0 {
		return object.NewFloat(0)
	}
	switch obj := args[0].(type) {
	case *object.Int:
		return object.NewFloat(float64(obj.Value()))
	case *object.Float:
		return obj
	case *object.String:
		if f, err := strconv.ParseFloat(obj.Value(), 64); err == nil {
			return object.NewFloat(f)
		}
		return object.Errorf("value error: invalid literal for float(): %q", obj.Value())
	}
	return object.Errorf("type error: float() argument must be a string, float, or int (%s given)", args[0].Type())
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
	nArgs := len(args)
	if nArgs < 1 || nArgs > 20 {
		return object.NewArgsRangeError("error", 1, 20, len(args))
	}
	msg, ok := args[0].(*object.String)
	if !ok {
		return object.Errorf("type error: error() expected a string (%s given)", args[0].Type())
	}
	var goArgs []interface{}
	for _, arg := range args[1:] {
		goArgs = append(goArgs, arg.Interface())
	}
	return object.Errorf(msg.Value(), goArgs...)
}

func Try(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs < 1 || nArgs > 2 {
		return object.NewArgsRangeError("try", 1, 2, len(args))
	}
	handleErr := func(arg object.Object, err *object.Error) object.Object {
		switch obj := arg.(type) {
		case *object.Function:
			callFunc, found := object.GetCallFunc(ctx)
			if !found {
				return object.Errorf("eval error: context did not contain a call function")
			}
			return callFunc(ctx, obj.Scope(), obj, []object.Object{err.Message()})
		default:
			return obj
		}
	}
	switch obj := args[0].(type) {
	case *object.Error:
		if nArgs == 2 {
			return handleErr(args[1], obj)
		}
		return object.Nil
	case *object.Result:
		if obj.IsErr() {
			if nArgs == 2 {
				return handleErr(args[1], obj.UnwrapErr())
			}
			return object.Nil
		}
		return obj.Unwrap()
	default:
		return obj
	}
}

func Iter(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs != 1 {
		return object.NewArgsError("iter", 1, len(args))
	}
	container, ok := args[0].(object.Container)
	if !ok {
		return object.Errorf("type error: iter() expected a container (%s given)", args[0].Type())
	}
	return container.Iter()
}

func Exit(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs > 1 {
		return object.Errorf("type error: exit() expected at most 1 argument (%d given)", nArgs)
	}
	if nArgs == 0 {
		os.Exit(0)
	}
	switch obj := args[0].(type) {
	case *object.Int:
		os.Exit(int(obj.Value()))
	case *object.Error:
		os.Exit(1)
	}
	return object.Errorf("type error: exit() argument must be an int or error (%s given)", args[0].Type())
}

func GlobalBuiltins() []*object.Builtin {
	type builtin struct {
		name string
		fn   object.BuiltinFunction
	}
	builtins := []builtin{
		{"all", All},
		{"any", Any},
		{"assert", Assert},
		{"bool", Bool},
		{"call", Call},
		{"chr", Chr},
		{"delete", Delete},
		{"err", Err},
		{"error", Error},
		{"exit", Exit},
		{"fetch", Fetch},
		{"float", Float},
		{"getattr", GetAttr},
		{"int", Int},
		{"iter", Iter},
		{"keys", Keys},
		{"len", Len},
		{"list", List},
		{"map", Map},
		{"ok", Ok},
		{"ord", Ord},
		{"print", Print},
		{"printf", Printf},
		{"reversed", Reversed},
		{"set", Set},
		{"sorted", Sorted},
		{"sprintf", Sprintf},
		{"string", String},
		{"type", Type},
		{"unwrap_or", UnwrapOr},
		{"unwrap", Unwrap},
	}
	var ret []*object.Builtin
	for _, b := range builtins {
		ret = append(ret, object.NewBuiltin(b.name, b.fn))
	}
	ret = append(ret, object.NewErrorHandler("try", Try))
	return ret
}
