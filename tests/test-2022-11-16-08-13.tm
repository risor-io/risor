// github issue: n/a
// expected value: '"hello"'
// expected type: string

let s = "\"hello\""
let j = json.unmarshal(s)
assert(type(j) == "result")
assert(j.is_ok())
assert(!j.is_err())

let v = j.unwrap()
assert(type(v) == "string")
assert(v == "hello")

s = json.marshal(v)
assert(type(s) == "result")
assert(j.is_ok())
assert(!j.is_err())

assert(s.unwrap() == "\"hello\"")
s.unwrap()
