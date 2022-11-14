# Tamarin Language Features

Here is an overview of Tamarin Language Features. This is not comprehensive.

## Print

Print to stdout:

```
print("Hello gophers!")
```

Print any number of variables:

```
print("x:", x, "y:", y)
```

Equivalent to `fmt.Println`.

## Assignment Statements

Both `let` and `const` statements are supported:

```
let x = 42
const y = "this is a constant"
```

Using the `:=` operator instead of `let` is encouraged:

```
x := 42
```

## Dynamic Typing

Variables may change type, similar to Python.

```
let x = 42
x = "now a string"
print(x)
```

## Semicolons

Semicolons are optional, so statements are ended by newlines if semicolons are not present.

```
let foo = "bar"; let baz = "qux"
```

## Comments

Lines are commented using `//`.

```
// This line is commented out
```

## Functions

Functions are defined using the `func` keyword. They may be passed around as values.
The `return` keyword is optional. If not present, the value of the last statement or
expression in the block is understood to be the return value. Expressions that do not
evaluate to a value will result in an `*object.Null` being returned.

```
func addOne(x) {
    x + 1
}

// This way of defining a function is equivalent to the above
const subOne = func(x) {
    return x - 1
}

addOne(100)
subOne(100)
```

Default parameter values are supported:

```
func increment(value, amount=1) {
    return value + amount
}

print(increment(100)) // 101
```

## Conditionals

Go style if-else statements are supported.

```
name := "ben"

if name == "noa" {
    print("the name is noa")
} else {
    print("the name is something else")
}
```

## Switch Statements

Go style switch statements are supported.

```
name := "ben"

switch name {
    case "ben":
        print("matched ben")
    case "noa":
        print("matched noa")
    default:
        print("default")
}
```

## Loops

Currently loops require parentheses but this will be updated to match Go's approach soon.

```
for (let i = 5; i < 10; i++) {
    print(i)
}
```

## Operations that may fail

`Result` objects wrap `Ok` and `Err` values for operations that may fail.

```
obj := json.unmarshal("true")
obj.unwrap() // returns true

failed := json.unmarshal("/not-valid/")
failed.is_err() // returns true
failed.unwrap() // raises error that stops execution
```

## Pipe Expressions

These execute a series of function calls, passing the result from one stage
in as the first argument to the next.

This pipe expression evalutes to the string `"50"`.

```
42 | math.sum(8) | string
```

## Array Methods

Arrays offer `map` and `filter` methods:

```
arr := [1, 2, 3, 4].filter(func(x) { x < 3 })
arr = arr.map(func(x) { x * x })
// arr is now [1, 4]
```

## Builtins

```
type(x)  // returns the string type name of x
len(s)   // returns the size of the string, array, hash, or set
any(arr) // true if any item in arr is truthy
all(arr) // true if all items in arr are truthy
```

## Types

A variety of built-in types are available.

```
101        // integer
1.1        // float
"1"        // string
[1,2,3]    // array
{1:2}      // hash
{1,2}      // set
false      // boolean
null       // null
func() {}  // function
time.now() // time
```

There are also `HttpResponse` and `DatabaseConnection` types in progress.

## Standard Library

Documentation for this is a work in progress. For now, browse the modules [here](../internal/modules).
