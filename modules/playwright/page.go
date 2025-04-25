package playwright

import (
	"context"
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

// Page wraps the playwright Page
type Page struct {
	page playwright.Page
}

func (p *Page) Type() object.Type { return "playwright.page" }

func (p *Page) Inspect() string {
	return fmt.Sprintf("playwright.page(%p)", &p.page)
}

func (p *Page) Interface() interface{} {
	return p.page
}

func (p *Page) Equals(other object.Object) object.Object {
	if otherPage, ok := other.(*Page); ok {
		return object.NewBool(&p.page == &otherPage.page)
	}
	return object.False
}

func (p *Page) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "goto":
		return object.NewBuiltin("goto", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) < 1 {
				return object.NewArgsError("page.goto", 1, len(args))
			}
			url, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			_, gotoErr := p.page.Goto(url)
			if gotoErr != nil {
				return object.NewError(gotoErr)
			}
			// TODO: Wrap response object
			return object.Nil
		}), true
	case "locator":
		return object.NewBuiltin("locator", func(ctx context.Context, args ...object.Object) object.Object {
			if len(args) < 1 {
				return object.NewArgsError("page.locator", 1, len(args))
			}
			selector, err := object.AsString(args[0])
			if err != nil {
				return err
			}
			locator := p.page.Locator(selector)
			return &Locator{locator: locator}
		}), true
	case "close":
		return object.NewBuiltin("close", func(ctx context.Context, args ...object.Object) object.Object {
			return p.Close()
		}), true
	}
	return nil, false
}

// Goto navigates the page to the given URL
func (p *Page) Goto(url string) object.Object {
	response, err := p.page.Goto(url)
	if err != nil {
		return object.NewError(err)
	}
	if response.Status() > 399 {
		return object.NewError(fmt.Errorf("page.goto: failed to navigate to %s: %s", url, response.StatusText()))
	}
	return object.Nil
}

func (p *Page) Close() object.Object {
	err := p.page.Close()
	if err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func (p *Page) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("cannot set attribute %q on playwright.page object", name))
}

func (p *Page) IsTruthy() bool {
	return true
}

func (p *Page) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("operation %v not supported on playwright.page object", opType))
}

func (p *Page) Cost() int {
	return 0
}
