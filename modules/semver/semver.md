import { Callout } from 'nextra/components';

# semver

<Callout type="info" emoji="ℹ️">
  This module requires that Risor has been compiled with the `semver` Go build tag.
  When compiling **manually**, [make sure you specify `-tags semver`](https://github.com/risor-io/risor#build-and-install-the-cli-from-source).
</Callout>

## Functions

### compare

```go filename="Function signature"
compare(v1 int, v2 int) int
```

Compares v1 and v2. Returns -1 if v1 is less than v2, 0 if both are equal, 1 if v1 is greater than v2.

```go copy filename="Example"
>>> import semver
>>> semver.compare("1.2.3", "1.2.4")
-1
```

### major

```go filename="Function signature"
major(version string) int
```

Returns the major version of the given version string.

```go copy filename="Example"
>>> import semver
>>> semver.major("1.2.3")
1
```

### minor

```go filename="Function signature"
minor(version string) int
```

Returns the minor version of the given version string.

```go copy filename="Example"
>>> import semver
>>> semver.minor("1.2.3")
2
```

### patch

```go filename="Function signature"
patch(version string) int
```

Returns the patch version of the given version string.

```go copy filename="Example"
>>> import semver
>>> semver.patch("1.2.3")
3
```

### build

```go filename="Function signature"
build(version string) string
```

Returns the build version of the given version string.

```go copy filename="Example"
>>> import semver
>>> semver.build("1.2.3+build")
"build"
```

### pre

```go filename="Function signature"
pre(version string) string
```

Pre returns the pre-release version of the given version string.

```go copy filename="Example"
>>> import semver
>>> semver.pre("1.2.3-pre")
"pre"
```

### validate

```go filename="Function signature"
validate(version string) bool
```

Returns an error if the version isn't valid.

```go copy filename="Example"
>>> import semver
>>> semver.validate("1.2.3invalid")
Invalid character(s) found in patch number "3invalid"
```

### parse

```go filename="Function signature"
parse(version string) map
```

Parses the given version string and returns a map with the major, minor, patch, pre-release, and build versions.

```go copy filename="Example"
>>> import semver
>>> semver.parse("1.2.3-pre+build")
{
  "major": 1,
  "minor": 2,
  "patch": 3,
  "pre": "pre",
  "build": "build"
}
```

### equals

```go filename="Function signature"
equals(v1 string, v2 string) bool
```

Returns whether v1 and v2 are equal.

```go copy filename="Example"
>>> import semver
>>> semver.equals("1.2.3", "1.2.3")
true
```
