package exec

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudcmds/tamarin/internal/evaluator"
	"github.com/cloudcmds/tamarin/internal/lexer"
	modJson "github.com/cloudcmds/tamarin/internal/modules/json"
	modMath "github.com/cloudcmds/tamarin/internal/modules/math"
	modSql "github.com/cloudcmds/tamarin/internal/modules/sql"
	modStrings "github.com/cloudcmds/tamarin/internal/modules/strings"
	modTime "github.com/cloudcmds/tamarin/internal/modules/time"
	"github.com/cloudcmds/tamarin/internal/parser"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
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

	timeModule, err := modTime.Module(s)
	if err != nil {
		return nil, err
	}
	s.Declare("time", timeModule, false)

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
