# strconv

String conversion functions from the Go standard library.

## Functions

### atoi

```go filename="Function signature"
atoi(s string) int
```

Converts the string s to an int.

```go copy filename="Example"
>>> strconv.atoi("123")
123
>>> strconv.atoi("nope")
strconv.Atoi: parsing "nope": invalid syntax
```

### parse_bool

```go filename="Function signature"
parse_bool(s string) bool
```

Converts the string s to a bool.

```go copy filename="Example"
>>> strconv.parse_bool("true")
true
>>> strconv.parse_bool("false")
false
>>> strconv.parse_bool("nope")
strconv.ParseBool: parsing "nope": invalid syntax
```

### parse_float

```go filename="Function signature"
parse_float(s string) float
```

Converts the string s to a float.

```go copy filename="Example"
>>> strconv.parse_float("3.14")
3.14
>>> strconv.parse_float("nope")
strconv.ParseFloat: parsing "nope": invalid syntax
```

### parse_int

```go filename="Function signature"
parse_int(s string, base int = 10, bit_size int = 64) int
```

Converts the string s to an int. 

```go copy filename="Example"
>>> strconv.parse_int("123")
123
>>> strconv.parse_int("nope")
strconv.ParseInt: parsing "nope": invalid syntax
>>> strconv.parse_int("ff", 16)
255
```
