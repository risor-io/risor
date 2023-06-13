package vm

import (
	"context"

	"github.com/cloudcmds/tamarin/v2/builtins"
	"github.com/cloudcmds/tamarin/v2/compiler"
	"github.com/cloudcmds/tamarin/v2/importer"
	"github.com/cloudcmds/tamarin/v2/modules/all"
	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/parser"
)

func run(ctx context.Context, code string) (object.Object, error) {

	builtins := builtins.Defaults()
	for k, v := range all.Defaults() {
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
