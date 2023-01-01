# Data Types

Tamarin includes a variety of built-in types. The core types are: int, float,
bool, error, string, list, map, set, result, function, and time. There are also
a handful of iterator types, one for each container type.

Container types may hold a heterogeneous mix of types within. There is not
currently a way to restrict the types a container may hold.

We are reviewing whether optional type hints like in Python or Typescript would
be a useful addition to Tamarin.

## Numeric Types

Int and Float types are the two numeric types in Tamarin. They correspond to
boxed `int64` and `float64` values in Go. Tamarin automatically converts Ints
to Floats in mixed type operations.

The standard set of numeric operators are available when working with these types.

| Operation | Result                | Notes |
| --------- | --------------------- | ----- |
| x + y     | sum of x and y        |       |
| x - y     | difference of x and y |       |
| x \* y    | product of x and y    |       |
| -x        | negation of x         |       |
| x \*\* y  | x to the power of y   |       |
| x += y    | add y to x            |       |
| x -= y    | subtract y from x     |       |

Many math functions are also available in the Tamarin `math` module.

###

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
