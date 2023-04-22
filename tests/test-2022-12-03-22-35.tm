// expected value: "101"
// expected type: string

inc := func(x) {
    x + 1
}

i := 100

s := '{inc(i)}'
s
