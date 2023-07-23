#!/usr/bin/env risor --

x := 0

// The first function should fail on the error, right after x is set to 10.
// After that, the second function will run and cause 33 to be returned from
// the try call.
result := try(func() {
    x = 10
    error("kaboom")
    x = 11
}, func() {
    33
})

print("result:", result)
print("x:", x)
