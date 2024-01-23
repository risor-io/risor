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
	Return       FuncReturn
	Params       []FuncParam
}

type FuncReturn struct {
	NewFunc  string
	CastFunc string
}

type FuncParam struct {
	Name         string
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

	for _, param := range decl.Type.Params.List {
		for _, name := range param.Names {
			p, err := m.parseFuncDeclParam(name.Name, param)
			if err != nil {
				return fmt.Errorf("param %q: %w", name.Name, err)
			}
			exported.Params = append(exported.Params, p)
		}
	}

	if decl.Type.Results != nil {
		if len(decl.Type.Results.List) > 1 {
			return fmt.Errorf("multiple return values are not supported")
		}
		field := decl.Type.Results.List[0]
		if len(field.Names) > 1 {
			return fmt.Errorf("multiple return values are not supported")
		}
		r, err := m.parseFuncDeclReturn(field)
		if err != nil {
			return err
		}
		exported.Return = r
	}

	m.addImport(importRisorObject)
	m.addImport("context")
	m.exportedFuncs = append(m.exportedFuncs, exported)
	return nil
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
	p, err := m.parseParamType(m.sprintExpr(param.Type))
	if err != nil {
		return FuncParam{}, err
	}
	p.Name = name
	return p, nil
}

func (m *Module) parseParamType(typeName string) (FuncParam, error) {
	switch typeName {
	case "string":
		return FuncParam{ReadFunc: "AsString"}, nil
	case "[]string":
		return FuncParam{ReadFunc: "AsStringSlice"}, nil
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
	p, err := m.parseReturnType(m.sprintExpr(ret.Type))
	if err != nil {
		return FuncReturn{}, err
	}
	return p, nil
}

func (m *Module) parseReturnType(typeName string) (FuncReturn, error) {
	switch typeName {
	case "string":
		return FuncReturn{NewFunc: "NewString"}, nil
	case "[]string":
		return FuncReturn{NewFunc: "NewStringList"}, nil
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
