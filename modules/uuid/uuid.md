# uuid

Module `uuid` provides generation of the different flavors of UUIDs. The
generated UUIDs are returned as strings.

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
