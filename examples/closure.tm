#!/usr/bin/env tamarin

var testValue = 100

func getint() {
    var foo = testValue + 1
    func inner() {
        foo
    }
    return inner
}

print(getint()())
