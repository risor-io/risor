package vm

import (
	"context"

	"github.com/cloudcmds/tamarin/v2/builtins"
	"github.com/cloudcmds/tamarin/v2/compiler"
	"github.com/cloudcmds/tamarin/v2/importer"
	modBytes "github.com/cloudcmds/tamarin/v2/modules/bytes"
	modFmt "github.com/cloudcmds/tamarin/v2/modules/fmt"
	modJson "github.com/cloudcmds/tamarin/v2/modules/json"
	modMath "github.com/cloudcmds/tamarin/v2/modules/math"
	modRand "github.com/cloudcmds/tamarin/v2/modules/rand"
	modStrconv "github.com/cloudcmds/tamarin/v2/modules/strconv"
	modStrings "github.com/cloudcmds/tamarin/v2/modules/strings"
	modTime "github.com/cloudcmds/tamarin/v2/modules/time"
	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/parser"
)

func run(ctx context.Context, code string) (object.Object, error) {

	builtins := builtins.Builtins()
	for k, v := range modFmt.Builtins() {
		builtins[k] = v
	}
	for k, v := range defaultModules() {
		builtins[k] = v
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
