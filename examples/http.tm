#!/usr/bin/env tamarin

resp := fetch("https://httpbin.org/post", {
    method: "POST",
    timeout: 10000,
    data: {foo: "bar"},
    params: {test: "123"},
})

print(resp)
print()
print(resp.json())
