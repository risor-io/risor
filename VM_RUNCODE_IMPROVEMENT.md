# VM RunCode Improvement

## Overview

This improvement adds the ability to run multiple different compiled code objects on a single Risor VM instance. Previously, each VM was tied to a single code object that was set during VM creation. Now, VMs can execute arbitrary compiled code objects sequentially while preserving global variables and state.

## Key Changes

### New Methods

1. **`RunCode(ctx context.Context, codeToRun *compiler.Code) error`**
   - Runs arbitrary compiled code on an existing VM
   - Resets VM execution state but preserves global variables
   - Can be called multiple times on the same VM

2. **`RunCodeOnVM(ctx context.Context, vm *VirtualMachine, code *compiler.Code) (object.Object, error)`**
   - Convenience wrapper function that runs code and returns the result
   - Similar to the existing `Run` function but works on existing VMs

### Internal Changes

- Added `resetForNewCode()` method to properly reset VM state between code executions
- Modified code loading to handle different code objects
- Ensured global variables are preserved across different code executions

## Use Cases

This improvement is particularly useful for:

1. **Game Development**: Running scripts for different game objects on a single VM
2. **REPL/Interactive Systems**: Executing multiple user inputs on the same VM
3. **Script Engines**: Running different scripts while maintaining shared state
4. **Testing**: Running multiple test cases on the same VM instance

## Example Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"
)

func main() {
	ctx := context.Background()
	
	// Create a VM with initial code
	source1 := `x := 10; y := 20; x + y`
	ast1, _ := parser.Parse(ctx, source1)
	code1, _ := compiler.Compile(ast1)
	machine := vm.New(code1)
	
	// Run the initial code
	machine.Run(ctx)
	result1, _ := machine.TOS()
	fmt.Printf("Result 1: %s\n", result1.Inspect()) // Output: 30
	
	// Run different code on the same VM
	source2 := `a := 5; b := 15; a * b`
	ast2, _ := parser.Parse(ctx, source2)
	code2, _ := compiler.Compile(ast2)
	machine.RunCode(ctx, code2)
	result2, _ := machine.TOS()
	fmt.Printf("Result 2: %s\n", result2.Inspect()) // Output: 75
	
	// Use convenience function
	source3 := `name := "Risor"; "Hello, " + name + "!"`
	ast3, _ := parser.Parse(ctx, source3)
	code3, _ := compiler.Compile(ast3)
	result3, _ := vm.RunCodeOnVM(ctx, machine, code3)
	fmt.Printf("Result 3: %s\n", result3.Inspect()) // Output: "Hello, Risor!"
}
```

## Benefits

1. **Memory Efficiency**: Reuse VM instances instead of creating new ones
2. **Performance**: Avoid VM initialization overhead for subsequent code executions
3. **State Persistence**: Global variables and modules remain available across executions
4. **Flexibility**: Support for dynamic code execution scenarios

## Backwards Compatibility

This change is fully backwards compatible. Existing code using `vm.Run()` will continue to work unchanged. The new functionality is additive and doesn't modify existing behavior.

## Testing

Comprehensive tests have been added to verify:
- Basic functionality of running multiple code objects
- Global variable preservation and isolation
- Function definitions and calls across different code objects
- Error handling and VM state management
- Integration with existing VM features

## Thread Safety

Like the existing `Run()` method, `RunCode()` is not thread-safe. The VM enforces single-threaded execution and will return an error if called while another execution is in progress.

## Implementation Details

The implementation properly handles:
- VM state reset between code executions (stack, frames, instruction pointer)
- Code object loading and caching
- Global variable preservation
- Error handling and cleanup
- Memory management for temporary objects

This improvement enables the game development use case mentioned in the GitHub issue while maintaining the VM's existing performance and reliability characteristics.