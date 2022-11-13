package exec

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudcmds/tamarin/internal/evaluator"
	"github.com/cloudcmds/tamarin/internal/lexer"
	modJson "github.com/cloudcmds/tamarin/internal/modules/json"
	modMath "github.com/cloudcmds/tamarin/internal/modules/math"
	modRand "github.com/cloudcmds/tamarin/internal/modules/rand"
	modSql "github.com/cloudcmds/tamarin/internal/modules/sql"
	modStrings "github.com/cloudcmds/tamarin/internal/modules/strings"
	modTime "github.com/cloudcmds/tamarin/internal/modules/time"
	modUuid "github.com/cloudcmds/tamarin/internal/modules/uuid"
	"github.com/cloudcmds/tamarin/internal/parser"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

type ModuleFunc func(*scope.Scope) (*object.Module, error)

var moduleFuncs = map[string]ModuleFunc{}

func init() {
	moduleFuncs["math"] = modMath.Module
	moduleFuncs["json"] = modJson.Module
	moduleFuncs["strings"] = modStrings.Module
	moduleFuncs["sql"] = modSql.Module
	moduleFuncs["time"] = modTime.Module
	moduleFuncs["uuid"] = modUuid.Module
	moduleFuncs["rand"] = modRand.Module
}

func Execute(ctx context.Context, input string, importer evaluator.Importer) (object.Object, error) {

	e := evaluator.New(evaluator.Opts{Importer: importer})
	s := scope.New(scope.Opts{Name: "global"})

	// Automatically "import" standard modules
	for name, fn := range moduleFuncs {
		mod, err := fn(s)
		if err != nil {
			return nil, err
		}
		if err := s.Declare(name, mod, false); err != nil {
			return nil, err
		}
	}

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
