package vm

import (
	"context"

	"github.com/cloudcmds/tamarin/internal/compiler"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/parser"
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

	c := compiler.New(compiler.Options{
		Builtins: builtins,
		Name:     "main",
	})
	mainScope, err := c.Compile(ast)
	if err != nil {
		return nil, err
	}

	vm := New(mainScope)
	if err := vm.Run(ctx); err != nil {
		return nil, err
	}
	result, exists := vm.TOS()
	if exists {
		return result, nil
	}
	return object.Nil, nil
}
