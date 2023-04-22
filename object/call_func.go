package object

import (
	"context"
)

// CallFunc is a type signature for a function that can call a Tamarin function.
type CallFunc func(ctx context.Context, fn *Function, args []Object) (Object, error)

type contextKey string

const callFuncKey = contextKey("evaluator")

// WithCallFunc adds an CallFunc to the context, which can be used by
// objects to call a Tamarin function at runtime.
func WithCallFunc(ctx context.Context, fn CallFunc) context.Context {
	return context.WithValue(ctx, callFuncKey, fn)
}

// GetCallFunc returns the CallFunc from the context, if it exists.
func GetCallFunc(ctx context.Context) (CallFunc, bool) {
	fn, ok := ctx.Value(callFuncKey).(CallFunc)
	return fn, ok
}
