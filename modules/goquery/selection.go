package goquery

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

var _ object.Object = (*Selection)(nil)

const SELECTION object.Type = "goquery.selection"

type Selection struct {
	value *goquery.Selection
}

func (s *Selection) Value() *goquery.Selection {
	return s.value
}

func (s *Selection) Type() object.Type {
	return SELECTION
}

func (s *Selection) Inspect() string {
	return fmt.Sprintf("%s()", SELECTION)
}

func (s *Selection) IsTruthy() bool {
	return s.value != nil && s.value.Length() > 0
}

func (s *Selection) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: cannot set %q on %s object", name, SELECTION)
}

func (s *Selection) String() string {
	if s.value == nil {
		return "nil"
	}
	html, err := s.value.Html()
	if err != nil {
		return fmt.Sprintf("%s()", SELECTION)
	}
	return html
}

func (s *Selection) Interface() interface{} {
	return s.value
}

func (s *Selection) Equals(other object.Object) object.Object {
	otherSel, ok := other.(*Selection)
	if !ok {
		return object.False
	}
	// Pointer comparison is the best we can do
	return object.NewBool(otherSel.Value() == s.value)
}

func (s *Selection) Cost() int {
	return 0
}

func (s *Selection) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for %s: %v", SELECTION, opType)
}

func (s *Selection) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "find":
		return object.NewBuiltin("find", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.find", 1, args); err != nil {
				return err
			}
			selector, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return NewSelection(s.value.Find(selector))
		}), true
	case "attr":
		return object.NewBuiltin("attr", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.attr", 1, args); err != nil {
				return err
			}
			attrName, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			val, exists := s.value.Attr(attrName)
			if !exists {
				return object.Nil
			}
			return object.NewString(val)
		}), true
	case "html":
		return object.NewBuiltin("html", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.html", 0, args); err != nil {
				return err
			}
			html, err := s.value.Html()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewString(html)
		}), true
	case "text":
		return object.NewBuiltin("text", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.RequireRange("goquery.selection.text", 0, 1, args); err != nil {
				return err
			}
			if len(args) == 0 {
				return object.NewString(s.value.Text())
			}
			selObj, ok := args[0].(*Selection)
			if !ok {
				return object.TypeErrorf("type error: expected selection (got %s)", args[0].Type())
			}
			return object.NewString(selObj.Value().Text())
		}), true
	case "each":
		return object.NewBuiltin("each", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.each", 1, args); err != nil {
				return err
			}

			var result object.Object = object.Nil
			s.value.Each(func(i int, sel *goquery.Selection) {
				fargs := []object.Object{object.NewInt(int64(i)), NewSelection(sel)}

				switch fn := args[0].(type) {
				case *object.Function:
					callFunc, found := object.GetCallFunc(ctx)
					if !found {
						result = object.EvalErrorf("eval error: context did not contain a call function")
						return
					}
					var err error
					result, err = callFunc(ctx, fn, fargs)
					if err != nil {
						result = object.NewError(err)
						return
					}
				case object.Callable:
					result = fn.Call(ctx, fargs...)
				default:
					result = object.TypeErrorf("type error: expected function or callable (got %s)", args[0].Type())
					return
				}

				// Check if we should abort early due to error
				if object.IsError(result) {
					return
				}
			})
			return result
		}), true
	case "eq":
		return object.NewBuiltin("eq", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.eq", 1, args); err != nil {
				return err
			}
			index, err := object.AsInt(args[0])
			if err != nil {
				return err
			}
			return NewSelection(s.value.Eq(int(index)))
		}), true
	case "length":
		return object.NewBuiltin("length", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.length", 0, args); err != nil {
				return err
			}
			return object.NewInt(int64(s.value.Length()))
		}), true
	case "first":
		return object.NewBuiltin("first", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.first", 0, args); err != nil {
				return err
			}
			return NewSelection(s.value.First())
		}), true
	case "last":
		return object.NewBuiltin("last", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.last", 0, args); err != nil {
				return err
			}
			return NewSelection(s.value.Last())
		}), true
	case "parent":
		return object.NewBuiltin("parent", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.parent", 0, args); err != nil {
				return err
			}
			return NewSelection(s.value.Parent())
		}), true
	case "children":
		return object.NewBuiltin("children", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.RequireRange("goquery.selection.children", 0, 1, args); err != nil {
				return err
			}

			if len(args) == 0 {
				return NewSelection(s.value.Children())
			}

			selector, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return NewSelection(s.value.Children().Filter(selector))
		}), true
	case "filter":
		return object.NewBuiltin("filter", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.filter", 1, args); err != nil {
				return err
			}
			selector, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return NewSelection(s.value.Filter(selector))
		}), true
	case "not":
		return object.NewBuiltin("not", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.not", 1, args); err != nil {
				return err
			}
			selector, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return NewSelection(s.value.Not(selector))
		}), true
	case "has_class":
		return object.NewBuiltin("has_class", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("goquery.selection.has_class", 1, args); err != nil {
				return err
			}
			class, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			return object.NewBool(s.value.HasClass(class))
		}), true
	default:
		return nil, false
	}
}

func NewSelection(sel *goquery.Selection) *Selection {
	return &Selection{value: sel}
}
