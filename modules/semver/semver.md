# semver

The `semver` module provides functions to work with semantic version strings.

Semantic version strings are in the form:

```
vMAJOR[.MINOR[.PATCH[-PRERELEASE][+BUILD]]]
```

The square brackets indicate optional fields with the string.

The MAJOR, MINOR, and PATCH strings must be decimal integers without leading zeros.
The PRERELEASE and BUILD strings may contain alphanumeric identifiers and hyphens.
If the PRERELEASE string is numeric, it must not have leading zeros.

This implementation matches [Semver version 2](https://semver.org/) with two
exceptions:

1. The version string must begin with a leading "v".
2. The MINOR and PATCH fields are optional and default to 0 if not specified.

The core functionality is provided by
[golang.org/x/mod/semver](https://pkg.go.dev/golang.org/x/mod/semver).

## Functions

### build

```go filename="Function signature"
build(v string) string
```

Returns the build metadata of the specified version string.

```go filename="Example"
>>> semver.build("v1.2.3+build234")
"+build234"
```

### canonical

```go filename="Function signature"
canonical(v string) string
```

Returns the canonical formatting of the semantic version string. The function
fills in any missing minor or patch fields with zeros and removes any build
metadata.

```go filename="Example"
>>> semver.canonical("v1")
"v1.0.0"
```

### compare

```go filename="Function signature"
compare(v1, v2 string) int
```

Compares two version strings and returns -1, 0, or 1 if `v1` is less than, equal
to, or greater than `v2`, respectively.

```go filename="Example"
>>> semver.compare("v1.2.3", "v1.2.4")
-1
```

### is_valid

```go filename="Function signature"
is_valid(v string) bool
```

Returns true if the specified version string is a valid semantic version.

```go filename="Example"
>>> semver.is_valid("v1.2.3")
true
>>> semver.is_valid("3")
false
```

### major

```go filename="Function signature"
major(v string) string
```

Returns the major version of the specified version string.

```go filename="Example"
>>> semver.major("v1.2.3")
"v1"
```

### major_minor

```go filename="Function signature"
major_minor(v string) string
```

Returns the major and minor version of the specified version string.

```go filename="Example"
>>> semver.major_minor("v1.2.3")
"v1.2"
```

### prerelease

```go filename="Function signature"
prerelease(v string) string
```

Returns the prerelease information of the specified version string.

```go filename="Example"
>>> semver.prerelease("v1.2.3-alpha")
"alpha"
```
