// github issue: n/a
// expected value: [1, 2, 3]
// expected type: list

a := [
    "1",
    "22",
    "333",
]

a.map(func(x) { len(x) } )
