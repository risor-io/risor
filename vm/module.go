package vm

import (
	"context"

	"github.com/cloudcmds/tamarin/v2/compiler"
	"github.com/cloudcmds/tamarin/v2/object"
	"github.com/cloudcmds/tamarin/v2/parser"
)

func CompileModule(name, code string, builtins map[string]object.Object) (*object.Module, error) {
	ast, err := parser.Parse(code)
	if err != nil {
		return nil, err
	}
	c := compiler.New(compiler.Options{
		Name:     name,
		Builtins: builtins,
	})
	main, err := c.Compile(ast)
	if err != nil {
		return nil, err
	}
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
		return CompileModule(name, code, i.builtins)
	}
	return nil, nil
}
