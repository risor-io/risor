package risor

import (
	"context"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/importer"
	"github.com/risor-io/risor/internal/cfg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"
)

// Option describes a function used to configure a Risor evaluation.
type Option func(*cfg.RisorConfig)

// WithGlobals provides global variables that are made available to Risor
// evaluations. This option is additive, so multiple WithGlobals options
// may be supplied. If the same key is supplied multiple times, the last
// supplied value is used.
func WithGlobals(globals map[string]any) Option {
	return func(r *cfg.RisorConfig) {
		for k, v := range globals {
			r.Globals[k] = v
		}
	}
}

// WithGlobal supplies a single named global variable to the Risor evaluation.
func WithGlobal(name string, value any) Option {
	return func(r *cfg.RisorConfig) {
		r.Globals[name] = value
	}
}

// WithoutDefaultGlobals opts out of all default global builtins and modules.
func WithoutDefaultGlobals() Option {
	return func(r *cfg.RisorConfig) {
		r.DefaultGlobals = map[string]object.Object{}
	}
}

// WithImporter supplies an Importer that will be used to execute import statements.
func WithImporter(i importer.Importer) Option {
	return func(r *cfg.RisorConfig) {
		r.Importer = i
	}
}

// WithLocalImporter enables importing Risor modules from the given directory.
func WithLocalImporter(path string) Option {
	return func(r *cfg.RisorConfig) {
		r.LocalImportPath = path
	}
}

// Eval evaluates the given source code and returns the result.
func Eval(ctx context.Context, source string, options ...Option) (object.Object, error) {
	r := cfg.NewRisorConfig()
	for _, opt := range options {
		opt(r)
	}
	// Parse the source code to create the AST
	ast, err := parser.Parse(ctx, source)
	if err != nil {
		return nil, err
	}
	// Compile the AST to bytecode, appending these new instructions after any
	// instructions that were previously compiled
	main, err := compiler.Compile(ast, r.CompilerOpts()...)
	if err != nil {
		return nil, err
	}
	// Eval the bytecode in a VM then return the top-of-stack (TOS) value
	return vm.Run(ctx, main, r.VMOpts()...)
}

// EvalCode evaluates the precompiled code and returns the result.
func EvalCode(ctx context.Context, main *compiler.Code, options ...Option) (object.Object, error) {
	r := cfg.NewRisorConfig()
	for _, opt := range options {
		opt(r)
	}
	// Eval the bytecode in a VM then return the top-of-stack (TOS) value
	return vm.Run(ctx, main, r.VMOpts()...)
}

// Call evaluates the precompiled code and then calls the named function.
// The supplied arguments are passed in the function call. The result of
// the function call is returned.
func Call(
	ctx context.Context,
	main *compiler.Code,
	functionName string,
	args []object.Object,
	options ...Option,
) (object.Object, error) {
	r := cfg.NewRisorConfig()
	for _, opt := range options {
		opt(r)
	}
	vm := vm.New(main, r.VMOpts()...)
	if err := vm.Run(ctx); err != nil {
		return nil, err
	}
	return vm.Call(ctx, functionName, args)
}
