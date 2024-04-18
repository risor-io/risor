package bcrypt

import (
	"context"

	"github.com/risor-io/risor/object"
	"golang.org/x/crypto/bcrypt"
)

func Hash(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsRangeError("bcrypt.hash", 1, 2, numArgs)
	}
	password, argErr := object.AsBytes(args[0])
	if argErr != nil {
		return argErr
	}
	cost := bcrypt.DefaultCost
	if numArgs > 1 {
		cost64, argErr := object.AsInt(args[1])
		if argErr != nil {
			return argErr
		}
		cost = int(cost64)
	}
	hash, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return object.NewError(err)
	}
	return object.NewByteSlice(hash)
}

func Compare(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs != 2 {
		return object.NewArgsError("bcrypt.compare", 2, numArgs)
	}
	hash, argErr := object.AsBytes(args[0])
	if argErr != nil {
		return argErr
	}
	password, argErr := object.AsBytes(args[1])
	if argErr != nil {
		return argErr
	}
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return object.NewError(err)
	}
	return object.NewBool(true)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("bcrypt", map[string]object.Object{
		"hash":         object.NewBuiltin("hash", Hash),
		"compare":      object.NewBuiltin("compare", Compare),
		"min_cost":     object.NewInt(int64(bcrypt.MinCost)),
		"max_cost":     object.NewInt(int64(bcrypt.MaxCost)),
		"default_cost": object.NewInt(int64(bcrypt.DefaultCost)),
	})
}
