# Syntax

Tamarin was designed to be like a more scripting-friendly version of Go.
At times, approaches from Python were referenced when deciding how Tamarin
should approach a particular situation as an interpreted language. As a result,
Tamarin may feel like a hybrid of Go and Python.

!!! Note

    In the examples below, when you see a `>>>` it indicates that input and output
    from a Tamarin REPL session is being shown. To the right of the `>>>` is the
    command the user entered. The command output is shown on the line below.

## Variables

Variables are dynamically typed and are declared using `const`, `var`, or `:=`.
After declaration, variables are updated using `=` like in Go.

```go
x := 42             // this declares an integer
x = "now a string"  // like Python, a variable may change types
const foo = "bar"   // const values cannot be updated
var name = "anne"   // this is equivalent to `name := "anne"`
```

Multiple variables may be assigned in one statement, where the right-hand
side of the assignment is a list with a matching size:

```go
>>> a, b, c := [1, 2, 3]
[1, 2, 3]
>>> a
1
>>> b
2
>>> c
3
```

## Semicolons

Semicolons are optional. Multiple statements can be on a single line if
separated by a semicolon.

```go
foo := "bar"; baz := "qux"
```

## Comments

Lines are commented using `//` or `#`.

```go
// This line is commented out
```

```python
# As is this one
```

Block comments are defined using `/*` and `*/`.

## Functions

Functions are defined using the `func` keyword and may be passed around as values.
The last statement or expression in the function body is understood to be the
return value, so using the `return` keyword is _optional_.

The syntax for invoking a function is the same as in Go.

Functions may be declared with default parameter values.

```go
>>> func increment(x, amount=1) { x + amount }
>>> increment(3)
4
>>> increment(3, 2)
5
```

Functions may also be assigned to variables:

```go
const say_hello = func() {
    print("hello")
}

say_hello()
```

## Closures

Closures store the environment associated with an outer function, allowing its
reuse for multiple invocations of an inner function.

```go
func get_counter(start) {
  return func() {
    start++
  }
}

c := get_counter(5)

print(c()) // prints 5
print(c()) // prints 6
print(c()) // prints 7
```

## Conditionals

Go style conditionals are supported, including `if`, `else if`, and `else` cases.
Parentheses are not required when defining the condition for each case.

```go
if age < 13 {
    print("this is a kid")
} else if age < 18 {
    print("this is a young adult")
} else {
	print("this is an adult")
}
```

## Switch Statements

Switch statements compare a value with multiple defined cases, executing the
matching case if there is one, and executing the `default` case if one exists
and no other cases match.

```go
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

Multiple styles of for loops are accepted. The `break` and `continue` keywords
may be used to control looping as you'd expect from other languages.

The "standard" style includes _init_, _condition_, and _post_ statements:

```go
for i := 0; i < 10; i++ {
    print(i)
}
```

The "simple" style will loop until a `break` statement is reached:

```go
for {
    if a > b {
        break
    }
}
```

The "condition-only" style loops until the condition evaluates to `false`:

```go
for a < b {
	a++
}
```

Finally, you may also use the `range` keyword to iterate over the contents of a container:

```go
l := [1, 2, 3]

for i, value := range l {
	print(i, value)
}
```

## Pipelines

Pipelines execute a series of function calls, passing each call's output as the
first argument to the next call. This syntax lends itself to applying one or more
transformations to input data.

If an error is encountered at any stage in the pipeline, its execution stops early,
and the error propagates as usual.

Each expression in the pipeline is expected to evaluate to a function or method
to call. Parentheses may be omitted in each call when the function accepts one
argument, since that argument is passed implicitly from the previous stage. If
the function accepts two or more arguments, the pipeline always provides the
first argument and the code author must supply the following arguments as they
would in a normal function invocation.

This is an example of a two stage string transformation:

```go
>>> "hello." | strings.to_upper | strings.replace_all(".", "!")
"HELLO!"
```

The expression prior to the first `|` receives no special treatment in pipelines.
That is, it's treated as the first argument to the subsequent function, even if
it evaluates to a function value itself.

The examples below are all equivalent and illustrate how values (which may even
be a function) are passed as the first argument to the next stage.

```go
>>> [3, 2, 9, 5] | math.max
9
>>> math.max | call([3, 2, 9, 5])
9
>>> call(math.max, [3, 2, 9, 5])
9
```

Pipelines can be used to build functions:

```go
>>> func normalize(s) { s | strings.fields | strings.join(" ") }
>>> normalize("  so   much   whitespace  ")
"so much whitespace"
```

## Attributes

Objects in Tamarin may have data attributes and method attributes. Both are
retrieved with a familiar `object.attribute` syntax. There is also a built-in
`getattr` function that supports retrieving a named attribute from an object.

A method accessed on an object remains bound to that object, as you can see here:

```go
>>> l := [1,2,3]
[1, 2, 3]
>>> append := l.append
builtin(list.append)
>>> append(4)
[1, 2, 3, 4]
>>> l
[1, 2, 3, 4]
>>> getattr(l, "append")
builtin(list.append)
```

## Indexing

Multiple Tamarin core types support index operations. These types are referred
to as _container types_ and include list, map, set, and string.

Lists use zero-based indexing. Negative indices operate relative to the end of the list:

```go
>>> l := ["a", "b", "c"]
["a", "b", "c"]
>>> l[2]
"c"
>>> l[-1]
"c"
```

Map indexing is used to retrieve the value for a given key:

```go
>>> m := {name:"evan", age:18}
{"age": 18, "name": "evan"}
>>> m["name"]
"evan"
```

Set indexing is used to check whether a value is present in the set:

```go
>>> s := {3,4,5}
{3, 4, 5}
>>> s[5]
true
>>> s[6]
false
```

String indexing is used to read unicode code points, known as _runes_ in Go,
from a given string. Note this behavior differs from Go string indexing, which
operates on the underlying bytes.

```go
>>> const s = "正體字"
"正體字"
>>> s[0]
"正"
>>> s[1]
"體"
```

## Slices

Lists and strings in Tamarin support slice operations to select a range of items.

```go
>>> l := ["a", "b", "c"]
["a", "b", "c"]
>>> l[0:2]
["a", "b"]
>>> l[-2:]
["b", "c"]
```

The syntax for this is `l[start:stop]` where `start` and `stop` may be omitted
in order to refer to the beginning or the end of the sequence, respectively.

## Import

Tamarin files may be imported as modules using the `import` keyword. All module
data and functions are available as attributes on the module after import. As a
convention, if an attribute is intended to be private to a module, prefix its
name with an underscore.

```go
>>> import library
module(library.tm)
>>> library.add(2,3)
5
```

## The in Keyword

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
