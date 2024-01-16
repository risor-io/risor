//go:build vault
// +build vault

package vault

import (
	"context"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func NewVault(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("vault.connect", 1, args); err != nil {
		return err
	}

	addr, objErr := object.AsString(args[0])
	if objErr != nil {
		return objErr
	}

	vault, err := New(addr)
	if err != nil {
		return object.NewError(err)
	}

	return vault
}

func Module() *object.Module {
	return object.NewBuiltinsModule("vault", map[string]object.Object{
		"connect": object.NewBuiltin("vault.connect", NewVault),
	})
}
