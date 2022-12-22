// expected value: [1, 2, 3, 4, 5]
// expected type: list

l := [1, 2, 3]

funcs := [l.append]

funcs[0](4)
funcs[0](5)

l
