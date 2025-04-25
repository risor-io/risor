package playwright

import (
	"context"
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

// Locator wraps the playwright Locator
type Locator struct {
	locator playwright.Locator
}

func (l *Locator) Type() object.Type {
	return "playwright.locator"
}

func (l *Locator) Inspect() string {
	return "playwright.locator()"
}

func (l *Locator) Interface() interface{} {
	return l.locator
}

func (l *Locator) Equals(other object.Object) object.Object {
	return object.NewBool(l == other)
}

func (l *Locator) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "text_content":
		return object.NewBuiltin("text_content", func(ctx context.Context, args ...object.Object) object.Object {
			text, err := l.locator.TextContent()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewString(text)
		}), true
	case "all":
		return object.NewBuiltin("all", func(ctx context.Context, args ...object.Object) object.Object {
			all, err := l.locator.All()
			if err != nil {
				return object.NewError(err)
			}
			locators := make([]*Locator, len(all))
			for i, loc := range all {
				locators[i] = &Locator{locator: loc}
			}
			return &LocatorArray{locators: locators}
		}), true
	case "click":
		return object.NewBuiltin("click", func(ctx context.Context, args ...object.Object) object.Object {
			if err := l.locator.Click(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "count":
		return object.NewBuiltin("count", func(ctx context.Context, args ...object.Object) object.Object {
			count, err := l.locator.Count()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewInt(int64(count))
		}), true
	case "fill":
		return object.NewBuiltin("fill", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError(fmt.Errorf("fill: expected 1 argument, got %d", len(args)))
			}
			text, ok := args[0].(*object.String)
			if !ok {
				return object.NewError(fmt.Errorf("fill: expected string argument, got %s", args[0].Type()))
			}
			if err := l.locator.Fill(text.Value()); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "press":
		return object.NewBuiltin("press", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError(fmt.Errorf("press: expected 1 argument, got %d", len(args)))
			}
			key, ok := args[0].(*object.String)
			if !ok {
				return object.NewError(fmt.Errorf("press: expected string argument, got %s", args[0].Type()))
			}
			if err := l.locator.Press(key.Value()); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "check":
		return object.NewBuiltin("check", func(ctx context.Context, args ...object.Object) object.Object {
			if err := l.locator.Check(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "uncheck":
		return object.NewBuiltin("uncheck", func(ctx context.Context, args ...object.Object) object.Object {
			if err := l.locator.Uncheck(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "is_checked":
		return object.NewBuiltin("is_checked", func(ctx context.Context, args ...object.Object) object.Object {
			checked, err := l.locator.IsChecked()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewBool(checked)
		}), true
	case "is_visible":
		return object.NewBuiltin("is_visible", func(ctx context.Context, args ...object.Object) object.Object {
			visible, err := l.locator.IsVisible()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewBool(visible)
		}), true
	case "is_disabled":
		return object.NewBuiltin("is_disabled", func(ctx context.Context, args ...object.Object) object.Object {
			disabled, err := l.locator.IsDisabled()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewBool(disabled)
		}), true
	case "is_enabled":
		return object.NewBuiltin("is_enabled", func(ctx context.Context, args ...object.Object) object.Object {
			enabled, err := l.locator.IsEnabled()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewBool(enabled)
		}), true
	case "focus":
		return object.NewBuiltin("focus", func(ctx context.Context, args ...object.Object) object.Object {
			if err := l.locator.Focus(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "hover":
		return object.NewBuiltin("hover", func(ctx context.Context, args ...object.Object) object.Object {
			if err := l.locator.Hover(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "inner_html":
		return object.NewBuiltin("inner_html", func(ctx context.Context, args ...object.Object) object.Object {
			html, err := l.locator.InnerHTML()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewString(html)
		}), true
	case "inner_text":
		return object.NewBuiltin("inner_text", func(ctx context.Context, args ...object.Object) object.Object {
			text, err := l.locator.InnerText()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewString(text)
		}), true
	case "get_attribute":
		return object.NewBuiltin("get_attribute", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError(fmt.Errorf("get_attribute: expected 1 argument, got %d", len(args)))
			}
			attr, ok := args[0].(*object.String)
			if !ok {
				return object.NewError(fmt.Errorf("get_attribute: expected string argument, got %s", args[0].Type()))
			}
			value, err := l.locator.GetAttribute(attr.Value())
			if err != nil {
				return object.NewError(err)
			}
			if value == "" {
				return object.Nil
			}
			return object.NewString(value)
		}), true
	case "first":
		return object.NewBuiltin("first", func(ctx context.Context, args ...object.Object) object.Object {
			first := l.locator.First()
			return &Locator{locator: first}
		}), true
	case "last":
		return object.NewBuiltin("last", func(ctx context.Context, args ...object.Object) object.Object {
			last := l.locator.Last()
			return &Locator{locator: last}
		}), true
	case "nth":
		return object.NewBuiltin("nth", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError(fmt.Errorf("nth: expected 1 argument, got %d", len(args)))
			}
			idx, ok := args[0].(*object.Int)
			if !ok {
				return object.NewError(fmt.Errorf("nth: expected integer argument, got %s", args[0].Type()))
			}
			nth := l.locator.Nth(int(idx.Value()))
			return &Locator{locator: nth}
		}), true
	case "wait_for":
		return object.NewBuiltin("wait_for", func(ctx context.Context, args ...object.Object) object.Object {
			if err := l.locator.WaitFor(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "screenshot":
		return object.NewBuiltin("screenshot", func(ctx context.Context, args ...object.Object) object.Object {
			bytes, err := l.locator.Screenshot()
			if err != nil {
				return object.NewError(err)
			}
			return object.NewByteSlice(bytes)
		}), true
	case "clear":
		return object.NewBuiltin("clear", func(ctx context.Context, args ...object.Object) object.Object {
			if err := l.locator.Clear(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	case "locator":
		return object.NewBuiltin("locator", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError(fmt.Errorf("locator: expected 1 argument, got %d", len(args)))
			}
			selector, ok := args[0].(*object.String)
			if !ok {
				return object.NewError(fmt.Errorf("locator: expected string argument, got %s", args[0].Type()))
			}
			locator := l.locator.Locator(selector.Value())
			return &Locator{locator: locator}
		}), true
	}
	return nil, false
}

func (l *Locator) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("cannot set attribute %q on playwright.locator object", name))
}

func (l *Locator) IsTruthy() bool {
	return true
}

func (l *Locator) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("operation %v not supported on playwright.locator object", opType))
}

func (l *Locator) Cost() int {
	return 0
}
