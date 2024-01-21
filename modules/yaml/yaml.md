# yaml

Module `yaml` provides YAML encoding and decoding.

## Functions

### marshal

```go filename="Function signature"
marshal(v object) string
```

Returns a YAML string representing the given value. Raises an error if the value
cannot be marshalled.

```go copy filename="Example"
>>> m := {one: 1, two: 2}
>>> yaml.marshal(m)
"one: 1\ntwo: 2\n"
```

### unmarshal

```go filename="Function signature"
unmarshal(s string) object
```

Returns the value represented by the given YAML string. Raises an error if the
string cannot be unmarshalled.

```go copy filename="Example"
>>> yaml.unmarshal("one: 1\ntwo: 2")
{"one": 1, "two": 2}
>>> yaml.unmarshal("{bad") // raises value error
```

### valid

```go filename="Function signature"
valid(s string) bool
```

Returns whether the given string is valid YAML.

```go copy filename="Example"
>>> yaml.valid("42")
true
>>> yaml.valid("{oops")
false
```
