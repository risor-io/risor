# Iterators

You can enumerate items in any container using an iterator. Iterators are
created using the `iter` built-in function or by using the `range` keyword.

```go
>>> my_set := {1,2,3}
{1, 2, 3}
>>> iter(my_set)
set_iter({1, 2, 3})
```

Unlike in Go, the `range` keyword is available outside of for loop definitions:

```go
>>> range {one: 1, two: 2}
map_iter({"one": 1, "two": 2})
```

Iterators offer a `next` method to retrieve the next entry in the sequence. Each
entry is returned as an `iter_entry` object, which has `key` and `value` attributes.
When the iterator is exhausted, `nil` is returned instead.

For loops recognize when they're working with iterators and automatically assign
each key and value to the loop variables:

```go
>>> s := "abc"
"abc"
>>> for i, c := range s { print("index:", i, "rune:", c) }
index: 0 rune: a
index: 1 rune: b
index: 2 rune: c
```

Iterators can be used directly as well:

```go
>>> s := "abc"
"abc"
>>> iterator := iter(s)
string_iter("abc")
>>> item := iterator.next()
iter_entry(0, "a")
>>> item.key
0
>>> item.value
"a"
```
