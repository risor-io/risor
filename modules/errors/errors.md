# errors

Module `errors` provides functions for creating and manipulating error values.

## Functions

### new

```go filename="Function signature"
new(string) error
```

Returns a new error with the given message.

```go filename="Example"
>>> err := errors.new("something went wrong")
>>> err
something went wrong
```
