//go:build google
// +build google

package google

import (
	"context"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func ClientFunc(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("aws.client", 2, args); err != nil {
		return err
	}
	serviceName, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	resourceName, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return getClient(ctx, serviceName, resourceName)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("google", map[string]object.Object{
		"client": object.NewBuiltin("google.client", ClientFunc),
	})
}
