package playwright

import (
	"context"
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

type Browser struct {
	browser playwright.Browser
}

func (b *Browser) Type() object.Type {
	return "playwright.browser"
}

func (b *Browser) Inspect() string {
	return fmt.Sprintf("playwright.browser(%p)", &b.browser)
}

func (b *Browser) Interface() interface{} {
	return b.browser
}

func (b *Browser) Equals(other object.Object) object.Object {
	if otherBrowser, ok := other.(*Browser); ok {
		return object.NewBool(&b.browser == &otherBrowser.browser)
	}
	return object.False
}

func (b *Browser) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "new_page":
		return object.NewBuiltin("new_page", func(ctx context.Context, args ...object.Object) object.Object {
			page, err := b.browser.NewPage()
			if err != nil {
				return object.NewError(err)
			}
			return &Page{page: page}
		}), true
	case "close":
		return object.NewBuiltin("close", func(ctx context.Context, args ...object.Object) object.Object {
			if err := b.browser.Close(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	}
	return nil, false
}

// NewPage creates a new browser page
func (b *Browser) NewPage() object.Object {
	page, err := b.browser.NewPage()
	if err != nil {
		return object.NewError(err)
	}
	return &Page{page: page}
}

// Close closes the browser
func (b *Browser) Close() object.Object {
	if err := b.browser.Close(); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func (b *Browser) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("cannot set attribute %q on playwright.browser object", name))
}

func (b *Browser) IsTruthy() bool {
	return true
}

func (b *Browser) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("operation %v not supported on playwright.browser object", opType))
}

func (b *Browser) Cost() int {
	return 0
}
