package context

import (
	"context"
	"time"

	"github.com/risor-io/risor/object"
)

func Background(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 0 {
		return object.TypeErrorf("type error: context.background() takes no arguments (%d given)", len(args))
	}
	// Don't actually use context.Background(), since we need it bound for control
	return object.NewContext(ctx, nil)
}

func WithCancel(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.TypeErrorf("type error: context.with_cancel() takes 1 argument (%d given)", len(args))
	}
	parent, ok := args[0].(*object.Context)
	if !ok {
		return object.TypeErrorf("type error: context.with_cancel() argument must be a context")
	}
	newCtx, cancel := context.WithCancel(parent.Value())
	return object.NewContext(newCtx, cancel)
}

func WithTimeout(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.TypeErrorf("type error: context.with_timeout() takes 2 arguments (%d given)", len(args))
	}
	parent, ok := args[0].(*object.Context)
	if !ok {
		return object.TypeErrorf("type error: context.with_timeout() first argument must be a context")
	}
	timeout, ok := args[1].(*object.Float)
	if !ok {
		return object.TypeErrorf("type error: context.with_timeout() second argument must be a float")
	}
	newCtx, cancel := context.WithTimeout(parent.Value(), time.Duration(timeout.Value()*float64(time.Second)))
	return object.NewContext(newCtx, cancel)
}

func WithDeadline(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 2 {
		return object.TypeErrorf("type error: context.with_deadline() takes 2 arguments (%d given)", len(args))
	}
	parent, ok := args[0].(*object.Context)
	if !ok {
		return object.TypeErrorf("type error: context.with_deadline() first argument must be a context")
	}
	deadline, ok := args[1].(*object.Time)
	if !ok {
		return object.TypeErrorf("type error: context.with_deadline() second argument must be a time")
	}
	newCtx, cancel := context.WithDeadline(parent.Value(), deadline.Value())
	return object.NewContext(newCtx, cancel)
}

func WithValue(ctx context.Context, args ...object.Object) object.Object {
	if len(args) != 3 {
		return object.TypeErrorf("type error: context.with_value() takes 3 arguments (%d given)", len(args))
	}
	parent, ok := args[0].(*object.Context)
	if !ok {
		return object.TypeErrorf("type error: context.with_value() first argument must be a context")
	}
	key, ok := args[1].(*object.String)
	if !ok {
		return object.TypeErrorf("type error: context.with_value() second argument must be a string")
	}
	newCtx := context.WithValue(parent.Value(), key.Value(), args[2])
	return object.NewContext(newCtx, nil) //  parent.Cancel())
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"background":    object.NewBuiltin("background", Background),
		"with_cancel":   object.NewBuiltin("with_cancel", WithCancel),
		"with_timeout":  object.NewBuiltin("with_timeout", WithTimeout),
		"with_deadline": object.NewBuiltin("with_deadline", WithDeadline),
		"with_value":    object.NewBuiltin("with_value", WithValue),
	}
}

func Module() *object.Module {
	return object.NewBuiltinsModule("context", map[string]object.Object{
		"background":    object.NewBuiltin("background", Background),
		"with_cancel":   object.NewBuiltin("with_cancel", WithCancel),
		"with_timeout":  object.NewBuiltin("with_timeout", WithTimeout),
		"with_deadline": object.NewBuiltin("with_deadline", WithDeadline),
		"with_value":    object.NewBuiltin("with_value", WithValue),
	})
}
