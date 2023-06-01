
x := 0
y := 0

result := try(func() {
    x = 10
    err("SMASH")
    x = 11
    22
}, func() {
    y = 20
    y = 21
    33
}, "nope")

print("HI", result, x, y, type(result))
