// github issue: n/a
// expected value: nil
// expected type: nil

id := uuid.v4()
assert(type(id) == "string")
assert(strings.count(id, "-") == 4)
assert(len(id) == 36)
