package http

import (
	"context"
	"net/http"

	"github.com/risor-io/risor/object"
)

func NewHttpRequest(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsRangeError("fetch", 1, 2, numArgs)
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
	req, errObj := NewRequestFromParams(urlArg, params)
	if errObj != nil {
		return errObj
	}

	return req
}

func MethodCmd(method string) object.BuiltinFunction {
	return func(ctx context.Context, args ...object.Object) object.Object {
		numArgs := len(args)
		if numArgs < 1 || numArgs > 3 {
			return object.NewArgsRangeError("fetch", 1, 3, numArgs)
		}
		urlArg, argErr := object.AsString(args[0])
		if argErr != nil {
			return argErr
		}

		var errObj *object.Error

		params := object.NewMap(map[string]object.Object{
			"method": object.NewString(method),
		})

		if numArgs > 1 {
			headers, errObj := object.AsMap(args[1])
			if errObj != nil {
				return errObj
			}
			params.Set("headers", headers)
		}

		if numArgs > 2 {
			var key string
			switch method {
			case http.MethodGet, http.MethodHead, http.MethodDelete:
				key = "params"
			case http.MethodPost, http.MethodPut, http.MethodPatch:
				key = "body"
			}
			params.Set(key, args[2])
		}

		req, errObj := NewRequestFromParams(urlArg, params)
		if errObj != nil {
			return errObj
		}

		return req
	}
}

func Handle(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 2 {
		return object.NewArgsRangeError("http.handle", 2, 3, numArgs)
	}
	pattern, errObj := object.AsString(args[0])
	if errObj != nil {
		return errObj
	}
	callFn, ok := object.GetCloneCallFunc(ctx)
	if !ok {
		return object.Errorf("http.handle: no clone-call function found in context")
	}
	var handler http.Handler
	switch fn := args[1].(type) {
	case http.Handler:
		handler = fn
	case *object.Function:
		handler = HandlerFunc(fn, callFn)
	default:
		return object.Errorf("type error: unsupported http handler type: %s", fn.Type())
	}
	http.Handle(pattern, handler)
	return object.Nil
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"fetch": object.NewBuiltin("fetch", Fetch),
	}
}

type ModuleOpts struct {
	ListenersAllowed bool
}

func Module(opts ...ModuleOpts) *object.Module {
	var listenersAllowed bool
	if len(opts) > 0 {
		listenersAllowed = opts[0].ListenersAllowed
	}
	builtins := map[string]object.Object{
		"delete":  object.NewBuiltin("http.delete", MethodCmd(http.MethodDelete)),
		"get":     object.NewBuiltin("http.get", MethodCmd(http.MethodGet)),
		"head":    object.NewBuiltin("http.head", MethodCmd(http.MethodHead)),
		"patch":   object.NewBuiltin("http.patch", MethodCmd(http.MethodPatch)),
		"post":    object.NewBuiltin("http.post", MethodCmd(http.MethodPost)),
		"put":     object.NewBuiltin("http.put", MethodCmd(http.MethodPut)),
		"request": object.NewBuiltin("http.request", NewHttpRequest),
	}
	if listenersAllowed {
		builtins["listen_and_serve"] = object.NewBuiltin("http.listen_and_serve", ListenAndServe)
		builtins["listen_and_serve_tls"] = object.NewBuiltin("http.listen_and_serve_tls", ListenAndServeTLS)
		builtins["handle"] = object.NewBuiltin("http.handle", Handle)
	}
	return object.NewBuiltinsModule("http", builtins)
}
