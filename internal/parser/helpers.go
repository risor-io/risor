package parser

import (
	"errors"
	"strings"

	"github.com/myzie/tamarin/internal/ast"
	"github.com/myzie/tamarin/internal/lexer"
)

// ParseProgram is a shortcut for getting the AST corresponding for some source
func ParseProgram(input string) (*ast.Program, error) {
	parser := New(lexer.New(input))
	program := parser.ParseProgram()
	if errs := parser.Errors(); len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}
	return program, nil
}
