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

Both `var` and `const` statements are supported:

```
var x = 42
const y = "this is a constant"
```

Using the `:=` operator instead of `var` is encouraged:

```
x := 42
```

## Dynamic Typing

Variables may change type, similar to Python.

```
x := 42
x = "now a string"
print(x)
```

## Optional Semicolons

Semicolons are optional, so statements are ended by newlines if semicolons are not present.

```
foo := "bar"; baz := "qux"
```

## Comments

Lines are commented using `//`.

```
// This line is commented out
```

Block comments are defined using `/*` and `*/`.

## Integers

Integers are represented by Go's `int64` type internally. Numeric operations
that involve a mix of int and float inputs generally produce a float output.
Integers are automaticaly converted to floats in these situations.

```
>>> 1 + 3.3
4.3
>>> type(1)
int
```

## Floats

Floating point numbers use Go's `float64` type internally.

```
>>> math.max([1.0, 2.0])
2.0
>>> type(2.0)
float
```

## Strings

Strings come in three varieties. The standard string uses double quotes and
behaves very similarly to Go's string. Single quoted strings are similar, but
have the additional feature of supporting string templating similar to Python's
f-strings. Finally, raw strings are defined using backticks. Use raw strings
when you want to include backslashes, single or double quotes, or newlines in
your string.

```
salutation := "hello there"        // double-quoted string
count_str := 'the count is {1+1}'  // single-quoted string
raw_str := `\t\r\n`                // raw string
```

## Functions

Functions are defined using the `func` keyword. They may be passed around as values.
The `return` keyword is optional. If not present, the value of the last statement or
expression in the block is understood to be the return value. Expressions that do not
evaluate to a value will result in an `*object.Nil` being returned.

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

Two forms of for loops are accepted. The `break` keyword may be used to
stop looping in either form.

```
for i := 0; i < 10; i++ {
    print(i)
}

for {
    if condition {
        break
    }
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

This pipe expression evalutes to the string `"HELLO"`.

```
"hello" | strings.to_upper
```

## List Methods

Lists offer `map` and `filter` methods:

```
list := [1, 2, 3, 4].filter(func(x) { x < 3 })
list = list.map(func(x) { x \* x })
// list is now [1, 4]
```

## Builtins

```
type(x)            // returns the string type name of x
len(s)             // returns the size of the string, list, map, or set
any(arr)           // true if any item in arr is truthy
all(arr)           // true if all items in arr are truthy
sprintf(msg, ...)  // equivalent to fmt.Sprintf
keys(map)          // returns an array of keys in the given map
delete(map, key)   // delete an item from the map
string(obj)        // convert an object to its string representation
bool(obj)          // evaluates an object's truthiness
ok(result)         // create a Result object containing the given object
err(message)       // create a Result error object
unwrap(result)     // unwraps the ok value from the Result if allowed
unwrap_or(obj)     // unwraps but returns the provided obj if the Result is an Error
sorted(obj)        // works with lists, maps, and sets
reversed(arr)      // returns a reversed version of the given array
assert(obj, msg)   // raises an error if obj is falsy
print(...)         // equivalent to fmt.Println
printf(...)        // equivalent to fmt.Printf
set(obj)           // create a new set populated with items from the given iterable
list(obj)          // create a new list populated with items from the given iterable
int(s)             // convert a string to an int
float(s)           // convert a string to a float
call(fn, ...)      // call the given function (can be useful in pipe expressions)
getattr(obj, name) // get the object's attribute with the given name
ord()              // convert a unicode character to its integer value
chr()              // convert an integer to its corresponding unicode rune
```

## Types

A variety of built-in types are available.

```
101        // int
1.1        // float
"1"        // string
[1,2,3]    // list
{"key":2}  // map
{1,2}      // set
false      // bool
nil        // nil
func() {}  // function
time.now() // time
```

There are also `HttpResponse` and `DatabaseConnection` types in progress.

## Standard Library

Documentation for this is a work in progress. For now, browse the modules [here](../internal/modules).

## Proxying Calls to Go Objects

You can expose arbitrary Go objects to Tamarin code in order to enable method
calls on those objects. This allows you to expose existing structs in your
application as Tamarin objects that scripts can be written against. Tamarin
automatically discovers public methods on your Go types and converts inputs and
outputs for primitive types and for structs that you register.

Input and output values are type-converted automatically, for a variety of types.
Go structs are mapped to Tamarin map objects. Go `context.Context` and `error`
values are handled automatically.

```go
	// Create a registry that tracks proxied Go types and their attributes
	registry, err := object.NewTypeRegistry()
	if err != nil {
		return err
	}

	// This is the Go service we will expose in Tamarin
	svc := &MyService{}

	// Wrap the service in a Tamarin Proxy
	proxy, err := object.NewProxy(registry, svc)
	if err != nil {
		return err
	}

	// Add the proxy to a Tamarin execution scope
	s := scope.New(scope.Opts{})
	s.Declare("svc", proxy, true)

	// Execute Tamarin code against that scope. By doing this, the Tamarin
	// code can call public methods on `svc` and retrieve its public fields.
	result, err := exec.Execute(ctx, exec.Opts{
		Input: string(scriptSourceCode),
		Scope: s,
	})
```

See [example-proxy](../cmd/example-proxy/main.go) for a complete example.
