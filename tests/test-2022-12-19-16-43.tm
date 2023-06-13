// expected value: 10
// expected type: int

func square(x) { x * x }

assert(square(2) == 4)

x := 10

// This confirms the temporary scope for function execution, which also uses
// a variable named `x` doesn't update the outer scope's `x` variable.
x
