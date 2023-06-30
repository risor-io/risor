package vm

import (
	"context"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/parser"
)

func CompileModule(ctx context.Context, name, code string, builtins map[string]object.Object) (*object.Module, error) {
	ast, err := parser.Parse(ctx, code)
	if err != nil {
		return nil, err
	}
	c, err := compiler.New(compiler.WithBuiltins(builtins))
	if err != nil {
		return nil, err
	}
	main, err := c.Compile(ast)
	if err != nil {
		return nil, err
	}
	main.Name = name
	return object.NewModule(name, main), nil
}

type SimpleImporter struct {
	moduleSource map[string]string
	builtins     map[string]object.Object
}

func NewSimpleImporter(moduleSource map[string]string, builtins map[string]object.Object) *SimpleImporter {
	return &SimpleImporter{
		moduleSource: moduleSource,
		builtins:     builtins,
	}
}

func (i *SimpleImporter) Import(ctx context.Context, name string) (*object.Module, error) {
	if code, ok := i.moduleSource[name]; ok {
		return CompileModule(ctx, name, code, i.builtins)
	}
	return nil, nil
}
