# Iterators

You can step through items in any container using an iterator. You can create
an iterator using the `iter` builtin function or by using the `range` keyword.

```go
>>> iter({1,2,3})
set_iter({1, 2, 3})
```

```go
>>> range {one: 1, two: 2}
map_iter({"one": 1, "two": 2})
```

Iterators offer a `next` method to retrieve the next entry in the sequence. Each
entry is returned as an `iter_entry` object, which has `key` and `value` attributes.
When the iterator is exhausted, `nil` is returned instead.

For loops work with these iterators natively, and automatically assign the
key and value to the loop variables. But you can use iterators directly as well.

```go
>>> entry := range {foo: "bar"}.next()
iter_entry("foo", "bar")
>>> entry.key
"foo"
>>> entry.value
"bar"
```
