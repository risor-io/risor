# regexp

Module `regexp` provides regular expression matching.

More info here.

## Functions

### compile

```go filename="Function signature"
compile(expr string) regexp
```

Compiles a regular expression string into a regexp object.

```go copy filename="Example"
>>> regexp.compile("a+")
regexp("a+")
>>> r := regexp.compile("a+"); r.match("a")
true
>>> r := regexp.compile("[0-9]+"); r.match("nope")
false
```

### match

```go filename="Function signature"
match(expr, s string) bool
```

Returns true if the string s contains any match of the regular expression pattern.

```go copy filename="Example"
>>> regexp.match("ab+a", "abba")
true
>>> regexp.match("[0-9]+", "nope")
false
```

Another line.

Great info.

## Go Compatibility

The supported regular expression syntax is exactly as described
in the [Go regexp](https://pkg.go.dev/regexp) documentation.
