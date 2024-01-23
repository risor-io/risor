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

//go:generate go run ./cmd/risor-modgen

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

// WithoutGlobal opts out of a given global builtin or module. If the name can't
// be resolved, this is a no-op. This does operate on nested modules.
func WithoutGlobal(name string) Option {
	return func(cfg *Config) {
		parts := strings.SplitN(name, ".", 2)
		if len(parts) == 1 {
			delete(cfg.Globals, name)
			delete(cfg.DefaultGlobals, name)
			return
		}
		// Resolve the module (which could be nested) and then remove the
		// named attribute from it. This could be in either or both of
		// the Globals and DefaultGlobals structures.
		moduleName, attr := parts[0], parts[1]
		if obj, ok := cfg.Globals[moduleName]; ok {
			if m, ok := obj.(*object.Module); ok {
				removeModuleAttr(m, attr)
			}
		}
		if obj, ok := cfg.DefaultGlobals[moduleName]; ok {
			if m, ok := obj.(*object.Module); ok {
				removeModuleAttr(m, attr)
			}
		}
	}
}

// WithoutGlobals removes multiple global builtins or modules.
func WithoutGlobals(names ...string) Option {
	return func(cfg *Config) {
		for _, name := range names {
			WithoutGlobal(name)(cfg)
		}
	}
}

// WithGlobalOverride replaces the a global or module builtin with the given value
func WithGlobalOverride(name string, value any) Option {
	return func(cfg *Config) {
		parts := strings.Split(name, ".")
		if len(parts) == 1 {
			cfg.Globals[name] = value
			delete(cfg.DefaultGlobals, name)
			return
		}
		valueObj := object.FromGoType(value)
		if valueObj == nil || valueObj.Type() == object.ERROR {
			panic(fmt.Sprintf("invalid value for global override: %v", value))
		}
		moduleName := parts[0]
		nestedModulePath := parts[1 : len(parts)-1]
		attrName := parts[len(parts)-1]
		if obj, ok := cfg.Globals[moduleName]; ok {
			if m, ok := obj.(*object.Module); ok {
				if targetMod, ok := resolveModule(m, nestedModulePath); ok {
					targetMod.Override(attrName, valueObj)
				}
			}
		}
		if obj, ok := cfg.DefaultGlobals[moduleName]; ok {
			if m, ok := obj.(*object.Module); ok {
				if targetMod, ok := resolveModule(m, nestedModulePath); ok {
					targetMod.Override(attrName, valueObj)
				}
			}
		}
	}
}

// WithoutDefaultGlobals opts out of all default global builtins and modules.
func WithoutDefaultGlobals() Option {
	return func(cfg *Config) {
		cfg.DefaultGlobals = map[string]object.Object{}
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
	ast, err := parser.Parse(ctx, source)
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

func resolveModule(m *object.Module, attr []string) (*object.Module, bool) {
	if len(attr) == 0 {
		return m, true
	}
	var result *object.Module
	for _, name := range attr {
		if obj, ok := m.GetAttr(name); ok {
			if modObj, ok := obj.(*object.Module); ok {
				result = modObj
				continue
			}
		}
		return nil, false
	}
	return result, true
}

func removeModuleAttr(m *object.Module, attr string) {
	parts := strings.Split(attr, ".")
	partsLen := len(parts)
	if partsLen == 1 {
		m.Override(attr, nil)
		return
	}
	name := parts[partsLen-1]
	modPath := parts[:partsLen-1]
	if mod, ok := resolveModule(m, modPath); ok {
		mod.Override(name, nil)
	}
}
