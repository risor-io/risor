package playwright

import (
	"context"
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

// BrowserType wraps the playwright BrowserType (like Chromium, Firefox, WebKit)
type BrowserType struct {
	browserType playwright.BrowserType
}

func (b *BrowserType) Type() object.Type {
	return "playwright.browser_type"
}

func (b *BrowserType) Inspect() string {
	return fmt.Sprintf("playwright.browser_type(%p)", &b.browserType)
}

func (b *BrowserType) Interface() interface{} {
	return b.browserType
}

func (b *BrowserType) Equals(other object.Object) object.Object {
	if otherBrowser, ok := other.(*BrowserType); ok {
		return object.NewBool(&b.browserType == &otherBrowser.browserType)
	}
	return object.False
}

func (b *BrowserType) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "launch":
		return object.NewBuiltin("launch", func(ctx context.Context, args ...object.Object) object.Object {
			// TODO: Handle launch options
			browser, err := b.browserType.Launch()
			if err != nil {
				return object.NewError(err)
			}
			return &Browser{browser: browser}
		}), true
	}
	return nil, false
}

// Launch launches a browser of this type
func (b *BrowserType) Launch() object.Object {
	browser, err := b.browserType.Launch()
	if err != nil {
		return object.NewError(err)
	}
	return &Browser{browser: browser}
}

func (b *BrowserType) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("cannot set attribute %q on playwright.browser_type object", name))
}

func (b *BrowserType) IsTruthy() bool {
	return true
}

func (b *BrowserType) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("operation %v not supported on playwright.browser_type object", opType))
}

func (b *BrowserType) Cost() int {
	return 0
}
