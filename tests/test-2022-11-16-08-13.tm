// github issue: n/a
// expected value: '"hello"'
// expected type: string

var s = "\"hello\""
var j = json.unmarshal(s)
assert(type(j) == "string")
assert(j == "hello")

s = json.marshal("hello")
assert(type(s) == "string")
assert(s == "\"hello\"")

s
