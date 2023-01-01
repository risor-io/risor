# Data Types

Tamarin includes a variety of built-in types. The core types are: int, float,
bool, error, string, list, map, set, result, function, and time. There are also
a handful of iterator types, one for each container type.

Container types may hold a heterogeneous mix of types within. There is not
currently a way to restrict the types a container may hold.

We are reviewing whether optional type hints like in Python or Typescript would
be a useful addition to Tamarin.

## Numerics

Int and Float types are the two numeric types in Tamarin. They correspond to
boxed `int64` and `float64` values in Go. Tamarin automatically converts Ints
to Floats in mixed type operations.

The standard set of numeric operators are available when working with these types.

| Operation | Result                |
| --------- | --------------------- |
| x + y     | sum of x and y        |
| x - y     | difference of x and y |
| x \* y    | product of x and y    |
| -x        | negation of x         |
| x \*\* y  | x to the power of y   |
| x += y    | add y to x            |
| x -= y    | subtract y from x     |
| x \*= y   | multiply x by y       |
| x /= y    | divide x by y         |

Many math functions are also available in the Tamarin `math` module.

## Bool

The `bool` type in Tamarin is a simple wrapper of the Go `bool` type. Tamarin
requires all `object.Object` types to implement `IsTruthy() bool` and
`Equals(other)`, which are two common situations dealing with booleans.

```go
>>> bool(0)
false
>>> bool(5)
true
>>> bool([])
false
>>> bool([1,2,3])
true
>>> if 5 { print("5 is truthy") }
5 is truthy
>>> 5 == 5.0
true
>>> [1,2] == [1,2]
true
```

## String

Strings in Tamarin are based on the underlying `string` type in Go. As such,
they support unicode and various operations like indexing operate on the
underlying runes within the string.

Much of the `strings` Go module is exposed in a Tamarin module of the same name.

There are three ways to quote strings in Tamarin source code:

```
'single quotes: supports template {vars}'
"double quotes: equivalent to Go strings"
`backticks: raw strings that may span multiple lines`
```

Strings in Tamarin implement the `object.Container` interface, which means they
support typical container-style operations:

```go
>>> s := "hello"
"hello"
>>> s[0]
"h"
>>> len(s)
5
>>> s[1:3]
"el"
>>> s[1:]
"ello"
>>> s[:1]
"h"
>>> iter(s)
string_iter("hello")
>>> iter(s).next()
iter_entry(0, "h")
```

## List

Lists in Tamarin behave very similarly to lists in Python. Methods and indexing
on lists are the primary way of mutating and working with these objects.

```go
>>> l := ["a", 1, {2,3,4}]
["a", 1, {2, 3, 4}]
>>> l[0]
"a"
>>> l[0] = "b"
>>> l
["b", 1, {2, 3, 4}]
>>> len(l)
3
>>> l.each(func(item) { print(item) })
b
1
{2, 3, 4}
>>> l.append("tail")
["b", 1, {2, 3, 4}, "tail"]
>>> iter(l)
list_iter(["b", 1, {2, 3, 4}, "tail"])
>>> iter(l).next()
iter_entry(0, "b")
```

## Quick Reference

```go
101         // int
1.1         // float
"1"         // string
[1,2,3]     // list
{"key":2}   // map
{1,2}       // set
false       // bool
nil         // nil
func() {}   // function
time.now()  // time
iter([1])   // list_iter
iter(set()) // set_iter
iter({})    // map_iter
iter("1")   // string_iter
```

There are also `HttpResponse` and `DatabaseConnection` types in progress.
