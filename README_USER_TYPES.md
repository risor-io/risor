# User-Defined Types in Risor

This document describes the new user-defined types feature in Risor, which brings TypeScript-inspired syntax for type declarations, interfaces, and method definitions to the language.

## Overview

Risor now supports first-class user-defined types with the following features:

- **Type declarations** with field specifications
- **Interface declarations** for defining contracts
- **Type annotations** on variables and function parameters
- **Method receivers** for associating functions with types
- **Return type annotations** for functions

## Type Declaration Syntax

### Basic Type Declaration

```risor
type Person {
    name: string,
    age: int,
    email: string
}
```

### Complex Types with Nested Structures

```risor
type Company {
    name: string,
    employees: []Person,
    address: Address
}

type Address {
    street: string,
    city: string,
    zipCode: string
}
```

## Interface Declaration Syntax

Interfaces define method signatures that types can implement:

```risor
interface Drawable {
    draw(): void,
    getArea(): float,
    setPosition(x: float, y: float): void
}
```

## Variable Type Annotations

### Explicit Variable Declarations

```risor
var name: string = "Alice"
var age: int = 30
var scores: []float = [95.5, 87.2, 92.1]
```

### Short Declarations with Type Annotations

```risor
name: string := "Alice"
age: int := 30
scores: []float := [95.5, 87.2, 92.1]
```

## Method Receivers (Go-like Syntax)

Functions can be associated with types using method receivers:

```risor
func (p Person) greet(): string {
    return "Hello, I'm " + p.name
}

func (p Person) getAge(): int {
    return p.age
}

func (p Person) celebrate(): void {
    p.age = p.age + 1
}
```

## Function Return Type Annotations

Functions can specify their return types:

```risor
func NewPerson(name: string, age: int): Person {
    return Person{
        name: name,
        age: age,
        email: name + "@example.com"
    }
}

func calculate(x: float, y: float): float {
    return x * y + 10.0
}
```

## Complete Example

```risor
// Type declaration
type Rectangle {
    width: float,
    height: float
}

// Interface declaration
interface Shape {
    getArea(): float,
    getPerimeter(): float
}

// Constructor function
func NewRectangle(w: float, h: float): Rectangle {
    return Rectangle{width: w, height: h}
}

// Methods with receivers
func (r Rectangle) getArea(): float {
    return r.width * r.height
}

func (r Rectangle) getPerimeter(): float {
    return 2 * (r.width + r.height)
}

func (r Rectangle) describe(): string {
    return "Rectangle " + string(r.width) + "x" + string(r.height)
}

// Usage
var rect: Rectangle = NewRectangle(10.0, 5.0)
area: float := rect.getArea()
description: string := rect.describe()

print(description)
print("Area:", area)
```

## Design Decisions

### TypeScript Inspiration

The syntax draws heavily from TypeScript for several reasons:

1. **Familiar syntax** for developers coming from web development
2. **Clear type annotations** with the `:` separator
3. **Optional typing** that doesn't interfere with existing code
4. **Interface-based design** for flexible type contracts

### Go-like Method Receivers

Method receivers follow Go's syntax `func (receiver Type) method()` because:

1. **Explicit receiver naming** avoids confusion about `this` binding
2. **Clear association** between methods and types
3. **No hidden variables** - the receiver is always explicitly named
4. **Familiar to Go developers** who might use Risor

### Avoiding Parsing Conflicts

The design carefully avoids the parsing ambiguity mentioned in the discussion:

- **No type prefixes on literals** like `[]int[1, 2, 3]`
- **Type annotations only on variables** using `: type` syntax
- **Clear token boundaries** that don't require lookahead parsing

## Implementation Notes

### AST Extensions

New AST nodes were added:

- `TypeDecl` - for type declarations
- `InterfaceDecl` - for interface declarations  
- `TypeField` - for fields in type declarations
- `InterfaceMethod` - for method signatures in interfaces
- `TypeAnnotation` - for type annotations on variables
- `MethodReceiver` - for method receivers on functions

### Parser Extensions

The parser was extended to:

- Recognize `type` and `interface` keywords
- Parse type declarations with field lists
- Parse interface declarations with method signatures
- Handle type annotations on variable declarations
- Support method receivers in function declarations
- Parse return type annotations

### Token Extensions

New tokens were added:

- `TYPE` - for the `type` keyword
- `INTERFACE` - for the `interface` keyword

## Future Enhancements

Potential future improvements could include:

1. **Generic types** - `type List<T> { items: []T }`
2. **Union types** - `type StringOrNumber = string | int`
3. **Type aliases** - `type UserId = string`
4. **Type inference** - automatically inferring types from usage
5. **Runtime type checking** - validating types at runtime

## Backward Compatibility

This feature is fully backward compatible:

- Existing Risor code continues to work unchanged
- Type annotations are optional
- No breaking changes to existing syntax
- Progressive adoption is possible