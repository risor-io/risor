package compiler

import (
	"bytes"
	"fmt"
	"strings"
)

type Function struct {
	id         string
	name       string
	parameters []string
	defaults   []any
	code       *Code
	lineNumber int // Line number where function is defined
	columnNumber int // Column number where function is defined
}

func (f *Function) ID() string {
	return f.id
}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Code() *Code {
	return f.code
}

func (f *Function) ParametersCount() int {
	return len(f.parameters)
}

func (f *Function) Parameter(index int) string {
	return f.parameters[index]
}

func (f *Function) DefaultsCount() int {
	return len(f.defaults)
}

func (f *Function) Default(index int) any {
	return f.defaults[index]
}

func (f *Function) RequiredArgsCount() int {
	return len(f.parameters) - len(f.defaults)
}

func (f *Function) LineNumber() int {
	return f.lineNumber
}

func (f *Function) ColumnNumber() int {
	return f.columnNumber
}

func (f *Function) LocalsCount() int {
	return int(f.code.LocalsCount())
}

func (f *Function) String() string {
	var out bytes.Buffer
	parameters := make([]string, 0)
	for i, name := range f.parameters {
		if def := f.defaults[i]; def != nil {
			name += "=" + fmt.Sprintf("%v", def)
		}
		parameters = append(parameters, name)
	}
	out.WriteString("func")
	if f.name != "" {
		out.WriteString(" " + f.name)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") {")
	lines := strings.Split(f.Code().Source(), "\n")
	if len(lines) == 1 {
		out.WriteString(" " + lines[0] + " }")
	} else if len(lines) == 0 {
		out.WriteString(" }")
	} else {
		for _, line := range lines {
			out.WriteString("\n    " + line)
		}
		out.WriteString("\n}")
	}
	return out.String()
}

type FunctionOpts struct {
	ID           string
	Name         string
	Parameters   []string
	Defaults     []any
	Code         *Code
	LineNumber   int
	ColumnNumber int
}

func NewFunction(opts FunctionOpts) *Function {
	return &Function{
		id:           opts.ID,
		name:         opts.Name,
		parameters:   opts.Parameters,
		defaults:     opts.Defaults,
		code:         opts.Code,
		lineNumber:   opts.LineNumber,
		columnNumber: opts.ColumnNumber,
	}
}
