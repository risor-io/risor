# playwright

This module provides a Risor wrapper for the [Playwright](https://playwright.dev) browser automation library.
It allows Risor scripts to control browsers for web automation, testing, and scraping.

## Installation

This module requires the Playwright library, which can be installed via npm:

```bash
npm install -g playwright
```

## Usage

```risor
// Initialize Playwright
pw := playwright.run()

// Launch a browser (Chromium, Firefox, or WebKit)
browser := pw.chromium.launch()

// Create a new page
page := browser.new_page()

// Navigate to a URL
page.goto("https://example.com")

// Interact with the page
element := page.locator("#my-element")
text := element.text_content()

// Clean up resources
browser.close()
pw.stop()
```

### Using the Simplified `with_page` Function

```risor
// A more concise approach using with_page
playwright.with_page(func(page) {
    // Navigate to a URL
    page.goto("https://example.com")
    
    // Interact with the page
    element := page.locator("#my-element")
    text := element.text_content()
    
    // No need to manually clean up - it's handled automatically
})
```

## API Reference

### `playwright`

- `playwright.run()`: Initializes the Playwright engine and returns a Playwright instance.
- `playwright.install()`: Installs the Playwright engine and browsers.
- `playwright.with_page(callback, [config])`: Executes a callback with a page and handles setup/cleanup automatically.
  - `callback`: A function that receives a page object.
  - `config` (optional): A map containing configuration options:
    - `playwright`: An existing Playwright instance to use.
    - `browser`: An existing browser instance to use.

### `PlaywrightInstance`

- `pw.chromium`: The Chromium browser type.
- `pw.firefox`: The Firefox browser type.
- `pw.webkit`: The WebKit browser type.
- `pw.stop()`: Stops the Playwright engine.

### `BrowserType`

- `browserType.launch()`: Launches a browser instance.

### `Browser`

- `browser.new_page()`: Creates a new page.
- `browser.close()`: Closes the browser.

### `Page`

- `page.goto(url)`: Navigates to the given URL.
- `page.locator(selector)`: Returns a locator for the given selector.

### `Locator`

- `locator.text_content()`: Returns the text content of the element.
- `locator.all()`: Returns all elements matching the locator.

## Examples

See the `examples` directory for sample scripts.

### Scraping Hacker News Headlines

```risor
// Initialize Playwright
pw := playwright.run()

// Launch a browser
browser := pw.chromium.launch()

// Create a new page
page := browser.new_page()

// Navigate to Hacker News
page.goto("https://news.ycombinator.com")

// Find all story entries
entries := page.locator(".athing").all()

// Display the top 10 headlines
for i := range len(entries) {
    if i >= 10 {
        break
    }
    
    // Get the title of each entry
    title := entries[i].locator("td.title > span > a").text_content()
    print(fmt.sprintf("%d: %s", i+1, title))
}

// Clean up resources
browser.close()
pw.stop()
```

### Using with_page for Simplified Script

```risor
// Simplified approach with with_page
playwright.with_page(func(page) {
    // Navigate to Hacker News
    page.goto("https://news.ycombinator.com")
    
    // Find all story entries
    entries := page.locator(".athing").all()
    
    // Display the top 10 headlines
    for i, entry := range entries {
        if i >= 10 {
            break
        }
        
        // Get the title of each entry
        title := entry.locator("td.title > span > a")
        printf("%d: %s\n", i+1, title.text_content())
    }
    
    // No need for cleanup - it's handled automatically
})
```

## License

This module is available under the same license as Risor. 