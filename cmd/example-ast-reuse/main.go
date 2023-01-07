// This program demonstrates using Tamarin as a library
// to run a simple script.
package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/cloudcmds/tamarin/core/exec"
	"github.com/cloudcmds/tamarin/core/object"
	"github.com/cloudcmds/tamarin/core/parser"
	"github.com/cloudcmds/tamarin/core/scope"
)

func main() {

	ctx := context.Background()

	// This is the Tamarin script supplied by the user. It assumes a variable
	// named `input` is available in the execution scope.
	userCode := "sorted(input)"

	// Parse the user code into an AST.
	program, err := parser.ParseWithOpts(ctx, parser.Opts{
		Input: userCode,
	})
	if err != nil {
		fmt.Println("Parse error:", err)
		os.Exit(1)
	}

	// Create multiple scopes, one for each execution. Set them up with a
	// different `input` value.
	inputs := []*scope.Scope{
		scope.New(scope.Opts{Name: "execution1-scope"}),
		scope.New(scope.Opts{Name: "execution2-scope"}),
	}
	inputs[0].Declare("input", object.NewStringList([]string{"b", "a", "c"}), false)
	inputs[1].Declare("input", object.NewStringList([]string{"z", "y", "x"}), false)

	// Execute the same AST multiple times, with a different scope each time.
	for i, input := range inputs {
		result, err := exec.Execute(ctx, exec.Opts{
			InputProgram: program,
			Scope:        input,
		})
		if err != nil {
			fmt.Println("Execution error:", err)
			os.Exit(1)
		}
		fmt.Printf("execution %d result: %s (type %v)\n",
			i, result.Inspect(), reflect.TypeOf(result))
	}
}
