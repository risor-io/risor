package vm

import (
	"context"

	"github.com/cloudcmds/tamarin/v2/compiler"
	"github.com/cloudcmds/tamarin/v2/importer"
	modJson "github.com/cloudcmds/tamarin/v2/modules/json"
	modMath "github.com/cloudcmds/tamarin/v2/modules/math"
	modPgx "github.com/cloudcmds/tamarin/v2/modules/pgx"
	modRand "github.com/cloudcmds/tamarin/v2/modules/rand"
	modStrconv "github.com/cloudcmds/tamarin/v2/modules/strconv"
	modStrings "github.com/cloudcmds/tamarin/v2/modules/strings"
	modTime "github.com/cloudcmds/tamarin/v2/modules/time"
	modUuid "github.com/cloudcmds/tamarin/v2/modules/uuid"
	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/parser"
)

func Run(ctx context.Context, code string) (object.Object, error) {

	ast, err := parser.Parse(code)
	if err != nil {
		return nil, err
	}

	builtins := map[string]object.Object{}
	for _, b := range GlobalBuiltins() {
		builtins[b.Key()] = b
	}

	builtins["math"] = modMath.Module()
	builtins["json"] = modJson.Module()
	builtins["strings"] = modStrings.Module()
	builtins["time"] = modTime.Module()
	builtins["uuid"] = modUuid.Module()
	builtins["rand"] = modRand.Module()
	builtins["strconv"] = modStrconv.Module()
	builtins["pgx"] = modPgx.Module()

	c := compiler.New(compiler.Options{
		Builtins: builtins,
		Name:     "main",
	})
	mainScope, err := c.Compile(ast)
	if err != nil {
		return nil, err
	}

	vm := New(Options{
		Main: mainScope,
		Importer: importer.NewLocalImporter(importer.LocalImporterOptions{
			SourceDir:  ".",
			Extensions: []string{".tm"},
			Builtins:   builtins,
		}),
	})
	if err := vm.Run(ctx); err != nil {
		return nil, err
	}
	result, exists := vm.TOS()
	if exists {
		return result, nil
	}
	return object.Nil, nil
}
