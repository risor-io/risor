package uuid

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
)

func V4(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("uuid.v4", 0, args); err != nil {
		return err
	}
	value, err := uuid.NewV4()
	if err != nil {
		return object.Errorf(err.Error())
	}
	return object.NewString(value.String())
}

func V5(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("uuid.v5", 2, args); err != nil {
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
	if nsErr != nil {
		return object.Errorf(nsErr.Error())
	}
	return object.NewString(uuid.NewV5(nsID, name).String())
}

func V6(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("uuid.v6", 0, args); err != nil {
		return err
	}
	value, err := uuid.NewV6()
	if err != nil {
		return object.Errorf(err.Error())
	}
	return object.NewString(value.String())
}

func V7(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("uuid.v7", 0, args); err != nil {
		return err
	}
	value, err := uuid.NewV7()
	if err != nil {
		return object.Errorf(err.Error())
	}
	return object.NewString(value.String())
}

func Module() *object.Module {
	return object.NewBuiltinsModule("uuid", map[string]object.Object{
		"v4": object.NewBuiltin("v4", V4),
		"v5": object.NewBuiltin("v5", V5),
		"v6": object.NewBuiltin("v6", V6),
		"v7": object.NewBuiltin("v7", V7),
	}, V4)
}
