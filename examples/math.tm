#!/usr/bin/env tamarin

// Simple function definition. Note that the `return` keyword is optional
// and if it's not present the last expression in the body is the return value.
func square(n) { n * n }

let ints = [0, 1, 2, 3, 4]

print("integers:", ints)

print("squares:", ints.map(square))

print("sum of integers:", math.sum(ints))

print("math.PI:", math.PI)
