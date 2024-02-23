package vm

import (
	"github.com/risor-io/risor/importer"
	"github.com/risor-io/risor/limits"
)

// Option is a configuration function for a Virtual Machine.
type Option func(*VirtualMachine)

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

// WithLimits sets the limits for the Virtual Machine.
func WithLimits(limits limits.Limits) Option {
	return func(vm *VirtualMachine) {
		vm.limits = limits
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

// WithConcurrency opts into allowing the spawning of Goroutines
func WithConcurrency() Option {
	return func(vm *VirtualMachine) {
		vm.concAllowed = true
	}
}
