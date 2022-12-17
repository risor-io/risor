
s := {1, 2, 3}

func test(item, x=10) {
    return len(item) + x
}

// Just a test for using call with a user-defined function.
// Ordinarily you would just call the function directly like "test(s)".
result := call(test, s)

print(result)
