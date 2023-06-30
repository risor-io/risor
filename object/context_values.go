package object

import (
	"context"
)

type contextKey string

////////////////////////////////////////////////////////////////////////////////
// Store and retrieve a function that can call a compiled Risor function
////////////////////////////////////////////////////////////////////////////////

// CallFunc is a type signature for a function that can call a Risor function.
type CallFunc func(ctx context.Context, fn *Function, args []Object) (Object, error)

const callFuncKey = contextKey("risor:call")

// WithCallFunc adds an CallFunc to the context, which can be used by
// objects to call a Risor function at runtime.
func WithCallFunc(ctx context.Context, fn CallFunc) context.Context {
	return context.WithValue(ctx, callFuncKey, fn)
}

// GetCallFunc returns the CallFunc from the context, if it exists.
func GetCallFunc(ctx context.Context) (CallFunc, bool) {
	fn, ok := ctx.Value(callFuncKey).(CallFunc)
	return fn, ok
}

////////////////////////////////////////////////////////////////////////////////
// Store and retrieve a function that can retrieve the active code
////////////////////////////////////////////////////////////////////////////////

// CodeFunc is a type signature for a function that can retrieve the active code.
type CodeFunc func(ctx context.Context) (*Code, error)

const codeFuncKey = contextKey("risor:code")

// WithCodeFunc adds an CodeFunc to the context, which can be used by
// objects to retrieve the active code at runtime
func WithCodeFunc(ctx context.Context, fn CodeFunc) context.Context {
	return context.WithValue(ctx, codeFuncKey, fn)
}

// GetCodeFunc returns the CodeFunc from the context, if it exists.
func GetCodeFunc(ctx context.Context) (CodeFunc, bool) {
	fn, ok := ctx.Value(codeFuncKey).(CodeFunc)
	return fn, ok
}
