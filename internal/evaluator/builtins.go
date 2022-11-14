package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/cloudcmds/tamarin/internal/arg"
	"github.com/cloudcmds/tamarin/object"
)

// The built-in functions / standard-library methods are stored here.
var builtins = map[string]*object.Builtin{}

// length of a string, array, set, or hash
func lenFun(ctx context.Context, args ...object.Object) object.Object {
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
		return object.NewInteger(int64(len(arg.Pairs)))
	default:
		return newError("type error: len() argument is unsupported (%s given)", args[0].Type())
	}
}

// regular expression match
func matchFun(ctx context.Context, args ...object.Object) object.Object {
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
	if len(res) > 0 {
		newHash := make(map[object.HashKey]object.HashPair)
		//
		// If we get a match then the output is an array
		// First entry is the match, any additional parts
		// are the capture-groups.
		//
		if len(res) > 1 {
			for i := 1; i < len(res); i++ {
				// Capture groups start at index 0.
				k := object.NewInteger(int64(i - 1))
				v := object.NewString(res[i])
				newHashPair := object.HashPair{Key: k, Value: v}
				newHash[k.HashKey()] = newHashPair
			}
		}
		return &object.Hash{Pairs: newHash}
	}
	return object.NULL
}

// sprintfFun is the implementation of our `sprintf` function
func sprintfFun(ctx context.Context, args ...object.Object) object.Object {
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
func hashKeys(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("keys", 1, args); err != nil {
		return err
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("type error: keys() argument must be a hash (%s given)", args[0].Type())
	}
	// The object we're working with
	hash := args[0].(*object.Hash)
	ents := len(hash.Pairs)
	// Create a new array for the results.
	array := make([]object.Object, ents)
	// Now copy the keys into it.
	i := 0
	for _, ent := range hash.Pairs {
		array[i] = ent.Key
		i++
	}
	return &object.Array{Elements: array}
}

// Delete a given hash key
func hashDelete(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("delete", 2, args); err != nil {
		return err
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("type error: delete() argument must be a hash (%s given)", args[0].Type())
	}
	hash := args[0].(*object.Hash)
	key, ok := args[1].(object.Hashable)
	if !ok {
		return newError("type error: delete() key argument must be hashable (%s given)", args[1].Type())
	}
	newHash := make(map[object.HashKey]object.HashPair, len(hash.Pairs))
	for k, v := range hash.Pairs {
		if k != key.HashKey() {
			newHash[k] = v
		}
	}
	return &object.Hash{Pairs: newHash}
}

func setFun(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("set", 1, args); err != nil {
		return err
	}
	set := &object.Set{Items: map[object.HashKey]object.Object{}}
	arg := args[0]
	switch arg := arg.(type) {
	case *object.String:
		for _, v := range arg.Value {
			vStr := object.NewString(string(v))
			set.Items[vStr.HashKey()] = vStr
		}
	case *object.Array:
		for _, obj := range arg.Elements {
			hashable, ok := obj.(object.Hashable)
			if !ok {
				return newError("type error: set() argument contains an object that is not hashable (of type %s)", obj.Type())
			}
			set.Items[hashable.HashKey()] = obj
		}
	}
	return set
}

func strFun(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("str", 1, args); err != nil {
		return err
	}
	return object.NewString(args[0].Inspect())
}

func typeFun(ctx context.Context, args ...object.Object) object.Object {
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

func okFun(ctx context.Context, args ...object.Object) object.Object {
	switch len(args) {
	case 0:
		return &object.Result{Ok: object.NULL}
	case 1:
		return &object.Result{Ok: args[0]}
	default:
		return newError("type error: ok() takes 0 or 1 arguments (%d given)", len(args))
	}
}

func errFun(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("err", 1, args); err != nil {
		return err
	}
	return &object.Result{Err: args[0]}
}

func assertFun(ctx context.Context, args ...object.Object) object.Object {
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
		for _, ent := range arg.Pairs {
			if isTruthy(ent.Value) {
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
		for _, ent := range arg.Pairs {
			if !isTruthy(ent.Value) {
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

func fetchFun(ctx context.Context, args ...object.Object) object.Object {
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
				for _, v := range headersObj.Pairs {
					hdr.Add(v.Key.Inspect(), v.Value.Inspect())
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

func init() {
	RegisterBuiltin("delete", hashDelete)
	RegisterBuiltin("keys", hashKeys)
	RegisterBuiltin("len", lenFun)
	RegisterBuiltin("match", matchFun)
	RegisterBuiltin("set", setFun)
	RegisterBuiltin("sprintf", sprintfFun)
	RegisterBuiltin("str", strFun)
	RegisterBuiltin("type", typeFun)
	RegisterBuiltin("ok", okFun)
	RegisterBuiltin("err", errFun)
	RegisterBuiltin("assert", assertFun)
	RegisterBuiltin("fetch", fetchFun)
	RegisterBuiltin("any", Any)
	RegisterBuiltin("all", All)
	RegisterBuiltin("bool", Bool)
}

// RegisterBuiltin registers a built-in function.  This is used to register
// our "standard library" functions.
func RegisterBuiltin(name string, fun object.BuiltinFunction) {
	builtins[name] = &object.Builtin{Fn: fun}
}
