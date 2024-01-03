package risor

import (
	"context"
	"fmt"
	"strings"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/importer"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"
)

// Option describes a function used to configure a Risor evaluation.
type Option func(*Config)

// WithGlobals provides global variables that are made available to Risor
// evaluations. This option is additive, so multiple WithGlobals options
// may be supplied. If the same key is supplied multiple times, the last
// supplied value is used.
func WithGlobals(globals map[string]any) Option {
	return func(cfg *Config) {
		for k, v := range globals {
			cfg.Globals[k] = v
		}
	}
}

// WithGlobal supplies a single named global variable to the Risor evaluation.
func WithGlobal(name string, value any) Option {
	return func(cfg *Config) {
		cfg.Globals[name] = value
	}
}

// WithoutGlobal opts out of a given global builtin or module.
func WithoutGlobal(name string) Option {
	return func(cfg *Config) {
		delete(cfg.Globals, name)
	}
}

// WithoutDefaultGlobals opts out of all default global builtins and modules.
func WithoutDefaultGlobals() Option {
	return func(cfg *Config) {
		cfg.DefaultGlobals = map[string]object.Object{}
	}
}

// WithoutDefaultGlobal opts out of a given default global builtin or module.
func WithoutDefaultGlobal(name string) Option {
	return func(cfg *Config) {
		delete(cfg.DefaultGlobals, name)
	}
}

// WithDenylist opts out of a given set of global or module builtins.
func WithDenylist(names ...string) Option {
	return func(cfg *Config) {
		for _, name := range names {
			delete(cfg.DefaultGlobals, name)
			delete(cfg.Globals, name)
			parts := strings.Split(name, ".")
			if len(parts) == 1 {
				continue
			}
			WithGlobalOverride(name, object.NewBuiltin(parts[1], func(ctx context.Context, args ...object.Object) object.Object {
				return object.Errorf("compile error: undefined variable %q", parts[1])
			}))(cfg)
		}
	}
}

// WithGlobalOverride replaces the a global or module builtin with the given value
func WithGlobalOverride(name string, value any) Option {
	return func(cfg *Config) {
		parts := strings.Split(name, ".")
		if len(parts) == 1 {
			if _, ok := cfg.Globals[name]; ok {
				cfg.Globals[name] = value
			}
			if _, ok := cfg.DefaultGlobals[name]; ok {
				if o, ok := value.(object.Object); ok {
					cfg.DefaultGlobals[name] = o
				}
			}
			return
		}
		value, ok := value.(*object.Builtin)
		if !ok {
			panic("value must be a Builtin object!")
		}
		if b, ok := cfg.Globals[parts[0]]; ok {
			if m, ok := b.(*object.Module); ok {
				if err := m.Override(parts[1], value); err != nil {
					panic(err)
				}
			}
		}
		if b, ok := cfg.DefaultGlobals[parts[0]]; ok {
			if m, ok := b.(*object.Module); ok {
				if err := m.Override(parts[1], value); err != nil {
					panic(err)
				}
			}
		}
	}
}

// WithImporter supplies an Importer that will be used to execute import statements.
func WithImporter(i importer.Importer) Option {
	return func(cfg *Config) {
		cfg.Importer = i
	}
}

// WithLocalImporter enables importing Risor modules from the given directory.
func WithLocalImporter(path string) Option {
	return func(cfg *Config) {
		cfg.LocalImportPath = path
	}
}

// Eval evaluates the given source code and returns the result.
func Eval(ctx context.Context, source string, options ...Option) (object.Object, error) {
	cfg := NewConfig()
	for _, opt := range options {
		opt(cfg)
	}
	// Parse the source code to create the AST
	ast, err := parser.Parse(ctx, source, cfg.ParserOpts()...)
	if err != nil {
		return nil, err
	}
	// Compile the AST to bytecode, appending these new instructions after any
	// instructions that were previously compiled
	main, err := compiler.Compile(ast, cfg.CompilerOpts()...)
	if err != nil {
		return nil, err
	}
	// Eval the bytecode in a VM then return the top-of-stack (TOS) value
	return vm.Run(ctx, main, cfg.VMOpts()...)
}

// EvalCode evaluates the precompiled code and returns the result.
func EvalCode(ctx context.Context, main *compiler.Code, options ...Option) (object.Object, error) {
	cfg := NewConfig()
	for _, opt := range options {
		opt(cfg)
	}
	// Eval the bytecode in a VM then return the top-of-stack (TOS) value
	return vm.Run(ctx, main, cfg.VMOpts()...)
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
	cfg := NewConfig()
	for _, opt := range options {
		opt(cfg)
	}
	vm := vm.New(main, cfg.VMOpts()...)
	if err := vm.Run(ctx); err != nil {
		return nil, err
	}
	obj, err := vm.Get(functionName)
	if err != nil {
		return nil, err
	}
	fn, ok := obj.(*object.Function)
	if !ok {
		return nil, fmt.Errorf("object is not a function (got: %s)", obj.Type())
	}
	return vm.Call(ctx, fn, args)
}
