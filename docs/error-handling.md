# Error Handling

There are two concepts used in Tamarin for dealing with situations where
something goes wrong: _errors_ and _results_. A Tamarin object type is used
for each of these, `object.Error` and `object.Result`. There are various
built-in functions and methods used to interact with these as values to
design an error resilient application.

In short, when an error object is generated, the Tamarin interpreter recognizes
the error, evaluation of the code stops, and the error is returned. This means
that, by default, a Tamarin program halts on an error.

Results are objects that may contain an _ok_ value or an _err_ value. These are
returned from operations that have common error conditions, such as I/O operations
that may fail due to network problems or similar. These result objects offer
methods to test if the result is _ok_ or _err_ and then the developer can decide
how the application should handle each of those cases.

The built-in function `try` can be used to work with both _errors_ and _results_.
It allows the developer to return a fallback value when an expression evaluates
to an error or an err result.

## Errors

Generate an error intentionally using the `error` built-in function. This stops
evaluation of the Tamarin program immediately:

```go
>>> error("kaboom")
kaboom
```

Wrap an operation with `try` to prevent an error from propagating:

```go
>>> try(error("kaboom"), "that failed")
"that failed"
```

## Results

Create a result containing an error message using the `err` built-in function.
This result is treated as any other value in Tamarin:

```go
>>> e := err("io problem")
err("io problem")
>>> type(e)
"result"
>>> e.is_err()
true
>>> e.err_msg()
"io problem"
```

Like with errors, the `try` function works with results:

```go
>>> try(err("io problem"), "fallback value")
"fallback value"
```

Successful results are instead created with the `ok` built-in function:

```go
>>> res := ok("result-value")
ok("result-value")
>>> res.unwrap()
"result-value"
>>> res.is_err()
false
>>> res.is_ok()
true
```

The `try` function returns ok results as-is:

```go
>>> res := ok("result-value")
ok("result-value")
>>> try(res, "fallback value")
"result-value"
```

### Proxying

Results containing an _ok_ value proxy to the wrapped value. This is a convenience
for scripting situations to avoid needing to explicitly unwrap the result. If this
is attempted with an _err_ result, an error is generated that stops execution.

This example shows how the result proxies to the `map.keys` method:

```go
>>> res := ok({foo: "bar"})
ok({"foo": "bar"})
>>> res.keys()
["foo"]
>>> m := res.unwrap()
{"foo": "bar"}
>>> m.keys()
["foo"]
```

## Try

A fallback function can be provided as the second argument to `try`:

```go
>>> e := err("nope")
err("nope")
>>> try(e, func(msg) { print('operation failed: {msg}') })
operation failed: nope
```

## Examples

The JSON module is designed to return result values for operations that may fail:

```go
>>> json.unmarshal("true")
ok(true)
>>> json.unmarshal("invalid-json")
err("value error: json.unmarshal failed with: invalid character 'i' looking for beginning of value")
```
