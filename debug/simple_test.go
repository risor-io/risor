package main

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"
)

func main() {
	// Test basic VM creation and stack operations
	machine, err := vm.NewEmpty()
	if err != nil {
		fmt.Printf("Failed to create VM: %v\n", err)
		return
	}

	fmt.Printf("VM created successfully\n")
	fmt.Printf("useDynamic: %v\n", getPrivateField(machine, "useDynamic"))

	// Test basic push/pop
	machine.Push(object.NewInt(42))
	if result, exists := machine.TOS(); exists {
		fmt.Printf("TOS: %v\n", result)
	} else {
		fmt.Printf("No TOS available\n")
	}
}

// Helper to access private fields for debugging
func getPrivateField(vm *vm.VirtualMachine, fieldName string) interface{} {
	// This won't work in real code, just for illustration
	return "unknown"
}