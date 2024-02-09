package object

import (
	"context"
	"errors"
	"fmt"
)

type callFuncAdapter struct {
	funcObj *Function
}

func (c *callFuncAdapter) Call(ctx context.Context, args ...Object) Object {
	callFunc, found := GetCallFunc(ctx)
	if !found {
		return Errorf("eval error: context did not contain a call function")
	}
	result, err := callFunc(ctx, c.funcObj, args)
	if err != nil {
		return NewError(err)
	}
	return result
}

func Spawn(ctx context.Context, fnObj Object, args []Object) (*Thread, error) {
	spawnFunc, found := GetSpawnFunc(ctx)
	if !found {
		return nil, errors.New("eval error: context did not contain a spawn function")
	}

	// create independent copy of args to guarantee the caller can't modify
	// the slice provided to the spawned function
	argsCopy := make([]Object, len(args))
	copy(argsCopy, args)

	switch fn := fnObj.(type) {
	case *Function:
		adapter := &callFuncAdapter{funcObj: fn}
		obj, err := spawnFunc(ctx, adapter, argsCopy)
		if err != nil {
			return nil, err
		}
		return obj, nil
	case Callable: // *Builtin is Callable
		obj, err := spawnFunc(ctx, fn, argsCopy)
		if err != nil {
			return nil, err
		}
		return obj, nil
	default:
		return nil, fmt.Errorf("type error: spawn() expected a function (%s given)", fn.Type())
	}
}
