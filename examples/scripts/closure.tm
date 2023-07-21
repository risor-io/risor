#!/usr/bin/env risor --

func incrementer(value=0) {
    func inner() {
        result := value
        value++
        return result
    }
    return inner
}

inc := incrementer()
print(inc())
print(inc())
print(inc())
print(inc())
