// This program demonstrates using Tamarin as a library
// to run a simple script.
package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/cloudcmds/tamarin/evaluator"
	"github.com/cloudcmds/tamarin/exec"
)

func main() {
	evaluator.RegisterPrintBuiltins()
	ctx := context.Background()

	input := `
		print("current time:", time.now())

		let a = [1, 4, 9]
		print("a:", a)

		let b = a.map(func(x) { math.sqrt(x) })
		print("b:", b)

		rand.shuffle(b)
		print("b shuffled:", b)

		print("uuid:", uuid.v4())

		// This will be the result value for the script
		json.marshal(b).unwrap()
	`

	result, err := exec.Execute(ctx, exec.Opts{Input: string(input)})
	if err != nil {
		fmt.Println("Execution error:", err)
		os.Exit(1)
	}

	fmt.Printf("script result: %s (type %v)\n",
		result.Inspect(), reflect.TypeOf(result))

	if result.Inspect() != "[1,3,2]" {
		fmt.Println("unexpected result")
		os.Exit(1)
	}
}
