# isatty

The `isatty` module provides functions to check if the current process is
connected to a terminal.

The core functionality is provided by
[github.com/mattn/go-isatty](https://github.com/mattn/go-isatty).

## Module

```go copy filename="Function signature"
isatty() bool
```

The `isatty` module object itself is callable, and returns a boolean indicating
whether the process is connected to a terminal. This module-level function does
not differentiate between cygwin and non-cygwin terminals, returning true in
both cases.

```go copy filename="Example"
>>> isatty()
true
```

## Functions

### is_terminal

```go filename="Function signature"
is_terminal() bool
```

Returns true if the current process is connected to a terminal.

```go copy filename="Example"
>>> isatty.is_terminal()
true
```

### is_cygwin_terminal

```go filename="Function signature"
is_cygwin_terminal() bool
```

Returns true if the current process is connected to a Cygwin terminal.

```go copy filename="Example"
>>> isatty.is_cygwin_terminal()
false
```
