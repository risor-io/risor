#!/usr/bin/env tamarin

sum := 0

for i := 0; i < 10; i++ {
    if i == 4 { break }
    sum += i
}

print("sum:", sum)
assert(sum == 6, "expected sum to be 6 (0 + 1 + 2 + 3)")
