# Data Types

Risor includes a variety of built-in types. The core types are: int, float,
bool, error, string, list, map, set, result, function, and time. There are also
a handful of iterator types, one for each container type.

Container types may hold a heterogeneous mix of types within. There is not
currently a way to restrict the types a container may hold.

Optional type hints like found in Python or Typescript may be a future addition
to Risor.

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
```

## Numerics

Int and Float types are the two numeric types in Risor. They correspond to
boxed `int64` and `float64` values in Go. Risor automatically converts Ints
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

```go
>>> x := 2
2
>>> y := 3
3
>>> x * y
6
>>> x + y
5
>>> type(x + y)
"int"
>>> type(x + float(y))
"float"
```

Many math functions are also available in the Risor `math` module.

### Related Built-ins

#### float(x)

Converts a String or Int object to a Float. An error is generated if the
operation fails.

```go
>>> float("4.4")
4.4
```

#### int(x)

Converts a String or Float to an Int. An error is generated if the operation
fails.

```go
>>> int(4.4)
4
>>> int("123")
123
```

## Bool

The `bool` type in Risor is a simple wrapper of the Go `bool` type.

All underlying object types in Risor implement the `object.Object` interface,
which includes `IsTruthy()` and `Equals(other)` methods. It's good to keep this
in mind since object "truthiness" can be leveraged in conditional statements.

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
>>> [1,2] != [1,2]
false
>>> false == false
true
```

### Related Built-ins

#### bool(x)

Returns `true` or `false` according to the given objects "truthiness". Container
types including lists, maps, sets, and strings evaluate to `false` when empty
and `true` otherwise. Iterators are truthy when there are more items remaining
to iterate over. Objects of other types are generally always considered to
be truthy.

## String

Strings in Risor are based on the underlying `string` type in Go. As such,
they support unicode and various operations like indexing operate on the
underlying runes within the string.

### Quote Types

There are three ways to quote strings in Risor source code:

```
'single quotes: supports interpolated {vars}'
"double quotes: equivalent to Go strings"
`backticks: raw strings that may span multiple lines`
```

