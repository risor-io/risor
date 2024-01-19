package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type ExportedFunc struct {
	ExportedName string
	FuncName     string
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
	switch expr := param.Type.(type) {
	case *ast.Ident:
		p, err := m.parseParamType(expr.Name)
		if err != nil {
			return FuncParam{}, err
		}
		p.Name = name
		return p, nil
	default:
		return FuncParam{}, fmt.Errorf("unsupported parameter expression type: %T", param.Type)
	}
}

func (m *Module) parseParamType(typeName string) (FuncParam, error) {
	switch typeName {
	case "string":
		return FuncParam{ReadFunc: "AsString"}, nil
	case "bool":
		return FuncParam{ReadFunc: "AsBool"}, nil
	case "int64":
		return FuncParam{ReadFunc: "AsInt"}, nil
	case "int32":
		return FuncParam{ReadFunc: "AsInt", CastFunc: "int32", CastMaxValue: "math.MaxInt32"}, nil
	case "int":
		return FuncParam{ReadFunc: "AsInt", CastFunc: "int", CastMaxValue: "math.MaxInt"}, nil
	case "float64":
		return FuncParam{ReadFunc: "AsFloat"}, nil
	case "float32":
		return FuncParam{ReadFunc: "AsFloat", CastFunc: "float32", CastMaxValue: "math.MaxFloat32"}, nil
	default:
		return FuncParam{}, fmt.Errorf("unsupported parameter type: %q", typeName)
	}
}

func (m *Module) parseFuncDeclReturn(ret *ast.Field) (FuncReturn, error) {
	switch expr := ret.Type.(type) {
	case *ast.Ident:
		p, err := m.parseReturnType(expr.Name)
		if err != nil {
			return FuncReturn{}, err
		}
		return p, nil
	default:
		return FuncReturn{}, fmt.Errorf("unsupported return expression type: %T", ret.Type)
	}
}

func (m *Module) parseReturnType(typeName string) (FuncReturn, error) {
	switch typeName {
	case "string":
		return FuncReturn{NewFunc: "NewString"}, nil
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
