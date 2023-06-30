# Built-in Functions

Risor includes this set of default built-in functions. The set of available
built-ins is easily customizable, depending on the goals for your project.

### all(container)

Returns `true` if all entries in the given container are "truthy".

```go
>>> all([true, 1, "ok"])
true
>>> all([true, 0, "ok"])
false
```

### any(container)

Returns `true` if any of the entries in the given container are "truthy".

```go
>>> any([false, 0, "ok"])
true
>>> any([false, 0, ""])
false
```

### assert(x, message)

Generates an error if `x` is "falsy". If a message is provided, it is used as
the assertion error message.

```go
>>> assert(1 == 1, "check failed")
>>> assert(1 == 2, "check failed")
check failed
```

### bool(object)

Returns `true` or `false` depending on whether the object is considered "truthy".
Container types including lists, maps, sets, and strings evaluate to `false` when
empty and `true` otherwise.

```go
>>> bool(1)
true
>>> bool(0)
false
>>> bool([1])
true
>>> bool([])
false
```

### call(function, ...any)

Calls the function with given arguments. This is primarily useful in pipe
expressions when a function is being passed through the pipe as a variable.

```go
>>> func inc(x) { x + 1 }
>>> call(inc, 99)
100
>>> inc | call(41)
42
```

### chr(int)

Converts an Int to the corresponding unicode rune, which is returned as a String.
The `ord` built-in performs the inverse transformation.

```go
>>> chr(8364)
"€"
>>> chr(97)
"a"
```

### delete(map, key)

Deletes the item with the specified key from the map. This operation has no
effect if the key is not present in the map.

```go
>>> m := {one: 1, two: 2}
{"one": 1, "two": 2}
>>> delete(m, "one")
{"two": 2}
>>> delete(m, "foo")
{"two": 2}
```

### err(message)

