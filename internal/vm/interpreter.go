package vm

import (
	"context"

	"github.com/cloudcmds/tamarin/importer"
	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
)

type Interpreter struct {
	c    *compiler.Compiler
	main *object.Code
}

func NewInterpreter(builtins []*object.Builtin) *Interpreter {

	bmap := map[string]object.Object{}
	for _, b := range GlobalBuiltins() {
		bmap[b.Key()] = b
	}

	// bmap["math"] = modMath.Module()
	// bmap["json"] = modJson.Module()
	// bmap["strings"] = modStrings.Module()
	// bmap["time"] = modTime.Module()
	// bmap["uuid"] = modUuid.Module()
	// bmap["rand"] = modRand.Module()
	// bmap["strconv"] = modStrconv.Module()
	// bmap["pgx"] = modPgx.Module()

	s := object.NewCode("main")

	c := compiler.New(compiler.Options{
		Builtins: bmap,
		Name:     "main",
		Code:     s,
	})

	return &Interpreter{c: c, main: s}
}

func (i *Interpreter) Eval(ctx context.Context, code string) (object.Object, error) {
	ast, err := parser.Parse(code)
	if err != nil {
		return nil, err
	}
	pos := len(i.c.Instructions())

	if _, err = i.c.Compile(ast); err != nil {
		return nil, err
	}

	v := New(Options{
		Main:              i.main,
		InstructionOffset: pos,
		Importer: importer.NewLocalImporter(importer.LocalImporterOptions{
			SourceDir: ".",
		}),
	})
	if err := v.Run(ctx); err != nil {
		return nil, err
	}

	result, exists := v.TOS()
	if exists {
		return result, nil
	}
	return object.Nil, nil
}
