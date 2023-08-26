package vm

import (
	"context"

	"github.com/risor-io/risor/builtins"
	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/importer"
	modBytes "github.com/risor-io/risor/modules/bytes"
	modFmt "github.com/risor-io/risor/modules/fmt"
	modJson "github.com/risor-io/risor/modules/json"
	modMath "github.com/risor-io/risor/modules/math"
	modRand "github.com/risor-io/risor/modules/rand"
	modStrconv "github.com/risor-io/risor/modules/strconv"
	modStrings "github.com/risor-io/risor/modules/strings"
	modTime "github.com/risor-io/risor/modules/time"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
)

type runOpts struct {
	Globals map[string]interface{}
}

func run(ctx context.Context, source string, opts ...runOpts) (object.Object, error) {
	vm, err := newVM(ctx, source, opts...)
	if err != nil {
		return nil, err
	}
	if err := vm.Run(ctx); err != nil {
		return nil, err
	}
	if result, exists := vm.TOS(); exists {
		return result, nil
	}
	return object.Nil, nil
}

func newVM(ctx context.Context, source string, opts ...runOpts) (*VirtualMachine, error) {
	ast, err := parser.Parse(ctx, source)
	if err != nil {
		return nil, err
	}
	globals := basicBuiltins()
	if len(opts) > 0 {
		for k, v := range opts[0].Globals {
			globals[k] = v
		}
	}
	var globalNames []string
	for k := range globals {
		globalNames = append(globalNames, k)
	}
	main, err := compiler.Compile(ast, compiler.WithGlobalNames(globalNames))
	if err != nil {
		return nil, err
	}
	im := importer.NewLocalImporter(importer.LocalImporterOptions{
		SourceDir:   "./fixtures",
		Extensions:  []string{".risor", ".rsr"},
		GlobalNames: globalNames,
	})
	return New(main, WithImporter(im), WithGlobals(globals)), nil
}

func basicBuiltins() map[string]any {
	globals := map[string]any{
		"math":    modMath.Module(),
		"json":    modJson.Module(),
		"strings": modStrings.Module(),
		"time":    modTime.Module(),
		"rand":    modRand.Module(),
		"strconv": modStrconv.Module(),
		"bytes":   modBytes.Module(),
	}
	for k, v := range builtins.Builtins() {
		globals[k] = v
	}
	for k, v := range modFmt.Builtins() {
		globals[k] = v
	}
	return globals
}
