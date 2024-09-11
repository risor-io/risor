import { Callout } from 'nextra/components'

# fmt

The `fmt` module provides functions for formatting and printing strings.

Usually you should use Risor's top-level `print` and `printf` functions instead
of using `fmt.print` and `fmt.printf` since they are equivalent.


## Functions

### errorf

```go filename="Function signature"
errorf(string, ...any) error
```

Returns a new error with the given message formatted according to the format.

```go filename="Example"
>>> err := fmt.errorf("something went wrong: %d", 42)
>>> err
something went wrong: 42
```

### printf

```go filename="Function signature"
printf(string, ...any)
```

Prints the formatted string to the standard output.

```go filename="Example"
>>> fmt.printf("Hello, %s!\n", "world")
Hello, world!
```

### print

```go filename="Function signature"
print(...any)
```

Prints the given values to the standard output. Note that in Risor the output
may not be printed to the terminal until a newline character is printed.

```go filename="Example"
>>> fmt.print("Hello, ", "world", "!\n")
Hello, world!
```
