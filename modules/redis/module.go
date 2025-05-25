package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
)

func Create(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("redis.client", 1, 1, args); err != nil {
		return err
	}
	redisURL, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	opts, err2 := redis.ParseURL(redisURL)
	if err2 != nil {
		return object.NewError(err2)
	}

	client := redis.NewClient(opts)
	return New(client)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("redis",
		map[string]object.Object{
			"client": object.NewBuiltin("client", Create),
		},
		Create,
	)
}
