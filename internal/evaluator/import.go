package evaluator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/myzie/tamarin/internal/ast"
	"github.com/myzie/tamarin/internal/lexer"
	"github.com/myzie/tamarin/internal/object"
	"github.com/myzie/tamarin/internal/parser"
	"github.com/myzie/tamarin/internal/scope"
)

type Importer interface {
	Import(ctx context.Context, e *Evaluator, name string) (*object.Module, error)
}

type SimpleImporter struct{}

func (si *SimpleImporter) Import(
	ctx context.Context,
	e *Evaluator,
	name string,
) (*object.Module, error) {
	contents, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	p := parser.New(lexer.New(string(contents)))
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}

	// Module scope
	s := scope.New(scope.Opts{Name: fmt.Sprintf("module:%s", name)})

	result := e.Evaluate(ctx, program, s)
	if result != nil && result.Type() == "ERROR" {
		return nil, errors.New(result.Inspect())
	}

	return &object.Module{Name: name, Scope: s}, nil
}

func (e *Evaluator) evalImportStatement(
	ctx context.Context,
	node *ast.ImportStatement,
	s *scope.Scope,
) object.Object {
	if e.importer == nil {
		return newError("import error: importing is disabled")
	}
	moduleName := node.Name.String()
	name := fmt.Sprintf("%s.mon", moduleName)
	module, err := e.importer.Import(ctx, e, name)
	if err != nil {
		return newError(err.Error())
	}
	if err := s.Declare(moduleName, module, true); err != nil {
		return newError(fmt.Sprintf("import error: %s", err.Error()))
	}
	return module
}
