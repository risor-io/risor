package spawn

import (
	"context"
	"errors"
	"fmt"

	"github.com/risor-io/risor/object"
)

type callFuncAdapter struct {
	funcObj *object.Function
}

func (c *callFuncAdapter) Call(ctx context.Context, args ...object.Object) object.Object {
	callFunc, found := object.GetCallFunc(ctx)
	if !found {
		return object.Errorf("eval error: context did not contain a call function")
	}
	result, err := callFunc(ctx, c.funcObj, args)
	if err != nil {
		return object.NewError(err)
	}
	return result
}

func Spawn(ctx context.Context, fnObj object.Object, args []object.Object) (*object.Thread, error) {
	spawnFunc, found := object.GetSpawnFunc(ctx)
	if !found {
		return nil, errors.New("eval error: context did not contain a spawn function")
	}

	// create independent copy of args to guarantee the caller can't modify
	// the slice provided to the spawned function
	argsCopy := make([]object.Object, len(args))
	copy(argsCopy, args)

	switch fn := fnObj.(type) {
	case *object.Function:
		adapter := &callFuncAdapter{funcObj: fn}
		obj, err := spawnFunc(ctx, adapter, argsCopy)
		if err != nil {
			return nil, err
		}
		return obj, nil
	case object.Callable: // *object.Builtin is Callable
		obj, err := spawnFunc(ctx, fn, argsCopy)
		if err != nil {
			return nil, err
		}
		return obj, nil
	default:
		return nil, fmt.Errorf("type error: spawn() expected a function (%s given)", fn.Type())
	}
}
