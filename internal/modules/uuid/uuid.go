package uuid

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
	"github.com/gofrs/uuid"
)

const Name = "uuid"

func V4(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("uuid.v4", 0, args); err != nil {
		return err
	}
	value, err := uuid.NewV4()
	if err != nil {
		return object.NewError(err.Error())
	}
	return object.NewString(value.String())
}

func V5(ctx context.Context, args ...object.Object) object.Object {
	if err := RequireArgs("uuid.v5", 2, args); err != nil {
		return err
	}
	namespace, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	name, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	nsID, nsErr := uuid.FromString(namespace)
	if err != nil {
		return object.NewError(nsErr.Error())
	}
	return object.NewString(uuid.NewV5(nsID, name).String())
}

func RequireArgs(funcName string, count int, args []object.Object) *object.Error {
	nArgs := len(args)
	if nArgs != count {
		return object.NewError(
			fmt.Sprintf("type error: %s() takes exactly %d argument (%d given)",
				funcName, count, nArgs))
	}
	return nil
}

func Module(parentScope *scope.Scope) (*object.Module, error) {
	s := scope.New(scope.Opts{
		Name:   fmt.Sprintf("module:%s", Name),
		Parent: parentScope,
	})
	if err := s.AddBuiltins([]scope.Builtin{
		{Name: "v4", Func: V4},
		{Name: "v5", Func: V5},
	}); err != nil {
		return nil, err
	}
	return &object.Module{Name: Name, Scope: s}, nil
}
