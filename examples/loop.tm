#!/usr/bin/env tamarin

sum := 0

for i := 0; i < 4; i++ {
    sum += i
}

print("sum:", sum)
assert(sum == 6, "expected sum to be 6 (1 + 2 + 3)")
