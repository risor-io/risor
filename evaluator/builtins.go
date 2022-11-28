package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/cloudcmds/tamarin/internal/arg"
	"github.com/cloudcmds/tamarin/object"
)

// length of a string, array, set, or hash
func Len(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("len", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.String:
		return object.NewInteger(int64(utf8.RuneCountInString(arg.Value)))
	case *object.Array:
		return object.NewInteger(int64(len(arg.Elements)))
	case *object.Set:
		return object.NewInteger(int64(len(arg.Items)))
	case *object.Hash:
		return object.NewInteger(int64(len(arg.Map)))
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
	newHash := object.NewHash(nil)
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
				newHash.Map[k] = v
			}
		}
		return newHash
	}
	return newHash
}

// Sprintf is the implementation of our `sprintf` function
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

// Get hash keys
func Keys(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("keys", 1, args); err != nil {
		return err
	}
	hash, err := object.AsHash(args[0])
	if err != nil {
		return err
	}
	return hash.Keys()
}

// Delete a given hash key
func Delete(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("delete", 2, args); err != nil {
		return err
	}
	hash, err := object.AsHash(args[0])
	if err != nil {
		return err
	}
	key, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	hash.Delete(key)
	return hash
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
	case *object.Array:
		for _, obj := range arg.Elements {
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
	case *object.Hash:
		for k := range arg.Map {
			if err := set.Add(object.NewString(k)); err != nil {
				return newError(err.Error())
			}
		}
	default:
		return newError("type error: set() argument is unsupported (%s given)", args[0].Type())
	}
	return set
}

func String(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("string", 1, args); err != nil {
		return err
	}
	return object.NewString(args[0].Inspect())
}

func Type(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("type", 1, args); err != nil {
		return err
	}
	switch args[0].(type) {
	case *object.String:
		return object.NewString("string")
	case *object.Regexp:
		return object.NewString("regexp")
	case *object.Boolean:
		return object.NewString("bool")
	case *object.Builtin:
		return object.NewString("builtin")
	case *object.Array:
		return object.NewString("array")
	case *object.Function:
		return object.NewString("function")
	case *object.Integer:
		return object.NewString("integer")
	case *object.Float:
		return object.NewString("float")
	case *object.Hash:
		return object.NewString("hash")
	case *object.Set:
		return object.NewString("set")
	case *object.Result:
		return object.NewString("result")
	case *object.HttpResponse:
		return object.NewString("http_response")
	case *object.Time:
		return object.NewString("time")
	case *object.Null:
		return object.NewString("null")
	case *object.DatabaseConnection:
		return object.NewString("db_connection")
	default:
		return newError("type error: type() argument not supported (%s given)", args[0].Type())
	}
}

func Ok(ctx context.Context, args ...object.Object) object.Object {
	switch len(args) {
	case 0:
		return &object.Result{Ok: object.NULL}
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
	return object.NULL
}

func Any(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("any", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Array:
		for _, obj := range arg.Elements {
			if isTruthy(obj) {
				return object.TRUE
			}
		}
	case *object.Hash:
		for _, v := range arg.Map {
			if isTruthy(v) {
				return object.TRUE
			}
		}
	case *object.Set:
		for _, obj := range arg.Items {
			if isTruthy(obj) {
				return object.TRUE
			}
		}
	default:
		return newError("type error: any() argument must be an array, hash, or set (%s given)", args[0].Type())
	}
	return object.FALSE
}

func All(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("all", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Array:
		for _, obj := range arg.Elements {
			if !isTruthy(obj) {
				return object.FALSE
			}
		}
	case *object.Hash:
		for _, v := range arg.Map {
			if !isTruthy(v) {
				return object.FALSE
			}
		}
	case *object.Set:
		for _, obj := range arg.Items {
			if !isTruthy(obj) {
				return object.FALSE
			}
		}
	default:
		return newError("type error: all() argument must be an array, hash, or set (%s given)", args[0].Type())
	}
	return object.TRUE
}

func Bool(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("bool", 1, args); err != nil {
		return err
	}
	if isTruthy(args[0]) {
		return object.TRUE
	}
	return object.FALSE
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
	var params *object.Hash
	if numArgs == 2 {
		objParams, ok := args[1].(*object.Hash)
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
		if timeout != object.NULL {
			switch value := timeout.(type) {
			case *object.Float:
				client.Timeout = time.Millisecond * time.Duration(value.Value*1000.0)
			case *object.Integer:
				client.Timeout = time.Second * time.Duration(value.Value)
			default:
				return newError("type error: timeout value should be an integer or float")
			}
		}
		if bodyObj := params.Get("body"); bodyObj != object.NULL {
			switch bodyObj := bodyObj.(type) {
			case *object.String:
				body = bytes.NewBufferString(bodyObj.Value)
			}
			// TODO: support more buffer and/or stream options
		}
		if headersObj := params.Get("headers"); headersObj != object.NULL {
			switch headersObj := headersObj.(type) {
			case *object.Hash:
				for k, v := range headersObj.Map {
					hdr.Add(k, v.Inspect())
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
		values[i] = arg.Inspect()
	}
	fmt.Println(values...)
	return object.NULL
}

// Printf is the implementation of our `printf` function.
func Printf(ctx context.Context, args ...object.Object) object.Object {
	// Convert to the formatted version, via our `sprintf` function
	out := Sprintf(ctx, args...)
	// If that returned a string then we can print it
	if out.Type() == object.STRING_OBJ {
		fmt.Print(out.(*object.String).Value)
	}
	return object.NULL
}

func Unwrap(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("unwrap", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Result:
		return arg.InvokeMethod("unwrap")
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
		return arg.InvokeMethod("unwrap_or", args[1])
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
	case *object.Array:
		items = arg.Elements
	case *object.Hash:
		items = arg.Keys().Elements
	case *object.Set:
		items = arg.Array().Elements
	default:
		return newError("type error: sorted() argument must be an array, hash, or set (%s given)", arg.Type())
	}
	var comparableErr error
	sort.SliceStable(items, func(a, b int) bool {
		itemA := items[a]
		itemB := items[b]
		compA, ok := itemA.(object.Comparable)
		if !ok {
			comparableErr = fmt.Errorf("type error: sorted() encountered a non-comparable item (%s)", itemA.Type())
		}
		if _, ok := itemB.(object.Comparable); !ok {
			comparableErr = fmt.Errorf("type error: sorted() encountered a non-comparable item (%s)", itemB.Type())
		}
		result, err := compA.Compare(itemB)
		if err != nil {
			comparableErr = err
		}
		return result == -1
	})
	if comparableErr != nil {
		return newError(comparableErr.Error())
	}
	return &object.Array{Elements: items}
}

func Reversed(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("reversed", 1, args); err != nil {
		return err
	}
	switch arg := args[0].(type) {
	case *object.Array:
		return arg.Reversed()
	case *object.String:
		return arg.Reversed()
	default:
		return newError("type error: reversed() argument must be an array or string (%s given)", arg.Type())
	}
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
		{Name: "fetch", Fn: Fetch},
		{Name: "any", Fn: Any},
		{Name: "all", Fn: All},
		{Name: "bool", Fn: Bool},
		{Name: "print", Fn: Print},
		{Name: "printf", Fn: Printf},
		{Name: "unwrap", Fn: Unwrap},
		{Name: "unwrap_or", Fn: UnwrapOr},
		{Name: "sorted", Fn: Sorted},
		{Name: "reversed", Fn: Reversed},
	}
}
