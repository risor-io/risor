// github issue: n/a
// expected value: null
// expected type: NULL

id := uuid.v4()
assert(type(id) == "string")
assert(strings.count(id, "-") == 4)
assert(len(id) == 36)