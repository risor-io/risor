package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const REDIS object.Type = "redis.client"

type Client struct {
	client *redis.Client
}

func (r *Client) Type() object.Type {
	return REDIS
}

func (r *Client) Inspect() string {
	return "redis.client()"
}

func (r *Client) Interface() interface{} {
	return r.client
}

func (r *Client) Equals(other object.Object) object.Object {
	if r == other {
		return object.True
	}
	return object.False
}

func (r *Client) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "ping":
		return object.NewBuiltin("ping", r.Ping), true
	case "get":
		return object.NewBuiltin("get", r.Get), true
	case "set":
		return object.NewBuiltin("set", r.Set), true
	case "del":
		return object.NewBuiltin("del", r.Del), true
	case "exists":
		return object.NewBuiltin("exists", r.Exists), true
	case "keys":
		return object.NewBuiltin("keys", r.Keys), true
	case "expire":
		return object.NewBuiltin("expire", r.Expire), true
	case "ttl":
		return object.NewBuiltin("ttl", r.TTL), true
	case "incr":
		return object.NewBuiltin("incr", r.Incr), true
	case "decr":
		return object.NewBuiltin("decr", r.Decr), true
	case "flushdb":
		return object.NewBuiltin("flushdb", r.FlushDB), true
	}
	return nil, false
}

func (r *Client) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("type error: cannot set %q on %s object", name, REDIS)
}

func (r *Client) IsTruthy() bool {
	return r.client != nil
}

func (r *Client) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("type error: unsupported operation for %s object", REDIS)
}

func (r *Client) Cost() int {
	return 0
}

func New(client *redis.Client) *Client {
	return &Client{
		client: client,
	}
}
