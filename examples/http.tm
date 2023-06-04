#!/usr/bin/env tamarin

var body = json.marshal([1,2,3]).unwrap()

print("issuing post request to http://httpbin.org/post\n")

resp := fetch("https://httpbin.org/post", {
    method: "POST",
    timeout: 10.0,
    body: body,
    headers: {
        "Content-Type": "application/json",
    },
})

print(resp)

print("response:\n", resp.json().unwrap())
