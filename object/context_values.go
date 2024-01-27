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

// SpawnFunc is a type signature for a function that can spawn a Risor thread.
type SpawnFunc func(ctx context.Context, fn Callable, args []Object) (*Thread, error)

const spawnFuncKey = contextKey("risor:spawn")

// WithSpawnFunc adds an SpawnFunc to the context, which can be used by
// objects to spawn themselves.
func WithSpawnFunc(ctx context.Context, fn SpawnFunc) context.Context {
	return context.WithValue(ctx, spawnFuncKey, fn)
}

// GetSpawnFunc returns the SpawnFunc from the context, if it exists.
func GetSpawnFunc(ctx context.Context) (SpawnFunc, bool) {
	fn, ok := ctx.Value(spawnFuncKey).(SpawnFunc)
	return fn, ok
}

////////////////////////////////////////////////////////////////////////////////

const threadKey = contextKey("risor:thread")

// WithThread returns a context with a Thread associated.
func WithThread(ctx context.Context, t *Thread) context.Context {
	return context.WithValue(ctx, threadKey, t)
}

// GetThread returns the thread associated with the context, if it exists.
func GetThread(ctx context.Context) (*Thread, bool) {
	t, ok := ctx.Value(threadKey).(*Thread)
	return t, ok
}
