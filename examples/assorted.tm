#!/usr/bin/env tamarin

var testValue = 100

// comment
func getint() {
    var foo = testValue + 1
    func inner() {
        foo
    }
    return inner
}

func inc(x) {
    x + 1
}

print(getint()() )

// var x = getint()() | inc
var x = 4

print(x)

var data = {
    "foo": "bar",
    "value": 42
}

print("data-foo:", data["foo"])

a := [1, 2, 3, 4, 5]
print(a[4])
print("length of a:", len(a))

func addOne(index, value) {
    return value + 1
}

var mapped = a.map(type)
print("mapped:", mapped, type(mapped[0]))

s := { 1, 3.3,
  -1, -99, 42, 101.3, 3.3, 99 }
print("S:", s)

import library
print("library call:", library.add(1, 2))

print(
    "type:", type(s),
    "len:", len(s),
    "max:", math.max(s),
    "min:", math.min(s),
)

var filtered = [1, 2, 3, 4, 5].filter(func(item) { item < 3 })
print("filtered:", filtered)

print("ok:", ok("yup").unwrap(), ok("yup").is_ok())
err("explosion").is_err()

v := json.unmarshal("{\"a\":\"b\"}")
print("v:", v.unwrap(), type(v))

b := json.unmarshal("true")

switch b.is_err() {
    case true:
        print("unmarshal failed")
    case 42, false:
        print("unmarshal ok")
    case 5:
        print("case 5")
}

print('b is {b}')

assert(true)

resp := fetch("https://httpbin.org/post", {
    "method": "POST",
    "timeout": 1.0,
    "body": "42",
    "headers": {
        "Content-Type": "application/json",
    },
}).json().unwrap()

print(resp, type(resp))

print("uuid:", uuid.v4())

[1, 1, 2, 3, 4, 4] | set | string

print(strings.join(["hey", "ya"], "_"))

print("the time is:", time.now())

print("json diff:", json.diff(["a"],["b"]))

print("pipe expr:",
    ["foo", "bar"] | strings.join("-"),
    ["a", "c"] | strings.join("b") | strings.contains("c"),
    [99, 98, 97] | math.max)

print("any?", any([0, false, []]))

{ "results": "\"nice!\"" | json.unmarshal, "pi": math.PI }
