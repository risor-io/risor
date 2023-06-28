#!/usr/bin/env risor

// sum := 0

// for i := 0; i < 10; i++ {
//     if i == 4 { break }
//     sum += i
// }

// print("sum:", sum)
// assert(sum == 6, "expected sum to be 6 (0 + 1 + 2 + 3)")

// for i := range [0, 10] {
//     if i == 4 { continue }
//     sum += i
// }

x, y := [99, 100]
print("x:", x, "y:", y)

// r := range [41, 42, 43]
d := {one: 1, two: 2, three: 3}

for k, v := range d {
    print("k:", k, "v:", v)
}

