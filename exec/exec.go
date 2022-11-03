package exec

import (
	"context"
	_ "embed"
	"errors"
	"strings"

	"github.com/myzie/tamarin/evaluator"
	"github.com/myzie/tamarin/lexer"
	modJson "github.com/myzie/tamarin/modules/json"
	modMath "github.com/myzie/tamarin/modules/math"
	modSql "github.com/myzie/tamarin/modules/sql"
	modStrings "github.com/myzie/tamarin/modules/strings"
	"github.com/myzie/tamarin/object"
	"github.com/myzie/tamarin/parser"
	"github.com/myzie/tamarin/scope"
)

func Execute(ctx context.Context, input string, importer evaluator.Importer) (object.Object, error) {

	e := evaluator.New(evaluator.Opts{Importer: importer})
	s := scope.New(scope.Opts{Name: "global"})

	// Automatically "import" standard modules
	mathModule, err := modMath.Module(s)
	if err != nil {
		return nil, err
	}
	s.Declare("math", mathModule, false)

	jsonModule, err := modJson.Module(s)
	if err != nil {
		return nil, err
	}
	s.Declare("json", jsonModule, false)

	stringsModule, err := modStrings.Module(s)
	if err != nil {
		return nil, err
	}
	s.Declare("strings", stringsModule, false)

	sqlModule, err := modSql.Module(s)
	if err != nil {
		return nil, err
	}
	s.Declare("sql", sqlModule, false)

	// Parse the user supplied program
	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}

	// Evaluate the program
	result := e.Evaluate(ctx, program, s)
	if result == nil {
		return nil, nil
	}
	if result.Type() == "ERROR" {
		return nil, errors.New(result.Inspect())
	}
	return result, nil
}
