package object

import (
	"context"
)

type contextKey string

// CallFunc is a type signature for a function that can call a Risor function.
type CallFunc func(ctx context.Context, fn *Function, args []Object) (Object, error)

// SpawnFunc is a type signature for a function that can spawn a Risor thread.
type SpawnFunc func(ctx context.Context, fn Callable, args []Object) (*Thread, error)

type CallFrame struct {
	Depth    int    `json:"depth"`
	Function string `json:"function"`
}

type CallStack struct {
	Frames []*CallFrame `json:"frames"`
}

// StackFunc is a type signature for a function that can get the current stack.
type StackFunc func(ctx context.Context) (*CallStack, error)

////////////////////////////////////////////////////////////////////////////////

const callFuncKey = contextKey("risor:call")

// WithCallFunc adds an CallFunc to the context, which can be used by
// objects to call a Risor function at runtime.
func WithCallFunc(ctx context.Context, fn CallFunc) context.Context {
	return context.WithValue(ctx, callFuncKey, fn)
}

// GetCallFunc returns the CallFunc from the context, if it exists.
func GetCallFunc(ctx context.Context) (CallFunc, bool) {
	if fn, ok := ctx.Value(callFuncKey).(CallFunc); ok {
		if fn != nil {
			return fn, ok
		}
	}
	return nil, false
}

////////////////////////////////////////////////////////////////////////////////

const spawnFuncKey = contextKey("risor:spawn")

// WithSpawnFunc adds an SpawnFunc to the context, which can be used by
// objects to spawn themselves.
func WithSpawnFunc(ctx context.Context, fn SpawnFunc) context.Context {
	return context.WithValue(ctx, spawnFuncKey, fn)
}

// GetSpawnFunc returns the SpawnFunc from the context, if it exists.
func GetSpawnFunc(ctx context.Context) (SpawnFunc, bool) {
	if fn, ok := ctx.Value(spawnFuncKey).(SpawnFunc); ok {
		if fn != nil {
			return fn, ok
		}
	}
	return nil, false
}

////////////////////////////////////////////////////////////////////////////////

const cloneCallKey = contextKey("risor:clone-call")

// WithCloneCallFunc returns a context with a "clone-call" function
// associated. This function can be used to clone a Risor VM and then call a
// function on it synchronously.
func WithCloneCallFunc(ctx context.Context, fn CallFunc) context.Context {
	return context.WithValue(ctx, cloneCallKey, fn)
}

// GetCloneCallFunc returns the "clone-call" function from the context,
// if it exists. This function can be used to clone a Risor VM and then call a
// function on it synchronously.
func GetCloneCallFunc(ctx context.Context) (CallFunc, bool) {
	if fn, ok := ctx.Value(cloneCallKey).(CallFunc); ok {
		if fn != nil {
			return fn, ok
		}
	}
	return nil, false
}

////////////////////////////////////////////////////////////////////////////////

const stackFuncKey = contextKey("risor:stack-func")

// WithStackFunc adds an StackFunc to the context, which can be used by
// objects to get the current stack.
func WithStackFunc(ctx context.Context, fn StackFunc) context.Context {
	return context.WithValue(ctx, stackFuncKey, fn)
}

// GetStackFunc returns the StackFunc from the context, if it exists.
func GetStackFunc(ctx context.Context) (StackFunc, bool) {
	if fn, ok := ctx.Value(stackFuncKey).(StackFunc); ok {
		if fn != nil {
			return fn, ok
		}
	}
	return nil, false
}
