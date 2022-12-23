#!/usr/bin/env tamarin

var body = json.marshal([1,2,3]).unwrap()

print("issuing post request to http://httpbin.org/post\n")

resp := fetch("https://httpbin.org/post", {
    method: "POST",
    timeout: 0.001,
    body: body,
    headers: {
        "Content-Type": "application/json",
    },
}).json().unwrap()

print("response:\n", resp)
