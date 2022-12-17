package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/cloudcmds/tamarin/arg"
	"github.com/cloudcmds/tamarin/object"
)

// Len returns the length of a string, list, set, or map
func Len(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("len", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.String:
		return object.NewInt(int64(utf8.RuneCountInString(arg.Value)))
	case *object.List:
		return object.NewInt(int64(len(arg.Items)))
	case *object.Set:
		return object.NewInt(int64(len(arg.Items)))
	case *object.Map:
		return object.NewInt(int64(len(arg.Items)))
	default:
		return newError("type error: len() argument is unsupported (%s given)", args[0].Type())
	}
}

// regular expression match
func Match(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("match", 2, args); err != nil {
		return err
	}
	arg0, argErr := object.AsString(args[0])
	if argErr != nil {
		return argErr
	}
	arg1, argErr := object.AsString(args[1])
	if argErr != nil {
		return argErr
	}
	reg := regexp.MustCompile(arg0)
	res := reg.FindStringSubmatch(arg1)
	newMap := object.NewMap(nil)
	if len(res) > 0 {
		//
		// If we get a match then the output is an array
		// First entry is the match, any additional parts
		// are the capture-groups.
		//
		if len(res) > 1 {
			for i := 1; i < len(res); i++ {
				// Capture groups start at index 0.
				k := fmt.Sprintf("%d", int64(i-1))
				v := object.NewString(res[i])
				newMap.Items[k] = v
			}
		}
		return newMap
	}
	return newMap
}

func Sprintf(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return newError("type error: sprintf() takes 1 or more arguments (%d given)", len(args))
	}
	fs, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	fmtArgs := make([]interface{}, len(args)-1)
	for i, v := range args[1:] {
		fmtArgs[i] = v.ToInterface()
	}
	return object.NewString(fmt.Sprintf(fs, fmtArgs...))
}

func Delete(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("delete", 2, args); err != nil {
		return err
	}
	m, err := object.AsMap(args[0])
	if err != nil {
		return err
	}
	key, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	m.Delete(key)
	return m
}

func Set(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("set", 1, args); err != nil {
		return err
	}
	set := object.NewSetWithSize(0)
	arg := args[0]
	switch arg := arg.(type) {
	case *object.String:
		for _, v := range arg.Value {
			set.Add(object.NewString(string(v)))
		}
	case *object.List:
		for _, obj := range arg.Items {
			if err := set.Add(obj); err != nil {
				return newError(err.Error())
			}
		}
	case *object.Set:
		for _, obj := range arg.Items {
			if err := set.Add(obj); err != nil {
				return newError(err.Error())
			}
		}
	case *object.Map:
		for k := range arg.Items {
			if err := set.Add(object.NewString(k)); err != nil {
				return newError(err.Error())
			}
		}
	default:
		return newError("type error: set() argument is unsupported (%s given)", args[0].Type())
	}
	return set
}

func List(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("list", 1, args); err != nil {
		return err
	}
	switch obj := args[0].(type) {
	case *object.String:
		var items []object.Object
		for _, v := range obj.Value {
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
		return newError("type error: list() argument is unsupported (%s given)", args[0].Type())
	}
}

func String(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("string", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.String:
		return object.NewString(arg.Value)
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
		return &object.Result{Ok: object.Nil}
	case 1:
		return &object.Result{Ok: args[0]}
	default:
		return newError("type error: ok() takes 0 or 1 arguments (%d given)", len(args))
	}
}

func Err(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("err", 1, args); err != nil {
		return err
	}
	return &object.Result{Err: args[0]}
}

func Assert(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return newError("type error: assert() takes 1 or 2 arguments (%d given)", len(args))
	}
	if !isTruthy(args[0]) {
		if numArgs == 2 {
			return newError(args[1].Inspect())
		}
		return newError("assertion failed")
	}
	return object.Nil
}

func Any(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("any", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.List:
		for _, obj := range arg.Items {
			if isTruthy(obj) {
				return object.True
			}
		}
	case *object.Set:
		for _, obj := range arg.Items {
			if isTruthy(obj) {
				return object.True
			}
		}
	default:
		return newError("type error: any() argument must be an array, hash, or set (%s given)", args[0].Type())
	}
	return object.False
}