The single quoted string approach offers string formatting via interpolation,
much like [f-strings](https://peps.python.org/pep-0498/) in Python. Arbitrary
Risor expressions can be embedded within parentheses and resolved during
evaluation. In Risor, the restriction on these expressions is that they
cannot contain curly braces, since those are used to mark the beginning and
end of the template expression.

An example with simple expressions:

```go
>>> name := "jean"
"jean"
>>> age := 30
30
>>> '{name} is {age} years old'
"jean is 30 years old"
```

Another example:

```go
>>> nums := [0, 1, 2, 3]
[0, 1, 2, 3]
>>> 'the max is {math.max(nums)} and the length is {len(nums)}'
"the max is 3 and the length is 4"
```

### Container Operations

Strings in Risor implement the `object.Container` interface, which means they
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

### Related Built-ins

#### chr(int)

Converts an Int to the corresponding unicode rune, which is returned as a String.
The `ord` built-in performs the inverse transformation.

#### ord(string)

Converts a unicode character to the corresponding Int value. The `chr` built-in
performs the inverse transformation. An error is generated if a multi-rune string is
provided.

#### sprintf(string, ...any)

Wraps `fmt.Sprintf` to format the string with the provided arguments. Risor
objects are converted to their corresponding Go types before being passed to
`fmt.Sprintf`.

#### string(x)

Returns the string representation of any Risor object.

### Methods

#### string.contains(s)

Returns a bool that indicates if `s` is a substring of this string.

#### string.has_prefix(s)

Checks whether the string begins with the prefix `s`.

#### string.has_suffix(s)

Checks whether the string ends with the suffix `s`.

#### string.count(s)

Returns the number of occurrences of `s` in this string.

#### string.join(list)

Return the joined result of the given list of strings, using this string
as the separator.

#### string.split(separator)

Splits this string on all occurrences of the given separator, returning
the resulting list of strings.

#### string.fields()

Splits this string on whitespace, returning the list of non-whitespace substrings. If this string is only whitespace, an empty list is returned.

#### string.index(s)

Returns the index of the first occurence of `s` in this string, or `-1`
if `s` is not present.

#### string.last_index(s)

Returns the index of the last occurence of `s` in this string, or `-1`
if `s` is not present.

#### string.replace_all(old, new)

Returns a copy of this string with all occurrences of `old` replaced by `new`.

#### string.to_lower()

Returns a copy of this string that is transformed to all lowercase.

#### string.to_upper()

Returns a copy of this string that is transformed to all uppercase.

#### string.trim(cutset)

Returns a copy of this string with all leading and trailing characters contained
in `cutset` removed.

#### string.trim_prefix(prefix)

Returns a copy of this string without the given prefix. This is a no-op if this
string doesn't start with `prefix`.

#### string.trim_space()

Returns a copy of this string without the leading and trailing whitespace.

#### string.trim_suffix(suffix)

Returns a copy of this string without the given suffix. This is a no-op if this
string doesn't end with `suffix`.

## List

Lists in Risor behave very similarly to lists in Python. Methods and indexing
on lists are the primary way to interact with these objects. A list can store
objects of any types, including a mix of types in one list.

```go
>>> l := ["a", 1, 2]
["a", 1, 2]
>>> l.append("tail")
["a", 1, 2, "tail"]
```

### Container Operations

```go
>>> l := ["a", "b", "c"]
["a", "b", "c"]
>>> len(l)
3
>>> "c" in l
true
>>> "d" in l
false
>>> l[2]
"c"
>>> l[2] = "d"
>>> l
["a", "b", "d"]
>>> l[1:]
["b", "d"]
>>> l[:1]
["a"]
```

### Related Built-ins

#### list(container)

Returns a new list populated with items from the given container. If a list is
provided, a shallow copy of the list is returned. It is also commonly used to
convert a set to a list.

```go
>>> s := {"a", "b", "c"}
{"a", "b", "c"}
>>> list(s)
["a", "b", "c"]
```

### Methods

#### list.append(x)

Adds x to the end of the list.

#### list.clear()

Empties all items from the list.

#### list.copy()

Returns a shallow copy of the list.

#### list.count(x)

Returns a count of how many times x is found in the list.

#### list.extend(x)

Adds all items contained in x to the end of the list.

#### list.index(x)

Returns the first index of x in the list, or -1 if not found.

#### list.insert(index, x)

Inserts x into the list at the specified index.

#### list.pop(index)

Removes the item at the given index from the list.

#### list.remove(x)

Removes the first occurence of x in the list.

#### list.reverse()

Reverses the list in place.

#### list.sort()

Sorts the list in place.

#### list.map(func)

Returns a transformed list, in which the given function is applied to each list item.

#### list.filter(func)

Returns a transformed list, in which the given function returns true for items that should be added to the output list.

#### list.each(func)

Calls the supplied function once with each item in the list.

## Map

Maps associate keys with values and provide fast lookups by key. Risor
maps use underlying Go maps of type `map[string]interface{}`. This means
Risor maps always operate with string keys, which provides equivalence with JSON.

```go
>>> m := {one: 1, two: 2}
{"one": 1, "two": 2}
>>> m["three"] = 3
>>> m
{"one": 1, "three": 3, "two": 2}
```

### Container Operations

```go
>>> m := {"name": "sean", "age": 27}
{"age": 27, "name": "sean"}
>>> len(m)
2
>>> "age" in m
true
>>> m["age"]
27
>>> m["age"] = 28
>>> m
{"age": 28, "name": "sean"}
>>> m.keys()
["age", "name"]
```

### Related Built-ins

#### delete(map, key)

Deletes the item with the specified key from the map. This operation has no
effect if the key is not present in the map.

#### map(container)

Returns a new map with the contents of the given container. Generally, containers
are transformed into the map by creating an iterator for the given container and
the key and value for each iterator entry are added to the map. As a special
case, if the container is a list then it is expected to be a nested list of
key-value pairs, e.g. `[["key1", "val1"]]`. Any non-string keys that are
encountered are automatically converted to their string representation.

```go
>>> map({"a", "b", "c"})
{"a": true, "b": true, "c": true}
>>> map("abc")
{"0": "a", "1": "b", "2": "c"}
>>> map([["name", "joe"], ["age", 18]])
{"age": 18, "name": "joe"}
```

### Methods

#### map.clear()

Removes all items from the map.

#### map.copy()

Returns a shallow copy of the map, containing the same keys and values.

#### map.get(key, default=nil)

Returns the value associated with the given key, if it exists in the map.
If the key is not in the map, the given default value is returned.

#### map.items()

Returns a list of [key, value] pairs containing the items from the map.

#### map.keys()

Returns a sorted list of keys contained in the map.

#### map.pop(key, default=nil)

Returns the value associated with the given key and then removes it from the
map. If the key is not in the map, the given default value is returned instead.

#### map.setdefault(key, default)

Sets the key to the given default value if the key is not already in the map.
If the key already is in the map, do nothing. Returns the value associated with
the key after the set action.

#### map.update(other)

Updates this map with the key-value pairs contained in the provided map,
overwriting any items with matching keys already in this map.

#### map.values()

Returns a sorted list of values contained in the map.

## Set

Sets represent an unordered collection of unique objects. Only hashable objects
can be added to sets, which includes bool, int, float, nil, and string. It is
not possible to add a list or map to a set, since they are not hashable.

### Container Operations

```go
>>> s := {1, 2, 3}
{1, 2, 3}
>>> 3 in s
true
>>> 4 in s
false
>>> s[3]
true
>>> s[4]
false
>>> delete(s, 3)
>>> s
{1, 2}
```

### Related Built-ins

#### delete(set, key)

Deletes the given key from the set. This is a no-op if the key is not present in the set.

#### set(container)

Returns a new set that is populated with the contents of the given container.
The loading behavior for each provided type is as follows:

- Given a list, the values within are added to the set.
- Given a string, the characters within are added to the set.
- Given a map, the keys within are added to the set.

```go
>>> set([1, 2, 3, 3])
{1, 2, 3}
>>> set("abc")
{"a", "b", "c"}
>>> set({one: 1, two: 2})
{"one", "two"}
```

### Methods

#### set.add(x)

Add `x` to the set. If `x` is not hashable, an error is generated.

#### set.clear()

Empties all items from the set.

#### set.remove(x)

Remove `x` from the set. This is a no-op if `x` is not in the set.

#### set.union(other)

Returns a new set which contains the union of all items from this set
and the other set.

#### set.intersection(other)

Returns a new set containing items that are present in both this set and the other set.
