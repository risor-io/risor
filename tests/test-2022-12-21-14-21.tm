// expected value: {1, 2, 3, 4}
// expected type: set

s := {1, 2, 3}

assert(s[1])
assert(!s[99])

updated := s.union({4})

updated
