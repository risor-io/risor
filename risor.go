package risor

import (
	"context"

	"github.com/risor-io/risor/builtins"
	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/importer"
	modBase64 "github.com/risor-io/risor/modules/base64"
	modBytes "github.com/risor-io/risor/modules/bytes"
	modHash "github.com/risor-io/risor/modules/hash"
	modJson "github.com/risor-io/risor/modules/json"
	modMath "github.com/risor-io/risor/modules/math"
	modRand "github.com/risor-io/risor/modules/rand"
	modStrconv "github.com/risor-io/risor/modules/strconv"
	modStrings "github.com/risor-io/risor/modules/strings"
	modTime "github.com/risor-io/risor/modules/time"
	modUuid "github.com/risor-io/risor/modules/uuid"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"
)

type Risor struct {
	compiler *compiler.Compiler
	main     *object.Code
	builtins map[string]object.Object
	importer importer.Importer
	offset   int
}

type Option func(*Risor)

func WithDefaultBuiltins() Option {
	return func(t *Risor) {
		for k, v := range builtins.Builtins() {
			t.builtins[k] = v
		}
	}
}

func WithDefaultModules() Option {
	return func(t *Risor) {
		for k, v := range defaultModules() {
			t.builtins[k] = v
		}
	}
}

func WithBuiltins(builtins map[string]object.Object) Option {
	return func(t *Risor) {
		for k, v := range builtins {
			t.builtins[k] = v
		}
	}
}

func WithCompiler(c *compiler.Compiler) Option {
	return func(t *Risor) {
		t.compiler = c
	}
}

func WithImporter(i importer.Importer) Option {
	return func(t *Risor) {
		t.importer = i
	}
}

func WithCode(c *object.Code) Option {
	return func(t *Risor) {
		t.main = c
	}
}

func WithInstructionOffset(offset int) Option {
	return func(t *Risor) {
		t.offset = offset
	}
}

func Eval(ctx context.Context, source string, options ...Option) (object.Object, error) {

	t := Risor{
		builtins: map[string]object.Object{},
	}
	for _, opt := range options {
		opt(&t)
	}

	// Initialize a compiler if one was not provided via opts.
	if t.compiler == nil {
		var err error
		var compilerOpts []compiler.Option
		if t.builtins != nil {
			compilerOpts = append(compilerOpts, compiler.WithBuiltins(t.builtins))
		}
		if t.main != nil {
			compilerOpts = append(compilerOpts, compiler.WithCode(t.main))
		}
		t.compiler, err = compiler.New(compilerOpts...)
		if err != nil {
			return nil, err
		}
	}

	// Parse the source code to create the AST.
	ast, err := parser.Parse(ctx, source)
	if err != nil {
		return nil, err
	}

	// Compile the AST to bytecode, appending these new instructions after any
	// instructions that were previously compiled.
	main, err := t.compiler.Compile(ast)
	if err != nil {
		return nil, err
	}

	// Eval the bytecode in a new VM then return the top-of-stack (TOS) value.
	var vmOpts []vm.Option
	if t.importer != nil {
		vmOpts = append(vmOpts, vm.WithImporter(t.importer))
	}
	if t.offset != 0 {
		vmOpts = append(vmOpts, vm.WithInstructionOffset(t.offset))
	}
	machine := vm.New(main, vmOpts...)
	if err := machine.Run(ctx); err != nil {
		return nil, err
	}
	if result, exists := machine.TOS(); exists {
		return result, nil
	}
	return object.Nil, nil
}

func defaultModules() map[string]object.Object {
	result := map[string]object.Object{
		"math":    modMath.Module(),
		"json":    modJson.Module(),
		"strings": modStrings.Module(),
		"time":    modTime.Module(),
		"rand":    modRand.Module(),
		"strconv": modStrconv.Module(),
		"uuid":    modUuid.Module(),
		"bytes":   modBytes.Module(),
		"base64":  modBase64.Module(),
	}
	for k, v := range modHash.Builtins() {
		result[k] = v
	}
	return result
}
