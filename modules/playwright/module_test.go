package playwright

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
)

func TestPlaywrightModule(t *testing.T) {
	// Test that the module initializes properly
	mod := Module()
	if mod.Type() != object.MODULE {
		t.Errorf("expected module type to be MODULE, got %s", mod.Type())
	}

	// Test that the run function exists
	runObj, ok := mod.GetAttr("run")
	if !ok {
		t.Fatal("expected 'run' function to exist in module")
	}

	// Verify run is a builtin function
	builtinRun, ok := runObj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'run' to be a builtin function, got %T", runObj)
	}

	// Test that the run function has the expected function signature
	if builtinRun.Name() != "run" {
		t.Errorf("expected function name to be 'run', got %s", builtinRun.Name())
	}

	// Test that the with_page function exists
	withPageObj, ok := mod.GetAttr("with_page")
	if !ok {
		t.Fatal("expected 'with_page' function to exist in module")
	}

	// Verify with_page is a builtin function
	builtinWithPage, ok := withPageObj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'with_page' to be a builtin function, got %T", withPageObj)
	}

	// Test that the with_page function has the expected function signature
	if builtinWithPage.Name() != "with_page" {
		t.Errorf("expected function name to be 'with_page', got %s", builtinWithPage.Name())
	}

	// Note: We don't actually run the function here since it would require
	// Playwright to be installed and would create browser instances
}

func TestPlaywrightTypes(t *testing.T) {
	// Test PlaywrightInstance type
	pw := &PlaywrightInstance{pw: nil, cleanup: func() error { return nil }}
	if pw.Type() != "playwright" {
		t.Errorf("expected PlaywrightInstance.Type() to be 'playwright', got %s", pw.Type())
	}

	// Test BrowserType type
	bt := &BrowserType{}
	if bt.Type() != "playwright.browserType" {
		t.Errorf("expected BrowserType.Type() to be 'playwright.browserType', got %s", bt.Type())
	}

	// Test Browser type
	b := &Browser{}
	if b.Type() != "playwright.browser" {
		t.Errorf("expected Browser.Type() to be 'playwright.browser', got %s", b.Type())
	}

	// Test Page type
	p := &Page{}
	if p.Type() != "playwright.page" {
		t.Errorf("expected Page.Type() to be 'playwright.page', got %s", p.Type())
	}

	// Test Locator type
	l := &Locator{}
	if l.Type() != "playwright.locator" {
		t.Errorf("expected Locator.Type() to be 'playwright.locator', got %s", l.Type())
	}

	// Test LocatorArray type
	la := &LocatorArray{locators: []*Locator{}}
	if la.Type() != "playwright.locatorArray" {
		t.Errorf("expected LocatorArray.Type() to be 'playwright.locatorArray', got %s", la.Type())
	}
}

func TestLocatorArrayIteration(t *testing.T) {
	// Create a LocatorArray with some test items
	la := &LocatorArray{
		locators: []*Locator{
			{},
			{},
			{},
		},
	}

	// Test Len()
	if la.Len().Value() != 3 {
		t.Errorf("expected LocatorArray.Len() to be 3, got %d", la.Len().Value())
	}

	// Test iteration
	it := la.Iter()
	count := 0
	ctx := context.Background()

	for {
		_, ok := it.Next(ctx)
		if !ok {
			break
		}
		entry, ok := it.Entry()
		if !ok {
			t.Fatal("expected iterator entry to be available after Next()")
		}

		// Verify the index matches our count
		index := entry.Key().(*object.Int).Value()
		if index != int64(count) {
			t.Errorf("expected index to be %d, got %d", count, index)
		}
		count++
	}

	if count != 3 {
		t.Errorf("expected to iterate over 3 items, got %d", count)
	}
}
