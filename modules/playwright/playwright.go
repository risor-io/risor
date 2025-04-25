package playwright

import (
	"context"
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

// PlaywrightInstance wraps the playwright.Playwright instance
type PlaywrightInstance struct {
	pw      *playwright.Playwright
	cleanup func() error
}

func (p *PlaywrightInstance) Type() object.Type { return "playwright" }

func (p *PlaywrightInstance) Inspect() string {
	return "playwright()"
}

func (p *PlaywrightInstance) Interface() interface{} {
	return p.pw
}

func (p *PlaywrightInstance) Equals(other object.Object) object.Object {
	return object.NewBool(p == other)
}

func (p *PlaywrightInstance) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "chromium":
		return &BrowserType{browserType: p.pw.Chromium}, true
	case "firefox":
		return &BrowserType{browserType: p.pw.Firefox}, true
	case "webkit":
		return &BrowserType{browserType: p.pw.WebKit}, true
	case "stop":
		return object.NewBuiltin("stop", func(ctx context.Context, args ...object.Object) object.Object {
			if err := p.cleanup(); err != nil {
				return object.NewError(err)
			}
			return object.Nil
		}), true
	}
	return nil, false
}

// Stop stops the Playwright instance and cleans up resources
func (p *PlaywrightInstance) Stop() object.Object {
	if err := p.cleanup(); err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

func (p *PlaywrightInstance) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("cannot set attribute %q on playwright object", name))
}

func (p *PlaywrightInstance) IsTruthy() bool {
	return p.pw != nil
}

func (p *PlaywrightInstance) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("operation %v not supported on playwright object", opType))
}

func (p *PlaywrightInstance) Cost() int {
	return 0
}
