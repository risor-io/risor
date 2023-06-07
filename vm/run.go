package vm

import (
	"context"

	"github.com/cloudcmds/tamarin/v2/importer"
	"github.com/cloudcmds/tamarin/v2/object"
)

func RunWithDefaults(ctx context.Context, code string) (object.Object, error) {

	builtins := GlobalBuiltins()
	builtinsMap := map[string]object.Object{}
	for _, b := range builtins {
		builtinsMap[b.Key()] = b
	}

	importerInst := importer.NewLocalImporter(importer.LocalImporterOptions{
		SourceDir:  ".",
		Extensions: []string{".tm"},
		Builtins:   builtinsMap,
	})

	interp := NewInterpreter(
		WithDefaultBuiltins(),
		WithDefaultModules(),
		WithImporter(importerInst),
	)
	return interp.Eval(ctx, code)
}

func run(ctx context.Context, code string) (object.Object, error) {
	return RunWithDefaults(ctx, code)
}
