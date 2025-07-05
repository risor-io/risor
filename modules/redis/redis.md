# redis

This module provides a client to interact with a Redis server.

Module `redis` provides a set of functions to perform common Redis operations.

## Functions

### client

```go filename="Function signature"
client(dns string)
```

Instantiates a new Redis client for the given Redis server address(DSN).

```go filename="Example"
rdb := redis.client("redis://localhost:6379")
```

### ping

```go filename="Function signature"
ping() string
```

Sends a `PING` command to the Redis server and returns the response.

```go filename="Example"
rdb.ping() // "PONG"
```

### get

```go filename="Function signature"
get(key string) string
```

Retrieves the value of the given key from Redis.

```go filename="Example"
rdb.get("my_key") // "my_value"
```

### set

```go filename="Function signature"
set(key string, value string, expiration int) string
```

Sets the value of the given key with an optional expiration time (in seconds).

```go filename="Example"
rdb.set("my_key", "my_value", 300) // "OK"
```

### del

```go filename="Function signature"
del(key string) int
```

Deletes the given key from Redis and returns the number of keys removed.

```go filename="Example"
rdb.del("my_key") // 1
```

### exists

```go filename="Function signature"
exists(key string) int
```

Checks if the given key exists in Redis.

```go filename="Example"
rdb.exists("my_key") // 1
```

### keys

```go filename="Function signature"
keys(pattern string) []string
```

Returns a list of keys matching the given pattern.

```go filename="Example"
rdb.keys("*") // ["key1", "key2"]
```

### expire

```go filename="Function signature"
expire(key string, seconds int) bool
```

Sets an expiration time (in seconds) for the given key.

```go filename="Example"
rdb.expire("my_key", 300) // true
```

### ttl

```go filename="Function signature"
ttl(key string) int
```

Returns the remaining time-to-live (TTL) of the given key in seconds.

```go filename="Example"
rdb.ttl("my_key") // 300
```

### incr

```go filename="Function signature"
incr(key string) int
```

Increments the value of the given key by 1.

```go filename="Example"
rdb.incr("counter") // 1
```

### decr

```go filename="Function signature"
decr(key string) int
```

Decrements the value of the given key by 1.

```go filename="Example"
rdb.decr("counter") // 0
```

### flushdb

```go filename="Function signature"
flushdb() string
```

Deletes all keys in the current database.

```go filename="Example"
rdb.flushdb() // "OK"
```
```