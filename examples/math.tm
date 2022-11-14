#!/usr/bin/env tamarin

// Simple function definition. Note that the `return` keyword is optional
// and if it's not present the last expression in the body is the return value.
func square(n) { n * n }

let ints = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

print("integers:", ints)

print("squares:", ints.map(square))

print("sum:", math.sum(ints))

import library
print("addition example via a library:", library.add(1, 2))

print("math.PI:", math.PI)