func All(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("all", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.List:
		for _, obj := range arg.Items {
			if !isTruthy(obj) {
				return object.False
			}
		}
	case *object.Set:
		for _, obj := range arg.Items {
			if !isTruthy(obj) {
				return object.False
			}
		}
	default:
		return newError("type error: all() argument must be an array, hash, or set (%s given)", args[0].Type())
	}
	return object.True
}

func Bool(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bool", 1, args); err != nil {
		return err
	}
	if isTruthy(args[0]) {
		return object.True
	}
	return object.False
}

func Fetch(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return newError("type error: fetch() takes 1 or 2 arguments (%d given)", len(args))
	}
	urlArg, argErr := object.AsString(args[0])
	if argErr != nil {
		return argErr
	}
	var params *object.Map
	if numArgs == 2 {
		objParams, ok := args[1].(*object.Map)
		if !ok {
			return newError("type error: expected a hash argument; got %v", args[1].Type())
		}
		params = objParams
	}
	client := &http.Client{Timeout: 10 * time.Second}
	method := "GET"
	var body io.Reader
	hdr := http.Header{}
	if params != nil {
		if value, ok := params.Get("method").(*object.String); ok {
			method = value.Value
		}
		timeout := params.Get("timeout")
		if timeout != object.Nil {
			switch value := timeout.(type) {
			case *object.Float:
				client.Timeout = time.Millisecond * time.Duration(value.Value*1000.0)
			case *object.Int:
				client.Timeout = time.Second * time.Duration(value.Value)
			default:
				return newError("type error: timeout value should be an integer or float")
			}
		}
		if bodyObj := params.Get("body"); bodyObj != object.Nil {
			switch bodyObj := bodyObj.(type) {
			case *object.String:
				body = bytes.NewBufferString(bodyObj.Value)
			}
			// TODO: support more buffer and/or stream options
		}
		if headersObj := params.Get("headers"); headersObj != object.Nil {
			switch headersObj := headersObj.(type) {
			case *object.Map:
				for k, v := range headersObj.Items {
					switch v := v.(type) {
					case *object.String:
						hdr.Add(k, v.Value)
					default:
						hdr.Add(k, v.Inspect())
					}
				}
			}
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, urlArg, body)
	if err != nil {
		return &object.Result{Err: &object.Error{Message: err.Error()}}
	}
	for k, values := range hdr {
		for _, value := range values {
			req.Header.Add(k, value)
		}
	}
	// req.Header = hdr
	resp, err := client.Do(req)
	if err != nil {
		return &object.Result{Err: &object.Error{Message: err.Error()}}
	}
	return &object.Result{Ok: &object.HttpResponse{Response: resp}}
}

// output a string to stdout
func Print(ctx context.Context, args ...object.Object) object.Object {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		switch arg := arg.(type) {
		case *object.String:
			values[i] = arg.Value
		default:
			values[i] = arg.Inspect()
		}
	}
	fmt.Println(values...)
	return object.Nil
}

// Printf is the implementation of our `printf` function.
func Printf(ctx context.Context, args ...object.Object) object.Object {
	// Convert to the formatted version, via our `sprintf` function
	out := Sprintf(ctx, args...)
	// If that returned a string then we can print it
	if out.Type() == object.STRING {
		fmt.Print(out.(*object.String).Value)
	}
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
			return newError("type error: unwrap() method not found")
		}
		return fn.(*object.Builtin).Fn(ctx, args...)
	default:
		return newError("type error: unwrap() argument must be a result (%s given)", arg.Type())
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
			return newError("type error: unwrap_or() method not found")
		}
		return fn.(*object.Builtin).Fn(ctx, args...)
	default:
		return newError("type error: unwrap_or() argument must be a result (%s given)", arg.Type())
	}
}

func Sorted(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("sorted", 1, args); err != nil {
		return err
	}
	var items []object.Object
	switch arg := args[0].(type) {
	case *object.List:
		items = arg.Items
	case *object.Map:
		items = arg.Keys().Items
	case *object.Set:
		items = arg.List().Items
	default:
		return newError("type error: sorted() argument must be an array, hash, or set (%s given)", arg.Type())
	}
	result := &object.List{Items: make([]object.Object, len(items))}
	copy(result.Items, items)
	if err := object.Sort(result.Items); err != nil {
		return err
	}
	return result
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
		return newError("type error: reversed() argument must be an array or string (%s given)", arg.Type())
	}
}

