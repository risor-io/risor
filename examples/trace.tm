
func checkSomething(value, whatever="foo") {
    if value > 100000 {
        print("Wow that's big")
    }
}

func add(x, y) {
    result := x + y
    checkSomething(result)
    return result
}

func sub(x, y) {
    result := x - y
    checkSomething(result)
    return result
}

func doWork(inputs) {
    var sum = 0.0
    for i := 0; i < len(inputs); i++ {
        x := inputs[i]
        r := rand.float()
        y := float()
        if r < 0.5 {
            y = add(x, r)
        } else {
            y = sub(x, r)
        }
        sum += y
    }
    return sum
}

inputs := []
for i := 0; i < 1000; i++ {
    inputs.append(rand.float())
}

result := doWork(inputs)
print("result:", result)

result
