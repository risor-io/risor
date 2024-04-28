# uuid

Module `uuid` provides generation of the different flavors of UUIDs. The
generated UUIDs are returned as strings.

The core functionality is provided by
[github.com/gofrs/uuid](https://github.com/gofrs/uuid).

## Module

```go copy filename="Function signature"
uuid() string
```

The `uuid` module object itself is callable, which is a shorthand for `uuid.v4()`.

```go copy filename="Example"
>>> uuid()
"83650166-58e3-4077-91a1-f176199f4954"
```

## Functions

### v4

```go filename="Function signature"
v4() string
```

Returns a randomly generated v4 UUID.

```go filename="Example"
>>> uuid.v4()
"83650166-58e3-4077-91a1-f176199f4954"
```

### v5

```go filename="Function signature"
v5(namespace, name string) string
```

Returns a UUID based on SHA-1 hash of the namespace UUID and name.

```go filename="Example"
>>> uuid.v5("c54e4fe3-ced2-4d07-a373-a95da134685e", "test")
"cd0c3ed6-1143-5e9f-bdc7-dd32215bb7ea"
```

### v6

```go filename="Function signature"
v6() string
```

Returns a v6 UUID. The v6 UUID is a reordering of UUIDv1 fields so it is
lexicographically sortable by time.

```go filename="Example"
>>> uuid.v6()
"1ef0171d-c5fb-6d5c-955f-d006d0ab5f0c"
```

### v7

```go filename="Function signature"
v7() string
```

Returns a v7 UUID. The v7 UUID is time-ordered and embeds a Unix timestamp
with millisecond precision. The time-ordered aspect makes these IDs useful
in some database scenarios, since database performance may be improved as
compared to v4 UUIDs.

```go filename="Example"
>>> uuid.v7()
"018f0b0e-1841-7560-9ec9-fa3272646ac7"
```
