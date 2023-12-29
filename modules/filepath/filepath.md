# filepath

The `filepath` module contains utilities thare are used to manipulate file paths
with operating system-specific separators.

## Functions

### abs

```go filename="Function signature"
abs(path string) string
```

Returns the absolute representation of path. If the path is not absolute it will
be joined with the current working directory to create the corresponding absolute path.

```go copy filename="Example"
>>> filepath.abs("foo/bar")
"/home/user/foo/bar"
>>> filepath.abs("../foo/bar")
"/home/foo/bar"
```

Learn more: [filepath.Abs](https://pkg.go.dev/path/filepath#Abs).

### base

```go filename="Function signature"
base(path string) string
```

Returns the last element of path with any trailing slashes removed.
If path is empty, "." is returned.

```go copy filename="Example"
>>> filepath.base("foo/bar")
"bar"
>>> filepath.base("foo/bar/")
"bar"
>>> filepath.base("test.txt")
"test.txt"
>>> filepath.base("")
"."
```

Learn more: [filepath.Base](https://pkg.go.dev/path/filepath#Base).

### clean

```go filename="Function signature"
clean(path string) string
```

Returns the shortest path name equivalent to path by purely lexical processing.

```go copy filename="Example"
>>> filepath.clean("foo/bar/../baz")
"foo/baz"
>>> filepath.clean("foo/bar/./baz")
"foo/bar/baz"
>>> filepath.clean("foo/bar/../../baz")
"baz"
```

Learn more: [filepath.Clean](https://pkg.go.dev/path/filepath#Clean).

### dir

```go filename="Function signature"
dir(path string) string
```

Returns all but the last element of path, typically the path's directory.

```go copy filename="Example"
>>> filepath.dir("foo/bar")
"foo"
>>> filepath.dir("foo/bar/")
"foo"
>>> filepath.dir("test.txt")
"."
```

Learn more: [filepath.Dir](https://pkg.go.dev/path/filepath#Dir).

### ext

```go filename="Function signature"
ext(path string) string
```

Returns the file name extension used by path. The extension is the suffix
beginning at the final dot in the final element of path. The result it is empty
if there is no dot.

```go copy filename="Example"
>>> filepath.ext("foo/bar.txt")
".txt"
>>> filepath.ext("foo/bar")
""
>>> filepath.ext("foo/bar.tar.gz")
".gz"
```

Learn more: [filepath.Ext](https://pkg.go.dev/path/filepath#Ext).

### is_abs

```go filename="Function signature"
is_abs(path string) bool
```

Returns true if the path is absolute.

```go copy filename="Example"
>>> filepath.is_abs("/foo/bar")
true
>>> filepath.is_abs("foo/bar")
false
```

Learn more: [filepath.IsAbs](https://pkg.go.dev/path/filepath#IsAbs).

### join

```go filename="Function signature"
join(paths ...string) string
```

Returns the result of joining the given path elements with the operating
system-specific path separator.

```go copy filename="Example"
>>> filepath.join("foo", "bar")
"foo/bar"
>>> filepath.join("foo", "bar", "baz")
"foo/bar/baz"
```

Learn more: [filepath.Join](https://pkg.go.dev/path/filepath#Join).

### match

```go filename="Function signature"
match(pattern, name string) bool
```

Returns true if the file name matches the shell pattern.

```go copy filename="Example"
>>> filepath.match("*.txt", "foo.txt")
true
>>> filepath.match("*.txt", "foo.tar.gz")
false
```

Learn more: [filepath.Match](https://pkg.go.dev/path/filepath#Match).

### rel

```go filename="Function signature"
rel(basepath, targpath string) string
```

Returns a relative path that is lexically equivalent to targpath when joined
to basepath with an intervening separator.

```go copy filename="Example"
>>> filepath.rel("/home/user", "/home/user/foo/bar")
"foo/bar"
>>> filepath.rel("/home/user", "/home/user/foo/../bar")
"bar"
>>> filepath.rel("/home/user", "/home/user/foo/../../bar")
"../bar"
```

Learn more: [filepath.Rel](https://pkg.go.dev/path/filepath#Rel).

### split_list

```go filename="Function signature"
split_list(path string) []string
```

Splits path immediately following the final separator, separating it into a
directory and file name component. If there is no separator in path, split_list
returns an empty dir and file set to path.

```go copy filename="Example"
>>> filepath.split_list("/home/user/foo/bar")
["/home/user/foo", "bar"]
>>> filepath.split_list("/home/user/foo")
["/home/user", "foo"]
>>> filepath.split_list("foo")
["", "foo"]
```

Learn more: [filepath.Split](https://pkg.go.dev/path/filepath#Split).

### split

```go filename="Function signature"
split(path string) []string
```

Splits the path immediately following the final separator, returning a list of two items: the directory and the file name. If there is no separator in the path, an empty directory and the file name are returned.

```go copy filename="Example"
>>> filepath.split("/home/user/foo/bar")
["/home/user/foo", "bar"]
>>> filepath.split("/home/user/foo")
["/home/user", "foo"]
>>> filepath.split("test.txt")
["", "test.txt"]
```

Learn more: [filepath.Split](https://pkg.go.dev/path/filepath#Split).

### walk_dir

```go filename="Function signature"
walk_dir(root string, fn func(path string))
```

Walks the file tree at root, calling fn for each file or directory in the tree,
including root. Files are walked in lexical order. Symbolic links are not
followed.

```go copy filename="Example"
>>> filepath.walk_dir("/home/user/foo", func(path, dir_entry, err) { print(path) })
"/home/user/foo"
"/home/user/foo/bar"
"/home/user/foo/test.txt"
```

Learn more: [filepath.WalkDir](https://pkg.go.dev/path/filepath#WalkDir).
