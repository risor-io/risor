package main

import (
	"context"
	"fmt"
	"log"

	"github.com/risor-io/risor"
)

func main() {
	// Test basic traceback functionality
	ctx := context.Background()
	
	// Simple test script
	code := `
try(
  func() {
    error("test error")
  },
  func(err) {
    print("Error message:", err.message())
    print("Traceback:", err.traceback())
  }
)
`
	
	// Evaluate the code
	result, err := risor.Eval(ctx, code)
	if err != nil {
		log.Fatalf("Runtime error: %v", err)
	}
	
	fmt.Printf("Result: %v\n", result)
	fmt.Println("Test completed successfully!")
}