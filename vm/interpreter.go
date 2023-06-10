package vm

import (
	"context"

	"github.com/cloudcmds/tamarin/v2/builtins"
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

type Interpreter struct {
	compiler *compiler.Compiler
	main     *object.Code
	builtins map[string]object.Object
	importer importer.Importer
}

type InterpreterOpt func(*Interpreter)

func WithDefaultBuiltins() InterpreterOpt {
	return func(i *Interpreter) {
		for _, b := range builtins.GlobalBuiltins() {
			i.builtins[b.Key()] = b
		}
	}
}

func WithDefaultModules() InterpreterOpt {
	return func(i *Interpreter) {
		i.builtins["math"] = modMath.Module()
		i.builtins["json"] = modJson.Module()
		i.builtins["strings"] = modStrings.Module()
		i.builtins["time"] = modTime.Module()
		i.builtins["uuid"] = modUuid.Module()
		i.builtins["rand"] = modRand.Module()
		i.builtins["strconv"] = modStrconv.Module()
		i.builtins["pgx"] = modPgx.Module()
	}
}

func WithImporter(im importer.Importer) InterpreterOpt {
	return func(i *Interpreter) {
		i.importer = im
	}
}

func NewInterpreter(opts ...InterpreterOpt) *Interpreter {
	i := &Interpreter{
		main:     object.NewCode("main"),
		builtins: map[string]object.Object{},
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *Interpreter) Eval(ctx context.Context, code string) (object.Object, error) {

	// Initialize a compiler as needed
	if i.compiler == nil {
		var err error
		i.compiler, err = compiler.New(
			compiler.WithBuiltins(i.builtins),
			compiler.WithCode(i.main),
		)
		if err != nil {
			return nil, err
		}
	}

	// Parse the source to create the AST
	ast, err := parser.Parse(ctx, code)
	if err != nil {
		return nil, err
	}

	// Compile the AST to bytecode, appending to the end of the main code
	offset := len(i.compiler.MainInstructions())
	if _, err = i.compiler.Compile(ast); err != nil {
		return nil, err
	}

	// Evaluate the bytecode in a VM then return the top-of-stack (TOS) value
	machine := New(Options{
		Main:              i.main,
		InstructionOffset: offset,
		Importer:          i.importer,
	})
	if err := machine.Run(ctx); err != nil {
		return nil, err
	}
	if result, exists := machine.TOS(); exists {
		return result, nil
	}
	return object.Nil, nil
}