Returns a new Result object containing the given error message. Results may
contain an `ok` value or an `err` value, and this built-in is used to construct
the latter. This is similar in some ways to [Rust's result type](https://doc.rust-lang.org/std/result/).

```go
>>> err("failed operation")
err("failed operation")
```

### error(message)

Generates an Error containing the given message. Errors in Risor stop evaluation
when they are generated, unless a `try` call is used to stop error propagation.

```go
>>> error("kaboom")
kaboom
```

### float(object)

Converts a String or Int object to a Float. An error is generated if the
operation fails.

```go
>>> float("4.4")
4.4
```

### getattr(object, name, default)

Returns the named attribute from the object, or the default value if the
attribute does not exist. The returned attribute is always a Risor object,
which may be a function. This is similar to
[getattr](https://docs.python.org/3/library/functions.html#getattr) in Python.

```go
>>> l := [1,2,3]
[1, 2, 3]
>>> append := getattr(l, "append")
builtin(list.append)
>>> append(4)
[1, 2, 3, 4]
>>> getattr(l, "unknown", "that doesn't exist")
"that doesn't exist"
```

### int(object)

Converts a String or Float to an Int. An error is generated if the operation
fails.

```go
>>> int(4.4)
4
>>> int("123")
123
```

### iter(container)

Returns an iterator for the given container object. This can be used to iterate
through items in a for loop or interacted with more directly. The returned
iterator has a `next()` method that will return `iter_entry` objects.

```go
>>> s := {"a", "b", "c"}
{"a", "b", "c"}
>>> iterator := iter(s)
set_iter({"a", "b", "c"})
>>> iterator.next().key
"a"
```

### keys(container)

Returns a list of all keys for items in the given map or list container.

```go
>>> m := {one: 1, two: 2}
{"one": 1, "two": 2}
>>> keys(m)
["one", "two"]
```

### len(container)

Returns the size of the string, list, map, or set.

```go
>>> len("ab")        // string length
2
>>> len([1,2,3])     // list length
3
>>> len({foo:"bar"}) // map length
1
>>> len({1,2,3,4})   // set length
4
```

### list(container)

Returns a new list populated with items from the given container. If a list is
provided, a shallow copy of the list is returned. It is also commonly used to
convert a set to a list.

```go
>>> s := {"a", "b", "c"}
{"a", "b", "c"}
>>> list(s)
["a", "b", "c"]
```

### ok(object)

Returns a new Result containing the given ok value. Results may contain an `ok`
value or an `err` value, and this built-in is used to construct the former. This
is similar in some ways to [Rust's result type](https://doc.rust-lang.org/std/result/).

```go
>>> result := ok("that worked")
ok("that worked")
>>> result.is_ok()
true
>>> result.unwrap()
"that worked"
```

### ord(string)

Converts a unicode character to the corresponding Int value. The `chr` built-in
performs the inverse transformation. An error is generated if a multi-rune string is
provided.

```go
>>> ord("€")
8364
>>> ord("a")
97
>>> chr(ord("€"))
"€"
```

### print(...any)

Prints the provided objects to stdout after converting them to their String
representations. Spaces are inserted between each object and a trailing newline
is output. This is a wrapper around `fmt.Println` in Go.

```go
>>> print(42, "is the answer")
42 is the answer
```

### printf(string, ...any)

Printf wraps `fmt.Printf` in order to print the formatted string and arguments
to stdout. In the Risor REPL you will currently not see the `printf` output
unless the string ends in a newline character.

```go
>>> printf("name: %s age: %d\n", "joe", 32)
name: joe age: 32
```

### reversed(list)

Returns a new list which is a reversed copy of the provided list.

```go
>>> l := ["a", "b", "c"]
["a", "b", "c"]
>>> reversed(l)
["c", "b", "a"]
>>> l
["a", "b", "c"]
```

### set(container)

Returns a new set containing the items from the given container object.

```go
>>> set("aabbcc")
{"a", "b", "c"}
>>> set([4,4,5])
{4, 5}
>>> set({one:1, two:2})
{"one", "two"}
```

### sorted(container)

Returns a sorted list of items from the given container object.

```go
>>> sorted("cba")
["a", "b", "c"]
>>> sorted([10, 3, -5])
[-5, 3, 10]
```

### sprintf(string, ...any)

Wraps `fmt.Sprintf` to format the string with the provided arguments. Risor
objects are converted to their corresponding Go types before being passed to
`fmt.Sprintf`.

```go
>>> sprintf("name: %s age: %d", "fred", 18)
"name: fred age: 18"
>>> sprintf("%v", [1, "a", 3.3])
"[1 a 3.3]"
```

### string(object)

Returns a string representation of the given Risor object.

```go
>>> string({one:1, two:2})
"{\"one\": 1, \"two\": 2}"
>>> string(4.4)
"4.4"
>>> string([1,2,3])
"[1, 2, 3]"
```

### try(expression, fallback)

Evaluates the expression and return the fallback value if an error occurs. This
works equally well with Error objects and Results containing an error value.

If a function is provided as the fallback value, that function is called when
an error occurs with the error message as the first argument. The value returned
from the function is used as the return value from the `try` call.

```go
>>> try("ok", "fallback")
"ok"
>>> try(error("boom"), "that failed")
"that failed"
>>> try(error("boom"), error("transformed err"))
transformed err
>>> r := err("err result")
err("err result")
>>> try(r, "the result contained an error")
"the result contained an error"
>>> try(r, func(msg) { 'failure: {msg}' })
"failure: err result"
```

### type(object)

Returns the type name of the given object as a String.

```go
>>> type(1)
"int"
>>> type(2.2)
"float"
>>> type("hi")
"string"
>>> type([])
"list"
>>> type({})
"map"
>>> type({1,2,3})
"set"
>>> type(ok("success"))
"result"
>>> type(err("failed"))
"result"
```

### unwrap_or(result, fallback)

Returns the wrapped "ok" value from the given Result object or the fallback
value if the Result contains an Error.

```go
>>> unwrap_or(ok("success"), "fallback")
"success"
>>> unwrap_or(err("boom"), "fallback")
"fallback"
```

### unwrap(result)

Returns the wrapped "ok" value from the given Result object. If the Result
contains an Error, an Error is generated instead.

```go
>>> unwrap(ok("success"))
"success"
>>> unwrap(err("boom"))
result error: unwrap() called on an error: error("boom")
```
