package main

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
)

type ExportedFunc struct {
	ExportedName string
	FuncName     string
	FuncGenName  string
	Return       *FuncReturn
	ReturnsError bool
	Params       []FuncParam
	NeedsContext bool
}

type FuncReturn struct {
	Type     string
	NewFunc  string
	CastFunc string
}

type FuncParam struct {
	Name         string
	Type         string
	ReadFunc     string
	CastFunc     string
	CastMaxValue string
	CastMinValue string
}

func (m *Module) parseFuncDecl(decl *ast.FuncDecl) error {
	exportedName, err := m.parseFuncDeclName(decl)
	if err != nil {
		return err
	}
	if exportedName == "" {
		// empty name = no risor:export comment found
		return nil
	}

	if decl.Recv != nil {
		// TODO: Add support for methods
		return fmt.Errorf("methods are not supported")
	}

	if decl.Type.TypeParams != nil {
		return fmt.Errorf("generic functions are not supported")
	}

	exported := ExportedFunc{
		ExportedName: exportedName,
		FuncName:     decl.Name.Name,
		FuncGenName:  firstRuneUpper(decl.Name.Name),
	}

	if exported.FuncName == exported.FuncGenName {
		wantName := firstRuneLower(exported.FuncName)
		return fmt.Errorf("collision of source function name and generated function name, should rename it to %q", wantName)
	}

	if strings.ToLower(exported.ExportedName) != exported.ExportedName {
		return fmt.Errorf("exported name should be snake_case, but got %q", exported.ExportedName)
	}

	for _, param := range decl.Type.Params.List {
		for _, name := range namesListOrUnnamed(param.Names) {
			p, err := m.parseFuncDeclParam(name, param)
			if err != nil {
				return fmt.Errorf("param %q: %w", name, err)
			}

			// Special treatment for context.Context parameter
			if p.Type == "context.Context" {
				if len(exported.Params) > 0 {
					return fmt.Errorf("param %q: context.Context must be the first parameter", name)
				}
				exported.NeedsContext = true
				continue
			}

			exported.Params = append(exported.Params, p)
		}
	}

	var returnTypes []FuncReturn
	if decl.Type.Results != nil {
		for _, result := range decl.Type.Results.List {
			for range namesListOrUnnamed(result.Names) {
				r, err := m.parseFuncDeclReturn(result)
				if err != nil {
					return err
				}
				returnTypes = append(returnTypes, r)
			}
		}
	}

	// Special treatment for error returns
	if len(returnTypes) > 0 && returnTypes[len(returnTypes)-1].Type == "error" {
		exported.ReturnsError = true
		returnTypes = returnTypes[:len(returnTypes)-1]
	}

	if len(returnTypes) > 1 {
		return fmt.Errorf("multiple return values are not supported")
	}

	if len(returnTypes) > 0 {
		exported.Return = &returnTypes[0]
	}

	m.addImport(importRisorObject)
	m.addImport("context")
	return m.addExportedFunc(exported)
}

func (m *Module) addExportedFunc(exported ExportedFunc) error {
	for _, f := range m.exportedFuncs {
		if f.FuncGenName == exported.FuncGenName {
			return fmt.Errorf("name collision: multiple functions with generated Go function name: %q", exported.FuncGenName)
		}
		if f.FuncName == exported.FuncName {
			return fmt.Errorf("name collision: multiple generated functions to the same function")
		}
		if f.ExportedName == exported.ExportedName {
			return fmt.Errorf("name collision: multiple generated functions to the same Risor function: %q", exported.ExportedName)
		}
	}
	m.exportedFuncs = append(m.exportedFuncs, exported)
	return nil
}

// namesListOrUnnamed is used to run loops through parameters and results
// at least once per item, so that unnamed parameters and unnamed return values
// are also considered.
func namesListOrUnnamed(names []*ast.Ident) []string {
	if names == nil {
		return []string{"_"}
	}
	var result []string
	for _, name := range names {
		result = append(result, name.Name)
	}
	return result
}

func (m *Module) parseFuncDeclName(decl *ast.FuncDecl) (string, error) {
	if decl.Doc == nil {
		return "", nil
	}

	for _, comment := range decl.Doc.List {
		after, ok := cutPrefixAndSpace(comment.Text, "//risor:export")
		if !ok {
			continue
		}
		fields := strings.Fields(after)
		if len(fields) > 1 {
			pos := m.fset.Position(comment.Pos())
			return "", fmt.Errorf("line %d: too many fields after //risor:export comment", pos.Line)
		}
		if len(fields) == 1 {
			return fields[0], nil
		}
		return strings.ToLower(decl.Name.Name), nil
	}
	return "", nil
}

