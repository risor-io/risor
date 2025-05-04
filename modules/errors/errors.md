# errors

Module `errors` provides functions for creating error values.

## Comparison with Go

Risor error values are similar to those in Go, in that they are values that
represent an error condition. However, error handling is different in Risor
because it has the concept of raising and catching errors.

This approach has two main benefits in Risor:

1. It keeps Risor code more concise, which is desirable for a scripting language.
2. The fact that functions in Risor always return exactly one value means that
   function results can be piped without having to check for errors.

## Functions

### new

```go filename="Function signature"
new(string) error
```

Returns a new error value with the given message.

```go filename="Example"
>>> err := errors.new("something went wrong")
>>> err
something went wrong
```

### is

```go filename="Function signature"
is(err error, target error) bool
```

Reports whether error matches target error.

```go filename="Example"
>>> errors.is(err, os.err_not_exist)
true
```