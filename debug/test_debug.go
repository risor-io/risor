package main

import (
	"context"
	"fmt"

	"github.com/risor-io/risor"
)

func main() {
	result, err := risor.Eval(context.Background(), "1 + 1")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Error type: %T\n", err)
	} else {
		fmt.Printf("Result: %v\n", result)
	}
}