// Package exec provides an Execute function that is used to
// run arbitrary Tamarin source code and return the result.
package exec

import (
	"context"
	"fmt"

	"github.com/cloudcmds/tamarin/evaluator"
	modJson "github.com/cloudcmds/tamarin/internal/modules/json"
	modMath "github.com/cloudcmds/tamarin/internal/modules/math"
	modRand "github.com/cloudcmds/tamarin/internal/modules/rand"
	modSql "github.com/cloudcmds/tamarin/internal/modules/sql"
	modStrconv "github.com/cloudcmds/tamarin/internal/modules/strconv"
	modStrings "github.com/cloudcmds/tamarin/internal/modules/strings"
	modTime "github.com/cloudcmds/tamarin/internal/modules/time"
	modUuid "github.com/cloudcmds/tamarin/internal/modules/uuid"
	"github.com/cloudcmds/tamarin/internal/parser"
	"github.com/cloudcmds/tamarin/internal/scope"
	"github.com/cloudcmds/tamarin/object"
)

// ModuleFunc is the signature of a function that returns a module
type ModuleFunc func(*scope.Scope) (*object.Module, error)

// Will contain module functions for default modules
var moduleFuncs = map[string]ModuleFunc{}

// Modules included here must not include any I/O operations or
// other operations that may be questionable in a secure environment.
// This is because these modules are imported automatically and
// some callers may want to have a limited core set of modules.
func init() {
	moduleFuncs["math"] = modMath.Module
	moduleFuncs["json"] = modJson.Module
	moduleFuncs["strings"] = modStrings.Module
	moduleFuncs["time"] = modTime.Module
	moduleFuncs["uuid"] = modUuid.Module
	moduleFuncs["rand"] = modRand.Module
	moduleFuncs["strconv"] = modStrconv.Module
	moduleFuncs["sql"] = modSql.Module
}

// Opts is used configure the execution of a Tamarin program.
type Opts struct {
	// Input is the main source code to execute.
	Input string

	// Importer may optionally be supplied as an interface
	// used to import modules. If not provided, any attempt
	// to import will fail, halting execution with an error.
	Importer evaluator.Importer

	// Scope may optionally be supplied as the top-level scope
	// used during execution. If not provided, an empty scope
	// will be created automatically.
	Scope *scope.Scope

	// If set to true, the default modules will not be imported
	// automatically.
	DisableAutoImport bool

	// If set to true, the default builtins will not be registered.
	DisableDefaultBuiltins bool

	// Supplies extra and/or override builtins for evaluation.
	Builtins []*object.Builtin
}

// Execute the given source code as input and return the result.
// If the execution is successful, a Tamarin object is returned
// as the final result. The context may be used to cancel the
// evaluation based on a timeout or otherwise.
//
// The opts should contain the required input as well as other
// optional parameters.
//
// Any panic is handled internally and propagated as an error.
//
// The result value is the final of the final statement or
// expression in the main source code, which may be object.NULL
// if the expression doesn't evaluate to a value.
func Execute(ctx context.Context, opts Opts) (result object.Object, err error) {

	// Translate any panic into an error so the caller has a good guarantee
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	// Create the top-level scope if one was not provided
	s := opts.Scope
	if s == nil {
		s = scope.New(scope.Opts{Name: "global"})
	}

	// Conditionally auto import standard modules
	if !opts.DisableAutoImport {
		for name, fn := range moduleFuncs {
			mod, err := fn(s)
			if err != nil {
				return nil, fmt.Errorf("init error: failed to create module %s: %w", name, err)
			}
			if err := s.Declare(name, mod, false); err != nil {
				return nil, fmt.Errorf("init error: failed to attach module %s: %w", name, err)
			}
		}
	}

	// Parse the program
	program, err := parser.Parse(opts.Input)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	// Evaluate the program
	result = evaluator.New(evaluator.Opts{
		Importer:               opts.Importer,
		DisableDefaultBuiltins: opts.DisableDefaultBuiltins,
		Builtins:               opts.Builtins,
	}).Evaluate(ctx, program, s)

	// Let's guarantee that if there's no error we return a
	// Tamarin object, so defaulting to object.NULL may make sense
	if result == nil {
		return object.NULL, nil
	}

	// If evaluation failed, we will have a Tamarin error object
	// and we should transform that into a Go error
	if errObj, ok := result.(*object.Error); ok {
		return nil, fmt.Errorf("eval error: %s", errObj.Message)
	}

	// At this point we know evaluation succeeded and we can
	// just return the final Tamarin object as-is
	return result, nil
}
