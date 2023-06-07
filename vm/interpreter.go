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

type Interpreter struct {
	Compiler *compiler.Compiler
	Main     *object.Code
	Builtins map[string]object.Object
	Importer importer.Importer
}

type InterpreterOpt func(*Interpreter)

func WithDefaultBuiltins() InterpreterOpt {
	return func(i *Interpreter) {
		for _, b := range GlobalBuiltins() {
			i.Builtins[b.Key()] = b
		}
	}
}

func WithDefaultModules() InterpreterOpt {
	return func(i *Interpreter) {
		i.Builtins["math"] = modMath.Module()
		i.Builtins["json"] = modJson.Module()
		i.Builtins["strings"] = modStrings.Module()
		i.Builtins["time"] = modTime.Module()
		i.Builtins["uuid"] = modUuid.Module()
		i.Builtins["rand"] = modRand.Module()
		i.Builtins["strconv"] = modStrconv.Module()
		i.Builtins["pgx"] = modPgx.Module()
	}
}

func WithImporter(im importer.Importer) InterpreterOpt {
	return func(i *Interpreter) {
		i.Importer = im
	}
}

func NewInterpreter(opts ...InterpreterOpt) *Interpreter {
	i := &Interpreter{
		Main:     object.NewCode("main"),
		Builtins: map[string]object.Object{},
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}

func (i *Interpreter) Eval(ctx context.Context, code string) (object.Object, error) {

	// Initialize a compiler as needed
	if i.Compiler == nil {
		i.Compiler = compiler.New(compiler.Options{
			Builtins: i.Builtins,
			Name:     "main",
			Code:     i.Main,
		})
	}

	// Parse the source to create the AST
	ast, err := parser.Parse(code)
	if err != nil {
		return nil, err
	}

	// Compile the AST to bytecode, appending to the end of the main code
	offset := len(i.Compiler.MainInstructions())
	if _, err = i.Compiler.Compile(ast); err != nil {
		return nil, err
	}

	// Evaluate the bytecode in a VM then return the top-of-stack (TOS) value
	machine := New(Options{
		Main:              i.Main,
		InstructionOffset: offset,
		Importer:          i.Importer,
	})
	if err := machine.Run(ctx); err != nil {
		return nil, err
	}
	if result, exists := machine.TOS(); exists {
		return result, nil
	}
	return object.Nil, nil
}
