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

Four forms of for loops are accepted. The `break` and `continue` keywords may
be used to control looping as you'd expect from other languages.

This form includes init, condition, and post statements:

```
for i := 0; i < 10; i++ {
    print(i)
}
```

This simple form will loop until a `break` is executed:

```
for {
    if a > b {
        break
    }
}
```

This form checks a condition before evaluating the loop body:

```
for a < b {

}
```

You may also use the `range` keyword to iterate through a container:

```
mylist := [1, 2, 3]
for index, value := range mylist { ... }
```

## Iterators

You can step through items in any container using an iterator. You can create
an iterator using the `iter` builtin function or by using the `range` keyword.

```
>>> iter({1,2,3})
set_iter({1, 2, 3})
```

```
>>> range {one: 1, two: 2}
map_iter({"one": 1, "two": 2})
```

Iterators offer a `next` method to retrieve the next entry in the sequence. Each
entry is returned as an `iter_entry` object, which has `key` and `value` attributes.
When the iterator is exhausted, `nil` is returned instead.

For loops work with these iterators natively, and automatically assign the
key and value to the loop variables. But you can use iterators directly as well.

```
>>> entry := range {foo: "bar"}.next()
iter_entry("foo", "bar")
>>> entry.key
"foo"
>>> entry.value
"bar"
```

## The "in" keyword

Check if an item exists is a container using the `in` keyword:

```
>>> 42 in [40, 41, 42]
true
>>> 3 in {2,3,4}
true
>>> "foo" in {foo: "bar"}
true
>>> "foo" in "bar foo baz"
true
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
all(arr)            // true if all items in arr are truthy
any(arr)            // true if any item in arr is truthy
assert(obj, msg)    // raises an error if obj is falsy
bool(obj)           // evaluates an object's truthiness
call(fn, ...)       // call the given function (can be useful in pipe expressions)
chr()               // convert an integer to its corresponding unicode rune
delete(map, key)    // delete an item from the map
err(message)        // create a Result error object
error(message)      // raise an error
float(s)            // convert a string to a float
getattr(obj, name)  // get the object's attribute with the given name
int(s)              // convert a string to an int
iter(obj)           // returns an iterator for the given container
keys(map)           // returns an array of keys in the given map
len(s)              // returns the size of the string, list, map, or set
list(obj)           // create a new list populated with items from the given iterable
ok(result)          // create a Result object containing the given object
ord()               // convert a unicode character to its integer value
print(...)          // equivalent to fmt.Println
printf(...)         // equivalent to fmt.Printf
reversed(arr)       // returns a reversed version of the given array
set(obj)            // create a new set populated with items from the given iterable
sorted(obj)         // return a sorted list of items from a container
sprintf(msg, ...)   // equivalent to fmt.Sprintf
string(obj)         // convert an object to its string representation
try(expr, fallback) // evaluate expr and return fallback if an error occurs
type(x)             // returns the string type name of x
unwrap_or(obj)      // unwraps but returns the provided obj if the Result is an Error
unwrap(result)      // unwraps the ok value from the Result if allowed
```

## Types

A variety of built-in types are available.

```
101         // int
1.1         // float
"1"         // string
[1,2,3]     // list
{"key":2}   // map
{1,2}       // set
false       // bool
nil         // nil
func() {}   // function
time.now()  // time
iter([1])   // list_iter
iter(set()) // set_iter
iter({})    // map_iter
iter("1")   // string_iter
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
