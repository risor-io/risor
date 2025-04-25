package playwright

import (
	"context"
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
)

// Run initializes the Playwright engine and returns a Playwright instance
func Run(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("playwright.run", 0, args); err != nil {
		return err
	}
	pw, err := playwright.Run()
	if err != nil {
		return object.NewError(err)
	}

	return &PlaywrightInstance{
		pw:      pw,
		cleanup: pw.Stop,
	}
}

// Install installs the Playwright engine and browsers
func Install(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("playwright.install", 0, args); err != nil {
		return err
	}
	err := playwright.Install()
	if err != nil {
		return object.NewError(err)
	}
	return object.Nil
}

// Goto initializes Playwright, launches a browser, creates a page,
// navigates to the specified URL, calls the provided callback function with the page,
// and automatically handles cleanup
func Goto(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("playwright.goto", 2, 3, args); err != nil {
		return err
	}

	// First argument must be a string (the URL)
	urlObj, ok := args[0].(*object.String)
	if !ok {
		return object.NewError(fmt.Errorf("playwright.goto: first argument must be a URL string"))
	}
	url := urlObj.Value()

	// Second argument must be a function (the callback)
	fn, ok := args[1].(*object.Function)
	if !ok {
		return object.NewError(fmt.Errorf("playwright.goto: second argument must be a function"))
	}

	// Check if we have a config map as the third argument
	var config *object.Map
	if len(args) > 2 {
		var ok bool
		config, ok = args[2].(*object.Map)
		if !ok {
			return object.NewError(fmt.Errorf("playwright.goto: third argument must be a map"))
		}
	}

	// Setup: initialize Playwright, launch browser, create page
	var pw *PlaywrightInstance
	var browser *Browser
	var existingPw bool
	var existingBrowser bool

	// Check if Playwright instance is provided in config
	if config != nil {
		pwObj := config.Get("playwright")
		if pwObj != nil {
			if pwInst, ok := pwObj.(*PlaywrightInstance); ok {
				pw = pwInst
				existingPw = true
			} else {
				return object.NewError(fmt.Errorf("playwright.goto: config['playwright'] must be a playwright instance"))
			}
		}

		// Check if browser is provided in config
		browserObj := config.Get("browser")
		if browserObj != nil {
			if b, ok := browserObj.(*Browser); ok {
				browser = b
				existingBrowser = true
			} else {
				return object.NewError(fmt.Errorf("playwright.goto: config['browser'] must be a browser instance"))
			}
		}
	}

	// Initialize Playwright if not provided
	if pw == nil {
		pwObj := Run(ctx)
		if err, ok := pwObj.(*object.Error); ok {
			return err
		}
		pw = pwObj.(*PlaywrightInstance)
	}

	// Launch browser if not provided
	if browser == nil {
		browserType, ok := pw.GetAttr("chromium")
		if !ok {
			return object.NewError(fmt.Errorf("playwright.goto: failed to get chromium browser type"))
		}

		browserTypeObj, ok := browserType.(*BrowserType)
		if !ok {
			return object.NewError(fmt.Errorf("playwright.goto: 'chromium' is not a browser type, got %T", browserType))
		}

		browserObj := browserTypeObj.Launch()
		if err, ok := browserObj.(*object.Error); ok {
			// Clean up Playwright if we created it
			if !existingPw {
				pw.Stop()
			}
			return err
		}
		browser = browserObj.(*Browser)
	}

	// Create a new page
	pageObj := browser.NewPage()
	if err, ok := pageObj.(*object.Error); ok {
		// Clean up browser and Playwright if we created them
		if !existingBrowser {
			browser.Close()
		}
		if !existingPw {
			pw.Stop()
		}
		return err
	}
	page := pageObj.(*Page)

	// Navigate to the URL
	gotoResult := page.Goto(url)
	if err, ok := gotoResult.(*object.Error); ok {
		// Clean up page, browser, and Playwright
		if !existingBrowser {
			browser.Close()
		}
		if !existingPw {
			pw.Stop()
		}
		return err
	}

	// Call the callback function with the page
	result := fn.Call(ctx, page)

	// Cleanup: close page, browser, and Playwright
	// Always close the page, even if it was created from an existing browser
	page.Close()
	if !existingBrowser {
		browser.Close()
	}
	if !existingPw {
		pw.Stop()
	}

	return result
}

// Module creates and returns the Playwright module
func Module() *object.Module {
	return object.NewBuiltinsModule("playwright", map[string]object.Object{
		"run":     object.NewBuiltin("run", Run),
		"install": object.NewBuiltin("install", Install),
		"goto":    object.NewBuiltin("goto", Goto),
	}, Run)
}
