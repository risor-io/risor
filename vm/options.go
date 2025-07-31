package vm

import (
	"github.com/risor-io/risor/importer"
	"github.com/risor-io/risor/os"
)

// Option is a configuration function for a Virtual Machine.
type Option func(*VirtualMachine)

// VMMemoryLimits defines configurable memory limits for VM dynamic allocation
type VMMemoryLimits struct {
	MaxStackSize      int
	MaxFrameCount     int
	MaxArgsLimit      int
	InitialStackSize  int
	InitialFrameCount int
}

// DefaultMemoryLimits returns the default memory limits
func DefaultMemoryLimits() VMMemoryLimits {
	return VMMemoryLimits{
		MaxStackSize:      64 * 1024, // 64K slots
		MaxFrameCount:     8 * 1024,  // 8K frames
		MaxArgsLimit:      1024,      // 1K args max
		InitialStackSize:  256,       // Start with 256 slots
		InitialFrameCount: 64,        // Start with 64 frames
	}
}

// WithMemoryLimits configures dynamic memory allocation limits
func WithMemoryLimits(limits VMMemoryLimits) Option {
	return func(vm *VirtualMachine) {
		vm.memoryLimits = limits
		vm.initializeDynamicMemory()
	}
}

// WithDefaultMemoryLimits configures VM to use default dynamic memory limits
func WithDefaultMemoryLimits() Option {
	return WithMemoryLimits(DefaultMemoryLimits())
}

// WithInstructionOffset sets the initial instruction offset.
func WithInstructionOffset(offset int) Option {
	return func(vm *VirtualMachine) {
		vm.ip = offset
	}
}

// WithImporter is used to supply an Importer to the Virtual Machine.
func WithImporter(importer importer.Importer) Option {
	return func(vm *VirtualMachine) {
		vm.importer = importer
	}
}

// WithGlobals provides global variables with the given names.
func WithGlobals(globals map[string]any) Option {
	return func(vm *VirtualMachine) {
		for name, value := range globals {
			vm.inputGlobals[name] = value
		}
	}
}

// WithConcurrency opts into allowing the spawning of goroutines.
func WithConcurrency() Option {
	return func(vm *VirtualMachine) {
		vm.concAllowed = true
	}
}

// WithOS sets custom OS implementation in the context. This context is present
// in the invocation of Risor builtins, this OS will be used for all related
// functionality.
func WithOS(os os.OS) Option {
	return func(vm *VirtualMachine) {
		vm.os = os
	}
}