func (m *Module) parseFuncDeclParam(name string, param *ast.Field) (FuncParam, error) {
	typ := m.sprintExpr(param.Type)
	p, err := m.parseParamType(typ)
	if err != nil {
		return FuncParam{}, err
	}
	p.Name = name
	p.Type = typ
	return p, nil
}

func (m *Module) parseParamType(typeName string) (FuncParam, error) {
	switch typeName {
	case "context.Context":
		// there's special handling of [context.Context] in [parseFuncDecl]
		return FuncParam{}, nil
	case "object.Object":
		// Just pass [object.Object] through, as-is
		return FuncParam{}, nil
	case "any", "interface{}":
		return FuncParam{}, fmt.Errorf("type 'any' is not allowed, use 'object.Object' instead")
	case "string":
		return FuncParam{ReadFunc: "AsString"}, nil
	case "[]string":
		return FuncParam{ReadFunc: "AsStringSlice"}, nil
	case "[]byte":
		return FuncParam{ReadFunc: "AsBytes"}, nil
	case "bool":
		return FuncParam{ReadFunc: "AsBool"}, nil
	case "int64":
		return FuncParam{ReadFunc: "AsInt"}, nil
	case "int32":
		m.addImport("math")
		return FuncParam{ReadFunc: "AsInt", CastFunc: "int32", CastMaxValue: "math.MaxInt32", CastMinValue: "math.MinInt32"}, nil
	case "int":
		m.addImport("math")
		return FuncParam{ReadFunc: "AsInt", CastFunc: "int", CastMaxValue: "math.MaxInt", CastMinValue: "math.MinInt"}, nil
	case "float64":
		return FuncParam{ReadFunc: "AsFloat"}, nil
	case "float32":
		m.addImport("math")
		return FuncParam{ReadFunc: "AsFloat", CastFunc: "float32", CastMaxValue: "math.MaxFloat32"}, nil
	default:
		return FuncParam{}, fmt.Errorf("unsupported parameter type: %q", typeName)
	}
}

func (m *Module) parseFuncDeclReturn(ret *ast.Field) (FuncReturn, error) {
	typ := m.sprintExpr(ret.Type)
	p, err := m.parseReturnType(typ)
	if err != nil {
		return FuncReturn{}, err
	}
	p.Type = typ
	return p, nil
}

func (m *Module) parseReturnType(typeName string) (FuncReturn, error) {
	switch typeName {
	case "error":
		// there's special handling of [error] in [parseFuncDecl]
		return FuncReturn{}, nil
	case "object.Object":
		// Just pass [object.Object] through, as-is
		return FuncReturn{}, nil
	case "any", "interface{}":
		return FuncReturn{}, fmt.Errorf("type 'any' is not allowed, use 'object.Object' instead")
	case "string":
		return FuncReturn{NewFunc: "NewString"}, nil
	case "[]string":
		return FuncReturn{NewFunc: "NewStringList"}, nil
	case "[]byte":
		return FuncReturn{NewFunc: "NewByteSlice"}, nil
	case "bool":
		return FuncReturn{NewFunc: "NewBool"}, nil
	case "int64":
		return FuncReturn{NewFunc: "NewInt"}, nil
	case "int", "int32":
		return FuncReturn{NewFunc: "NewInt", CastFunc: "int64"}, nil
	case "float64":
		return FuncReturn{NewFunc: "NewFloat"}, nil
	case "float32":
		return FuncReturn{NewFunc: "NewFloat", CastFunc: "float64"}, nil
	default:
		return FuncReturn{}, fmt.Errorf("unsupported return type: %q", typeName)
	}
}

func firstRuneUpper(s string) string {
	if s == "" {
		return ""
	}
	var firstRune rune
	for _, r := range s {
		firstRune = r
		break
	}
	firstRune = unicode.ToUpper(firstRune)
	return string(firstRune) + s[1:]
}

func firstRuneLower(s string) string {
	if s == "" {
		return ""
	}
	var firstRune rune
	for _, r := range s {
		firstRune = r
		break
	}
	firstRune = unicode.ToLower(firstRune)
	return string(firstRune) + s[1:]
}
