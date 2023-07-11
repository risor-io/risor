#!/usr/bin/env risor

d := {
    one: 1, 
    two: 2, 
    three: 3,
}

for k, v := range d {
    print("k:", k, "v:", v)
}
