# bytes

## Functions

### clone

```go filename="Function signature"
clone(b byte_slice) byte_slice
```

Clone returns a new byte slice containing the same bytes as the given byte slice.

```go copy filename="Example"
>>> bytes.clone(byte_slice([1, 2, 3, 4]))
byte_slice("\x01\x02\x03\x04")
```

### contains_any

```go filename="Function signature"
contains_any(b byte_slice, chars string) bool
```

Reports whether any of the UTF-8-encoded code points in the string
are present in the byte_slice.

```go copy filename="Example"
>>> bytes.contains_any(byte_slice("Hello"), "abco")
true
>>> bytes.contains_any(byte_slice("Hello"), "abc")
false
```

### contains_rune

```go filename="Function signature"
contains_rune(b byte_slice, r rune) bool
```

Reports whether the rune is contained in the UTF-8-encoded byte slice.

```go copy filename="Example"
>>> bytes.contains_rune(byte_slice("Hello"), "H")
true
>>> bytes.contains_rune(byte_slice("Hello"), "h")
false
```

### contains

```go filename="Function signature"
contains(b, subslice byte_slice) bool
```

Contains reports whether subslice is within b.

```go copy filename="Example"
>>> bytes.contains(byte_slice("seafood"), byte_slice("foo"))
true
>>> bytes.contains(byte_slice("seafood"), byte_slice("bar"))
false
```

### count

```go filename="Function signature"
count(s, sep byte_slice) int
```

Counts the number of non-overlapping instances of sep in s. If sep is an empty
slice, Count returns 1 + the number of UTF-8-encoded code points in s.

```go copy filename="Example"
>>> bytes.count(byte_slice("cheese"), byte_slice("e"))
3
```

### equals

```go filename="Function signature"
equals(a, b byte_slice) bool
```

Reports whether a and b are the same length and contain the same bytes.

```go copy filename="Example"
>>> bytes.equals(byte_slice("Hello"), byte_slice("Hello"))
true
>>> bytes.equals(byte_slice("Hello"), byte_slice("hello"))
false
```

### has_prefix

```go filename="Function signature"
has_prefix(s byte_slice, prefix byte_slice) bool
```

Tests whether the byte slice s begins with prefix.

```go copy filename="Example"
>>> bytes.has_prefix(byte_slice("Gopher"), byte_slice("Go"))
true
>>> bytes.has_prefix(byte_slice("Gopher"), byte_slice("C"))
false
```

### has_suffix

```go filename="Function signature"
has_suffix(s byte_slice, suffix byte_slice) bool
```

Tests whether the byte slice s ends with suffix.

```go copy filename="Example"
>>> bytes.has_suffix(byte_slice("Amigo"), byte_slice("go"))
true
>>> bytes.has_suffix(byte_slice("Amigo"), byte_slice("O"))
false
```

### index_any

```go filename="Function signature"
index_any(s byte_slice, chars string) int
```

Interprets s as a sequence of UTF-8-encoded code points. Returns the byte
index of the first occurrence in s of any of the code points in chars. Returns
-1 if chars is empty or if there are no code point in common.

```go copy filename="Example"
>>> bytes.index_any(byte_slice("chicken"), "aeiou")
2
>>> bytes.index_any(byte_slice("bcd"), "aeiou")
-1
```

### index_byte

```go filename="Function signature"
index_byte(b byte_slice, c byte) int
```

Returns the index of the first occurrence of c in b, or -1 if c is not present.

```go copy filename="Example"
>>> bytes.index_byte(byte_slice("golang"), "g")
0
>>> bytes.index_byte(byte_slice("golang"), "x")
-1
```

### index_rune

```go filename="Function signature"
index_rune(s byte_slice, r rune) int
```

Interprets s as a sequence of UTF-8-encoded code points. Returns the byte index
of the first occurrence in s of the given rune. Returns -1 if rune is not present.

```go copy filename="Example"
>>> bytes.index_rune(byte_slice("chicken"), "k")
4
>>> bytes.index_rune(byte_slice("chicken"), "d")
-1
```

### index

```go filename="Function signature"
index(s, sep byte_slice) int
```

Returns the index of the first occurrence of sep in s, or -1 if sep is not present.

```go copy filename="Example"
>>> bytes.index(byte_slice("chicken"), byte_slice("ken"))
4
>>> bytes.index(byte_slice("chicken"), byte_slice("kex"))
-1
```

### repeat

```go filename="Function signature"
repeat(b byte_slice, count int) byte_slice
```

Returns a new byte slice consisting of count copies of b.

```go copy filename="Example"
>>> bytes.repeat(byte_slice("a"), 3)
byte_slice("aaa")
```

### replace_all

```go filename="Function signature"
replace_all(s, old, new byte_slice) byte_slice
```

Returns a copy of the slice s with all non-overlapping instances of old
replaced by new.

```go copy filename="Example"
>>> bytes.replace_all(byte_slice("aaa"), byte_slice("a"), byte_slice("b"))
byte_slice("bbb")
```

### replace

```go filename="Function signature"
replace(s, old, new byte_slice, n int) byte_slice
```

Returns a copy of the slice s with the first n non-overlapping instances of old
replaced by new.

```go copy filename="Example"
>>> bytes.replace(byte_slice("aaa"), byte_slice("a"), byte_slice("b"), 2)
byte_slice("bba")
```
