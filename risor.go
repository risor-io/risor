package risor

import (
	"context"

	"github.com/risor-io/risor/builtins"
	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/importer"
	"github.com/risor-io/risor/internal/cfg"
	modAws "github.com/risor-io/risor/modules/aws"
	modBase64 "github.com/risor-io/risor/modules/base64"
	modBytes "github.com/risor-io/risor/modules/bytes"
	modFetch "github.com/risor-io/risor/modules/fetch"
	modFmt "github.com/risor-io/risor/modules/fmt"
	modGoogle "github.com/risor-io/risor/modules/google"
	modHash "github.com/risor-io/risor/modules/hash"
	modImage "github.com/risor-io/risor/modules/image"
	modJson "github.com/risor-io/risor/modules/json"
	modMath "github.com/risor-io/risor/modules/math"
	modOs "github.com/risor-io/risor/modules/os"
	modPgx "github.com/risor-io/risor/modules/pgx"
	modRand "github.com/risor-io/risor/modules/rand"
	modStrconv "github.com/risor-io/risor/modules/strconv"
	modStrings "github.com/risor-io/risor/modules/strings"
	modTime "github.com/risor-io/risor/modules/time"
	modUuid "github.com/risor-io/risor/modules/uuid"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
	"github.com/risor-io/risor/vm"
)

type Option func(*cfg.RisorConfig)

func WithDefaultBuiltins() Option {
	return func(r *cfg.RisorConfig) {
		for k, v := range builtins.Builtins() {
			r.Builtins[k] = v
		}
		for k, v := range modFetch.Builtins() {
			r.Builtins[k] = v
		}
		for k, v := range modFmt.Builtins() {
			r.Builtins[k] = v
		}
		for k, v := range modHash.Builtins() {
			r.Builtins[k] = v
		}
		for k, v := range modOs.Builtins() {
			r.Builtins[k] = v
		}
	}
}

func WithDefaultModules() Option {
	return func(r *cfg.RisorConfig) {
		for k, v := range defaultModules() {
			r.Builtins[k] = v
		}
	}
}

func WithBuiltins(builtins map[string]object.Object) Option {
	return func(r *cfg.RisorConfig) {
		for k, v := range builtins {
			r.Builtins[k] = v
		}
	}
}

func WithCompiler(c *compiler.Compiler) Option {
	return func(r *cfg.RisorConfig) {
		r.Compiler = c
	}
}

func WithImporter(i importer.Importer) Option {
	return func(r *cfg.RisorConfig) {
		r.Importer = i
	}
}

func WithLocalImporter(path string) Option {
	return func(r *cfg.RisorConfig) {
		r.LocalImportPath = path
	}
}

func WithCode(c *object.Code) Option {
	return func(r *cfg.RisorConfig) {
		r.Main = c
	}
}

func WithInstructionOffset(offset int) Option {
	return func(r *cfg.RisorConfig) {
		r.Offset = offset
	}
}

func Eval(ctx context.Context, source string, options ...Option) (object.Object, error) {

	r := &cfg.RisorConfig{
		Builtins: map[string]object.Object{},
	}
	for _, opt := range options {
		opt(r)
	}

	// Set up a local module importer if LocalImportPath is set.
	if r.Importer == nil && r.LocalImportPath != "" {
		r.Importer = importer.NewLocalImporter(importer.LocalImporterOptions{
			Builtins:   r.Builtins,
			SourceDir:  r.LocalImportPath,
			Extensions: []string{".risor", ".rsr"},
		})
	}

	// Initialize a compiler if one was not provided via opts.
	if r.Compiler == nil {
		var err error
		var compilerOpts []compiler.Option
		if r.Builtins != nil {
			compilerOpts = append(compilerOpts, compiler.WithBuiltins(r.Builtins))
		}
		if r.Main != nil {
			compilerOpts = append(compilerOpts, compiler.WithCode(r.Main))
		}
		r.Compiler, err = compiler.New(compilerOpts...)
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
	main, err := r.Compiler.Compile(ast)
	if err != nil {
		return nil, err
	}

	// Eval the bytecode in a new VM then return the top-of-stack (TOS) value.
	var vmOpts []vm.Option
	if r.Importer != nil {
		vmOpts = append(vmOpts, vm.WithImporter(r.Importer))
	}
	if r.Offset != 0 {
		vmOpts = append(vmOpts, vm.WithInstructionOffset(r.Offset))
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
		"base64":  modBase64.Module(),
		"bytes":   modBytes.Module(),
		"fmt":     modFmt.Module(),
		"image":   modImage.Module(),
		"json":    modJson.Module(),
		"math":    modMath.Module(),
		"os":      modOs.Module(),
		"pgx":     modPgx.Module(),
		"rand":    modRand.Module(),
		"strconv": modStrconv.Module(),
		"strings": modStrings.Module(),
		"time":    modTime.Module(),
		"uuid":    modUuid.Module(),
	}
	if awsMod := modAws.Module(); awsMod != nil {
		result["aws"] = awsMod
	}
	if googleMod := modGoogle.Module(); googleMod != nil {
		result["google"] = googleMod
	}
	return result
}
