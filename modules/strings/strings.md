import { Callout } from 'nextra/components';

# strings

String manipulation functions from the Go standard library.

<Callout type="info" emoji="ℹ️">
    Note that, unlike in Go, many of these functions are also available as
    methods on string objects. See the related documentation [here](/docs/types/string). That said, using these functions **is** handy within pipe expressions.
</Callout>

## Functions

### compare

```go filename="Function signature"
compare(s1, s2 string) int
```

Compares two strings lexicographically. Returns -1 if s1 < s2, 0 if s1 == s2, and 1 if s1 > s2.

```go copy filename="Example"
>>> strings.compare("abc", "abc")
0
>>> strings.compare("abc", "abd")
-1
```

### contains

```go filename="Function signature"
contains(s, substr string) bool
```

Returns true if the string s contains substr.

```go copy filename="Example"
>>> strings.contains("abc", "b")
true
>>> strings.contains("abc", "d")
false
>>> "abc" | strings.contains("b")
true
```

### count

```go filename="Function signature"
count(s, substr string) int
```

Returns the number of non-overlapping instances of substr in s.

```go copy filename="Example"
>>> strings.count("abc", "b")
1
>>> strings.count("ababab", "ab")
3
```

### fields

```go filename="Function signature"
fields(s string) []string
```

Splits the string s around each instance of one or more consecutive white space
characters, returning a slice of substrings or any empty slice if s contains only
white space.

```go copy filename="Example"
>>> strings.fields("a b c")
["a", "b", "c"]
>>> strings.fields("")
[]
```

### has_prefix

```go filename="Function signature"
has_prefix(s, prefix string) bool
```

Returns true if the string s begins with prefix.

```go copy filename="Example"
>>> strings.has_prefix("abc", "a")
true
```

### has_suffix

```go filename="Function signature"
has_suffix(s, suffix string) bool
```

Returns true if the string s ends with suffix.

```go copy filename="Example"
>>> strings.has_suffix("abc", "c")
true
```

### index

```go filename="Function signature"
index(s, substr string) int
```

Returns the index of the first instance of substr in s, or -1 if substr is not
present in s.

```go copy filename="Example"
>>> strings.index("abc", "b")
1
>>> strings.index("abc", "d")
-1
```

### join

```go filename="Function signature"
join(a []string, sep string) string
```

Concatenates the elements of a to create a single string. The separator string
sep is placed between elements in the resulting string.

```go copy filename="Example"
>>> strings.join(["a", "b", "c"], ", ")
"a, b, c"
```

### last_index

```go filename="Function signature"
last_index(s, substr string) int
```

Returns the index of the last instance of substr in s, or -1 if substr is not
present in s.

```go copy filename="Example"
>>> strings.last_index("abc", "b")
1
>>> strings.last_index("abc", "d")
-1
```

### replace_all

```go filename="Function signature"
replace_all(s, old, new string) string
```

Returns a copy of the string s with all non-overlapping instances of old
replaced by new.

```go copy filename="Example"
>>> strings.replace_all("oink oink oink", "oink", "moo")
"moo moo moo"
```

### split

```go filename="Function signature"
split(s, sep string) []string
```

Splits the string s around each instance of sep, returning a slice of
substrings or any empty slice if s does not contain sep.

```go copy filename="Example"
>>> strings.split("a,b,c", ",")
["a", "b", "c"]
```

### to_lower

```go filename="Function signature"
to_lower(s string) string
```

Returns a copy of the string s with all Unicode letters mapped to their lower
case.

```go copy filename="Example"
>>> strings.to_lower("HELLO")
"hello"
```

### to_upper

```go filename="Function signature"
to_upper(s string) string
```

Returns a copy of the string s with all Unicode letters mapped to their upper
case.

```go copy filename="Example"
>>> strings.to_upper("hello")
"HELLO"
```

### trim_prefix

```go filename="Function signature"
trim_prefix(s, prefix string) string
```

Returns s without the provided leading prefix string. If s doesn't start with
prefix, s is returned unchanged.

```go copy filename="Example"
>>> strings.trim_prefix("foo", "f")
"oo"
>>> strings.trim_prefix("foo", "b")
"foo"
```

### trim_space

```go filename="Function signature"
trim_space(s string) string
```

Returns a slice of the string s, with all leading and trailing white space
removed, as defined by Unicode.

```go copy filename="Example"
>>> strings.trim_space("  hello  ")
"hello"
```

### trim_suffix

```go filename="Function signature"
trim_suffix(s, suffix string) string
```

Returns s without the provided trailing suffix string. If s doesn't end with
suffix, s is returned unchanged.

```go copy filename="Example"
>>> strings.trim_suffix("foo", "o")
"f"
```

### trim

```go filename="Function signature"
trim(s, cutset string) string
```

Returns a slice of the string s, with all leading and trailing Unicode code
points contained in cutset removed.

```go copy filename="Example"
>>> strings.trim("¡¡¡Hello, Gophers!!!", "!¡")
"Hello, Gophers"
```
