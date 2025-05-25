package redis

import (
	"context"
	"github.com/risor-io/risor/object"
	"time"
)

func (r *Client) Ping(ctx context.Context, args ...object.Object) object.Object {
	if len(args) > 0 {
		return object.TypeErrorf("type error: redis.client.ping() does not accept any arguments (%d given)", len(args))
	}

	result := r.client.Ping(ctx)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewString(result.Val())
}

func (r *Client) Get(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.get() requires one argument")
	}
	if len(args) > 1 {
		return object.TypeErrorf("type error: redis.client.get() requires one argument (%d given)", len(args))
	}

	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	result := r.client.Get(ctx, key)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewString(result.Val())
}

func (r *Client) Set(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.set() requires at least two arguments")
	}
	if len(args) < 2 {
		return object.TypeErrorf("type error: redis.client.set() requires at least two arguments (%d given)", len(args))
	}
	if len(args) > 3 {
		return object.TypeErrorf("type error: redis.client.set() requires at most three arguments (%d given)", len(args))
	}

	var expiration int64
	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	value, err := object.AsString(args[1])
	if err != nil {
		return err
	}

	if len(args) == 3 {
		expiration, err = object.AsInt(args[2])
		if err != nil {
			return err
		}
	}

	result := r.client.Set(ctx, key, value, time.Duration(expiration)*time.Second)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewString(result.Val())
}
func (r *Client) Del(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.del() requires at least one argument")
	}
	if len(args) > 1 {
		return object.TypeErrorf("type error: redis.client.del() requires at most one argument (%d given)", len(args))
	}

	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	result := r.client.Del(ctx, key)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewInt(result.Val())
}

func (r *Client) Exists(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.exists() requires at least one argument")
	}
	if len(args) > 1 {
		return object.TypeErrorf("type error: redis.client.exists() requires at most one argument (%d given)", len(args))
	}

	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	result := r.client.Exists(ctx, key)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewInt(result.Val())
}

func (r *Client) Keys(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.keys() requires one argument")
	}
	if len(args) > 1 {
		return object.TypeErrorf("type error: redis.client.keys() requires at most one argument (%d given)", len(args))
	}

	pattern, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	result := r.client.Keys(ctx, pattern)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewStringList(result.Val())
}

func (r *Client) Expire(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.expire() requires at least two arguments")
	}
	if len(args) < 2 {
		return object.TypeErrorf("type error: redis.client.expire() requires at least two arguments (%d given)", len(args))
	}
	if len(args) >= 3 {
		return object.TypeErrorf("type error: redis.client.expire() requires at most two arguments (%d given)", len(args))
	}

	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	expiration, err := object.AsInt(args[1])
	if err != nil {
		return err
	}

	result := r.client.Expire(ctx, key, time.Duration(expiration)*time.Second)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewBool(result.Val())
}

func (r *Client) TTL(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.ttl() requires one argument")
	}
	if len(args) > 1 {
		return object.TypeErrorf("type error: redis.client.ttl() requires one argument (%d given)", len(args))
	}

	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	result := r.client.TTL(ctx, key)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewInt(int64(result.Val().Seconds()))
}

func (r *Client) Incr(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.incr() requires one argument")
	}
	if len(args) > 1 {
		return object.TypeErrorf("type error: redis.client.incr() requires one argument (%d given)", len(args))
	}

	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	result := r.client.Incr(ctx, key)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewInt(result.Val())
}

func (r *Client) Decr(ctx context.Context, args ...object.Object) object.Object {
	if len(args) == 0 {
		return object.TypeErrorf("type error: redis.client.decr() requires one argument")
	}
	if len(args) > 1 {
		return object.TypeErrorf("type error: redis.client.decr() requires one argument (%d given)", len(args))
	}

	key, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	result := r.client.Decr(ctx, key)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewInt(result.Val())
}

func (r *Client) FlushDB(ctx context.Context, args ...object.Object) object.Object {
	if len(args) > 0 {
		return object.TypeErrorf("type error: redis.client.flushdb() does not accept any arguments (%d given)", len(args))
	}

	result := r.client.FlushDB(ctx)
	if result.Err() != nil {
		return object.NewError(result.Err())
	}
	return object.NewString(result.Val())
}
