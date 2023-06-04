package uuid

import (
	"context"

	"github.com/cloudcmds/tamarin/v2/arg"
	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/gofrs/uuid"
)

const Name = "uuid"

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
	if err != nil {
		return object.Errorf(nsErr.Error())
	}
	return object.NewString(uuid.NewV5(nsID, name).String())
}

func Module() *object.Module {
	m := object.NewBuiltinsModule(Name, map[string]object.Object{
		"v4": object.NewBuiltin("v4", V4),
		"v5": object.NewBuiltin("v5", V5),
	})
	return m
}
