# Syntax

Tamarin syntax was designed to be like Go, but with some functionality and
behavior borrowed from Python, given that Tamarin is an interpreted language.

In the examples below, when you see a `>>>` it indicates that input and output
from a Tamarin REPL session is being shown. To the right of the `>>>` is what
the user entered, and the output is shown on the line below.

## Variables

Variables are dynamically typed and are declared using `const`, `var`, or `:=`.
After declaration, variables are updated using `=` like in Go.

```go
x := 42             // this declares an integer
x = "now a string"  // like Python, a variable may change types
const foo = "bar"   // const values cannot be updated
var name = "anne"   // this is equivalent to `name := "anne"`
```

## Semicolons

Semicolons are optional. A statement ends on a semicolon if a newline is not present,
except in some specific situations.

```go
foo := "bar"; baz := "qux"
```

## Print

Print to stdout:

```go
print("Hello gophers!")
```

Print any number of variables:

```go
print("x:", x, "y:", y)
```

The `print` built-in is equivalent to `fmt.Println` in Go.

## Comments

Lines are commented using `//`.

```go
// This line is commented out
```

Block comments are defined using `/*` and `*/`.

## Ints and Floats

Tamarin `int` and `float` types correspond to `int64` and `float64` in Go. These
values are boxed as Tamarin objects. Generally, Tamarin libraries convert bewtween
`int` and `float` automatically in numeric operations.

```go
>>> 1 + 3.3
4.3
>>> type(1)
int
>>> type(2.0)
float
>>> math.max([1, 2.0])
2
```

## Strings

String behavior overall is quite similar to that in Go, but single quoted
string templates are available to provide templating similar to Python f-strings.
Raw strings are defined using backticks and these opt-out of all escape character
behaviors.

```js
salutation := "hello there"        // double-quoted string
count_str := 'the count is {1+1}'  // single-quoted string with template variable
raw_str := `\t\r\n`                // raw string
```

## Functions

Functions are defined using the `func` keyword. They may be passed around as values.
The `return` keyword is optional. If not present, the value of the last statement or
expression in the block is understood to be the return value. Expressions that do not
evaluate to a value will result in an `*object.Nil` being returned.

```go
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

```go
func increment(value, amount=1) {
    return value + amount
}

print(increment(100)) // 101
```

## Conditionals

Go style if-else statements are supported.

```go
name := "ben"

if name == "noa" {
    print("the name is noa")
} else {
    print("the name is something else")
}
```

## Switch Statements

Go style switch statements are supported.

```go
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

```go
for i := 0; i < 10; i++ {
    print(i)
}
```

This simple form will loop until a `break` is executed:

```go
for {
    if a > b {
        break
    }
}
```

This form checks a condition before evaluating the loop body:

```go
for a < b {

}
```

You may also use the `range` keyword to iterate through a container:

```go
mylist := [1, 2, 3]
for index, value := range mylist { ... }
```

## Iterators

You can step through items in any container using an iterator. You can create
an iterator using the `iter` builtin function or by using the `range` keyword.

```go
>>> iter({1,2,3})
set_iter({1, 2, 3})
```

```go
>>> range {one: 1, two: 2}
map_iter({"one": 1, "two": 2})
```

Iterators offer a `next` method to retrieve the next entry in the sequence. Each
entry is returned as an `iter_entry` object, which has `key` and `value` attributes.
When the iterator is exhausted, `nil` is returned instead.

For loops work with these iterators natively, and automatically assign the
key and value to the loop variables. But you can use iterators directly as well.

```go
>>> entry := range {foo: "bar"}.next()
iter_entry("foo", "bar")
>>> entry.key
"foo"
>>> entry.value
"bar"
```

## The "in" keyword

Check if an item exists is a container using the `in` keyword:

```go
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

```go
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

```go
"hello" | strings.to_upper
```

## List Methods

Lists offer `map` and `filter` methods:

```go
list := [1, 2, 3, 4].filter(func(x) { x < 3 })
list = list.map(func(x) { x \* x })
// list is now [1, 4]
```

## Types

A variety of built-in types are available.

```go
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

```go title="proxy_service.go" linenums="1"
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
