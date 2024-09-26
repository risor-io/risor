package template

import (
	"context"
	"fmt"
	"strings"
	"text/template"

	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const TEMPLATE object.Type = "template"

type Template struct {
	tpl *template.Template
}

func (t *Template) Type() object.Type {
	return TEMPLATE
}

func (t *Template) Inspect() string {
	return "template"
}

func (t *Template) Interface() interface{} {
	return t.tpl
}

func (t *Template) IsTruthy() bool {
	return true
}

func (t *Template) Cost() int {
	return 8
}

func (t *Template) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal template")
}

func (db *Template) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.EvalErrorf("eval error: unsupported operation for %s: %v", TEMPLATE, opType)
}

func (t *Template) Equals(other object.Object) object.Object {
	if other.Type() != TEMPLATE {
		return object.False
	}

	return object.NewBool(t.tpl == other.(*Template).tpl)
}

func (db *Template) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: %s object has no attribute %q", TEMPLATE, name)
}

func (t *Template) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "delims":
		return object.NewBuiltin("template.delims", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("template.delims", 2, args); err != nil {
				return err
			}

			left, errObj := object.AsString(args[0])
			if errObj != nil {
				return errObj
			}

			right, errObj := object.AsString(args[1])
			if errObj != nil {
				return errObj
			}

			t.tpl.Delims(left, right)

			return object.Nil
		}), true
	case "parse":
		return object.NewBuiltin("template.parse", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("template.parse", 1, args); err != nil {
				return err
			}

			template, argsErr := object.AsString(args[0])
			if argsErr != nil {
				return argsErr
			}

			if _, err := t.tpl.Parse(template); err != nil {
				return object.NewError(err)
			}

			return object.Nil
		}), true
	case "add":
		return object.NewBuiltin("template.add", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("template.add", 2, args); err != nil {
				return err
			}

			name, argsErr := object.AsString(args[0])
			if argsErr != nil {
				return argsErr
			}

			template, argsErr := object.AsString(args[1])
			if argsErr != nil {
				return argsErr
			}

			if _, err := t.tpl.New(name).Parse(template); err != nil {
				return object.NewError(err)
			}

			return object.Nil
		}), true
	case "execute_template":
		return object.NewBuiltin("template.execute_template", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("template.execute_template", 2, args); err != nil {
				return err
			}

			data, argsErr := object.AsMap(args[0])
			if argsErr != nil {
				return argsErr
			}

			name, argsErr := object.AsString(args[1])
			if argsErr != nil {
				return argsErr
			}

			buf := new(strings.Builder)

			if err := t.tpl.ExecuteTemplate(buf, name, data.Interface()); err != nil {
				return object.NewError(err)
			}

			return object.NewString(buf.String())
		}), true
	case "execute":
		return object.NewBuiltin("template.execute", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("template.execute", 1, args); err != nil {
				return err
			}

			data, argsErr := object.AsMap(args[0])
			if argsErr != nil {
				return argsErr
			}

			buf := new(strings.Builder)

			if err := t.tpl.Execute(buf, data.Interface()); err != nil {
				return object.NewError(err)
			}

			return object.NewString(buf.String())
		}), true
	}
	return nil, false
}

func New(ctx context.Context, args ...object.Object) object.Object {
	var name string
	if len(args) > 0 {
		var errObj *object.Error
		name, errObj = object.AsString(args[0])
		if errObj != nil {
			return errObj
		}
	}
	return &Template{
		tpl: newTemplate(name),
	}
}
