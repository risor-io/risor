# User-Defined Types in Risor

This document describes the new user-defined types feature in Risor, which brings TypeScript-inspired syntax for type declarations, interfaces, and method definitions to the language.

## ✅ Currently Working Features

Risor now supports the following type system features:

- **Type declarations** with field specifications
- **Interface declarations** for defining contracts  
- **Type annotations** on variables (both `var` and `:=` declarations)
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

### Short Declarations with Type Annotations (Walrus Operator)

```risor
name: string := "Alice"
age: int := 30
scores: []float := [95.5, 87.2, 92.1]
```

## Function Return Type Annotations

Functions can specify their return types:

```risor
func NewPerson(name, age): Person {
    return Person{
        name: name,
        age: age,
        email: name + "@example.com"
    }
}

func calculate(x, y): float {
    return x * y + 10.0
}
```

## Working Example

```risor
// Type declarations
type User {
    id: int,
    username: string,
    email: string
}

// Interface declarations
interface Identifiable {
    getId(): int,
    getUsername(): string
}

// Variable declarations with type annotations
var adminUser: string = "admin"
var maxUsers: int = 1000
var systemActive: bool = true

// Walrus declarations with type annotations
currentUser: string := "alice"
userCount: int := 150
serverLoad: float := 67.5

// Functions with return type annotations
func getUserCount(): int {
    return userCount
}

func formatUsername(name): string {
    return "[" + name + "]"
}

// Using the functions
formattedUser: string := formatUsername(currentUser)
totalUsers: int := getUserCount()

print("Current user:", formattedUser)
print("Total users:", totalUsers)
print("System status:", systemActive)
```

## 🚧 Planned Features (Not Yet Implemented)

The following features are designed but not yet implemented:

### Method Receivers

```risor
// Planned syntax (not working yet)
func (p Person) greet(): string {
    return "Hello, I'm " + p.name
}
```

### Function Parameter Type Annotations

```risor
// Planned syntax (not working yet)  
func add(x: int, y: int): int {
    return x + y
}
```

### Struct Literal Construction

```risor
// Planned syntax (not working yet)
var person: Person = Person{
    name: "Alice",
    age: 30,
    email: "alice@example.com"
}
```

### Runtime Type Checking

Currently, type annotations are parsed and stored but not enforced at runtime. Future versions will include:

- Type validation during assignment
- Type checking for function calls
- Interface implementation verification

## Design Decisions

### TypeScript Inspiration

The syntax draws heavily from TypeScript for several reasons:

1. **Familiar syntax** for developers coming from web development
2. **Clear type annotations** with the `:` separator
3. **Optional typing** that doesn't interfere with existing code
4. **Interface-based design** for flexible type contracts

### Go-like Method Receivers (Planned)

Method receivers will follow Go's syntax `func (receiver Type) method()` because:

1. **Explicit receiver naming** avoids confusion about `this` binding
2. **Clear association** between methods and types
3. **No hidden variables** - the receiver is always explicitly named
4. **Familiar to Go developers** who might use Risor

### Avoiding Parsing Conflicts

The design carefully avoids the parsing ambiguity mentioned in the original discussion:

- **No type prefixes on literals** like `[]int[1, 2, 3]`
- **Type annotations only on variables** using `: type` syntax
- **Clear token boundaries** that don't require lookahead parsing

## Implementation Status

### ✅ Completed

- Token support for `type` and `interface` keywords
- AST nodes for all type system constructs
- Parser support for type declarations and interface declarations
- Parser support for variable type annotations (both `var` and `:=`)
- Parser support for function return type annotations
- Compiler support (no-op compilation for now)
- Comprehensive test suite
- Working examples and documentation

### 🔄 In Progress

- Method receiver parsing (complex due to lookahead requirements)
- Function parameter type annotations
- Struct literal construction syntax

### 📋 Planned

- Runtime type checking and validation
- Type inference capabilities
- Generic types support
- Union types support

## Testing

The implementation includes comprehensive tests covering:

- Type declaration parsing
- Interface declaration parsing
- Variable type annotation parsing (both `var` and `:=`)
- Function return type annotation parsing
- Integration with existing Risor features

Run the tests with:

```bash
go test ./parser -v -run "TestType|TestInterface|TestVariable.*TypeAnnotation|TestWalrus|TestFunction.*ReturnType"
```

## Examples

See the working examples in:
- `examples/user_types_working.risor` - Comprehensive demonstration
- `test_comprehensive.risor` - Feature testing

## Backward Compatibility

This feature is fully backward compatible:

- Existing Risor code continues to work unchanged
- Type annotations are optional
- No breaking changes to existing syntax
- Progressive adoption is possible