func GetAttr(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 2 || numArgs > 3 {
		return newError("type error: getattr() takes 2 or 3 arguments (%d given)", len(args))
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
	return newError("attribute error: %s object has no attribute %q", args[0].Type(), attrName)
}

// Call the given builtin function with the provided arguments
func Call(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 {
		return newError("type error: call() takes 1 or more arguments (%d given)", len(args))
	}
	switch fn := args[0].(type) {
	case *object.Builtin:
		return fn.Fn(ctx, args...)
	case *object.Function:
		// TODO: pass in applyer
	}
	return newError("type error: unable to call object (%s given)", args[0].Type())
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
		return newError("type error: keys() argument must be a map or list (%s given)", arg.Type())
	}
}

func Int(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("int", 1, args); err != nil {
		return err
	}
	switch obj := args[0].(type) {
	case *object.Int:
		return obj
	case *object.Float:
		return object.NewInt(int64(obj.Value))
	case *object.String:
		if i, err := strconv.ParseInt(obj.Value, 0, 64); err == nil {
			return object.NewInt(i)
		}
		return newError("value error: invalid literal for int(): %q", obj.Value)
	}
	return newError("type error: int() argument must be a string, float, or int (%s given)", args[0].Type())
}

func Float(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("float", 1, args); err != nil {
		return err
	}
	switch obj := args[0].(type) {
	case *object.Int:
		return object.NewFloat(float64(obj.Value))
	case *object.Float:
		return obj
	case *object.String:
		if f, err := strconv.ParseFloat(obj.Value, 64); err == nil {
			return object.NewFloat(f)
		}
		return newError("value error: invalid literal for float(): %q", obj.Value)
	}
	return newError("type error: float() argument must be a string, float, or int (%s given)", args[0].Type())
}

func Ord(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("ord", 1, args); err != nil {
		return err
	}
	switch obj := args[0].(type) {
	case *object.String:
		runes := []rune(obj.Value)
		if len(runes) != 1 {
			return newError("value error: ord() expected a character, but string of length %d found", len(obj.Value))
		}
		return object.NewInt(int64(runes[0]))
	}
	return newError("type error: ord() expected a string of length 1 (%s given)", args[0].Type())
}

func Chr(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("chr", 1, args); err != nil {
		return err
	}
	switch obj := args[0].(type) {
	case *object.Int:
		if obj.Value < 0 {
			return newError("value error: chr() argument out of range (%d given)", obj.Value)
		}
		if obj.Value > unicode.MaxRune {
			return newError("value error: chr() argument out of range (%d given)", obj.Value)
		}
		return object.NewString(string(rune(obj.Value)))
	}
	return newError("type error: chr() expected an int (%s given)", args[0].Type())
}

func GlobalBuiltins() []*object.Builtin {
	return []*object.Builtin{
		{Name: "delete", Fn: Delete},
		{Name: "keys", Fn: Keys},
		{Name: "len", Fn: Len},
		{Name: "match", Fn: Match},
		{Name: "set", Fn: Set},
		{Name: "sprintf", Fn: Sprintf},
		{Name: "string", Fn: String},
		{Name: "type", Fn: Type},
		{Name: "ok", Fn: Ok},
		{Name: "err", Fn: Err},
		{Name: "assert", Fn: Assert},
		{Name: "any", Fn: Any},
		{Name: "all", Fn: All},
		{Name: "bool", Fn: Bool},
		{Name: "print", Fn: Print},
		{Name: "printf", Fn: Printf},
		{Name: "unwrap", Fn: Unwrap},
		{Name: "unwrap_or", Fn: UnwrapOr},
		{Name: "sorted", Fn: Sorted},
		{Name: "reversed", Fn: Reversed},
		{Name: "getattr", Fn: GetAttr},
		{Name: "call", Fn: Call},
		{Name: "list", Fn: List},
		{Name: "int", Fn: Int},
		{Name: "float", Fn: Float},
		{Name: "ord", Fn: Ord},
		{Name: "chr", Fn: Chr},
		{Name: "fetch", Fn: Fetch},
	}
}
