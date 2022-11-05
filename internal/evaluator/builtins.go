package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gofrs/uuid"
	"github.com/cloudcmds/tamarin/object"
)

// The built-in functions / standard-library methods are stored here.
var builtins = map[string]*object.Builtin{}

// convert a string, boolean, or float to an int
func intFun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch args[0].(type) {
	case *object.String:
		input := args[0].(*object.String).Value
		i, err := strconv.Atoi(input)
		if err == nil {
			return &object.Integer{Value: int64(i)}
		}
		return newError("Converting string '%s' to int failed %s", input, err.Error())
	case *object.Boolean:
		input := args[0].(*object.Boolean).Value
		if input {
			return &object.Integer{Value: 1}
		}
		return &object.Integer{Value: 0}
	case *object.Integer:
		return args[0]
	case *object.Float:
		input := args[0].(*object.Float).Value
		return &object.Integer{Value: int64(input)}
	default:
		return newError("argument to `int` not supported, got=%s", args[0].Type())
	}
}

// length of item
func lenFun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(utf8.RuneCountInString(arg.Value))}
	case *object.Null:
		return &object.Integer{Value: 0}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	case *object.Set:
		return &object.Integer{Value: int64(len(arg.Items))}
	default:
		return newError("argument to `len` not supported, got=%s", args[0].Type())
	}
}

// regular expression match
func matchFun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}
	if args[0].Type() != object.STRING_OBJ {
		return newError("argument to `match` must be STRING, got %s", args[0].Type())
	}
	if args[1].Type() != object.STRING_OBJ {
		return newError("argument to `match` must be STRING, got %s", args[1].Type())
	}
	reg := regexp.MustCompile(args[0].(*object.String).Value)
	res := reg.FindStringSubmatch(args[1].(*object.String).Value)
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
				k := &object.Integer{Value: int64(i - 1)}
				v := &object.String{Value: res[i]}
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
		return &object.Null{}
	}
	if args[0].Type() != object.STRING_OBJ {
		return &object.Null{}
	}
	fs := args[0].(*object.String).Value
	// Convert the arguments to something go's sprintf code will understand
	argLen := len(args)
	fmtArgs := make([]interface{}, argLen-1)
	for i, v := range args[1:] {
		fmtArgs[i] = v.ToInterface()
	}
	return &object.String{Value: fmt.Sprintf(fs, fmtArgs...)}
}

// Get hash keys
func hashKeys(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("argument to `keys` must be HASH, got=%s", args[0].Type())
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
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2", len(args))
	}
	if args[0].Type() != object.HASH_OBJ {
		return newError("argument to `delete` must be HASH, got=%s", args[0].Type())
	}
	// The object we're working with
	hash := args[0].(*object.Hash)
	// The key we're going to delete
	key, ok := args[1].(object.Hashable)
	if !ok {
		return newError("key `delete` into HASH must be Hashable, got=%s", args[1].Type())
	}
	// Make a new hash
	newHash := make(map[object.HashKey]object.HashPair)
	// Copy the values EXCEPT the one we have.
	for k, v := range hash.Pairs {
		if k != key.HashKey() {
			newHash[k] = v
		}
	}
	return &object.Hash{Pairs: newHash}
}

func setFun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	set := &object.Set{Items: map[object.HashKey]object.Object{}}
	arg := args[0]
	switch arg := arg.(type) {
	case *object.String:
		for _, v := range arg.Value {
			vStr := &object.String{Value: string(v)}
			set.Items[vStr.HashKey()] = vStr
		}
	case *object.Array:
		for _, obj := range arg.Elements {
			hashable, ok := obj.(object.Hashable)
			if !ok {
				return newError("type error: object is not hashable: %v", obj.Type())
			}
			set.Items[hashable.HashKey()] = obj
		}
	}
	return set
}

func strFun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	out := args[0].Inspect()
	return &object.String{Value: out}
}

func typeFun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	switch args[0].(type) {
	case *object.String:
		return &object.String{Value: "string"}
	case *object.Regexp:
		return &object.String{Value: "regexp"}
	case *object.Boolean:
		return &object.String{Value: "bool"}
	case *object.Builtin:
		return &object.String{Value: "builtin"}
	case *object.Array:
		return &object.String{Value: "array"}
	case *object.Function:
		return &object.String{Value: "function"}
	case *object.Integer:
		return &object.String{Value: "integer"}
	case *object.Float:
		return &object.String{Value: "float"}
	case *object.Hash:
		return &object.String{Value: "hash"}
	case *object.Set:
		return &object.String{Value: "set"}
	case *object.Result:
		return &object.String{Value: "result"}
	case *object.HttpResponse:
		return &object.String{Value: "http_response"}
	default:
		return newError("argument to `type` not supported, got=%s", args[0].Type())
	}
}

func randomFun(ctx context.Context, args ...object.Object) object.Object {
	return &object.Float{Value: rand.Float64()}
}

func okFun(ctx context.Context, args ...object.Object) object.Object {
	switch len(args) {
	case 0:
		return &object.Result{Ok: object.NULL}
	case 1:
		return &object.Result{Ok: args[0]}
	default:
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
}

func errFun(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}
	return &object.Result{Err: args[0]}
}

func assertFun(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
	}
	if !isTruthy(args[0]) {
		if numArgs == 2 {
			return newError(args[1].Inspect())
		}
		return newError("assertion failed")
	}
	return object.NULL
}

func fetchFun(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return newError("wrong number of arguments. got=%d, want=1 or 2", len(args))
	}
	urlArg, ok := args[0].(*object.String)
	if !ok {
		return newError("type error: expected a string argument; got %v", args[0].Type())
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
	req, err := http.NewRequestWithContext(ctx, method, urlArg.Value, body)
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

func uuidFun(ctx context.Context, args ...object.Object) object.Object {
	value, err := uuid.NewV4()
	if err != nil {
		return newError("failed to generate uuid: %v", err)
	}
	return &object.String{Value: value.String()}
}

func init() {
	RegisterBuiltin("delete", hashDelete)
	RegisterBuiltin("int", intFun)
	RegisterBuiltin("keys", hashKeys)
	RegisterBuiltin("len", lenFun)
	RegisterBuiltin("match", matchFun)
	RegisterBuiltin("set", setFun)
	RegisterBuiltin("sprintf", sprintfFun)
	RegisterBuiltin("str", strFun)
	RegisterBuiltin("type", typeFun)
	RegisterBuiltin("random", randomFun)
	RegisterBuiltin("ok", okFun)
	RegisterBuiltin("err", errFun)
	// RegisterBuiltin("json_unmarshal", object.JsonUnmarshal)
	RegisterBuiltin("assert", assertFun)
	RegisterBuiltin("fetch", fetchFun)
	RegisterBuiltin("uuid", uuidFun)

	// TODO:
	// any
	// all
	// bool
	// chr
	// ord
	// filter
	// float
	// hex
	// map
	// oct
	// pow
	// round
	// sorted
	// sum
	// isinstance
}

// RegisterBuiltin registers a built-in function.  This is used to register
// our "standard library" functions.
func RegisterBuiltin(name string, fun object.BuiltinFunction) {
	builtins[name] = &object.Builtin{Fn: fun}
}
