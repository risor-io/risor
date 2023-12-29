# json

Module `json` provides JSON encoding and decoding.

## Functions

### marshal

```go filename="Function signature"
marshal(v object) string
```

Returns a JSON string representing the given value. Raises an error if the value
cannot be marshalled.

```go copy filename="Example"
>>> m := {one: 1, two: 2}
>>> json.marshal(m)
"{\"one\":1,\"two\":2}"
```

### unmarshal

```go filename="Function signature"
unmarshal(s string) object
```

Returns the value represented by the given JSON string. Raises an error if the
string cannot be unmarshalled.

```go copy filename="Example"
>>> json.unmarshal("{\"one\":1,\"two\":2}")
{"one": 1, "two": 2}
>>> json.unmarshal("{bad") // raises value error
```

### valid

```go filename="Function signature"
valid(s string) bool
```

Returns whether the given string is valid JSON.

```go copy filename="Example"
>>> json.valid("42")
true
>>> json.valid("{oops")
false
```
