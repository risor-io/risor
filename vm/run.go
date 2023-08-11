package vm

import (
	"context"
	"reflect"

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
	Inject map[string]interface{}
}

func run(ctx context.Context, code string, opts ...runOpts) (object.Object, error) {

	builtins := builtins.Builtins()
	for k, v := range modFmt.Builtins() {
		builtins[k] = v
	}
	for k, v := range defaultModules() {
		builtins[k] = v
	}
	if len(opts) > 0 {
		for k, v := range opts[0].Inject {
			conv, err := object.NewTypeConverter(reflect.TypeOf(v))
			if err != nil {
				return nil, err
			}
			wrapped, err := conv.From(v)
			if err != nil {
				return nil, err
			}
			builtins[k] = wrapped
		}
	}

	im := importer.NewLocalImporter(importer.LocalImporterOptions{
		SourceDir:  ".",
		Extensions: []string{".tm"},
		Builtins:   builtins,
	})

	// Parse
	ast, err := parser.Parse(ctx, code)
	if err != nil {
		return nil, err
	}

	// Compile
	main, err := compiler.Compile(ast, compiler.WithBuiltins(builtins))
	if err != nil {
		return nil, err
	}

	// Execute
	machine := New(main, WithImporter(im))
	if err := machine.Run(ctx); err != nil {
		return nil, err
	}
	if result, exists := machine.TOS(); exists {
		return result, nil
	}
	return object.Nil, nil
}

func defaultModules() map[string]object.Object {
	return map[string]object.Object{
		"math":    modMath.Module(),
		"json":    modJson.Module(),
		"strings": modStrings.Module(),
		"time":    modTime.Module(),
		"rand":    modRand.Module(),
		"strconv": modStrconv.Module(),
		"bytes":   modBytes.Module(),
	}
}
