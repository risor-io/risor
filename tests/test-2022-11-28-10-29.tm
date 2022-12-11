// expected value: [1, 2, 2.2, 3, a]
// expected type: LIST

s1 := {1, "a", 2.2}
s2 := {2, "a", 2.2}

s1.add(3)

union := s1.union(s2) | sorted
