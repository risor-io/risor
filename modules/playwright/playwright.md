# playwright

The `playwright` module provides a Risor wrapper for the [Playwright](https://playwright.dev) browser automation library.
It allows Risor scripts to control browsers for web automation, testing, and scraping.

## Functions

### run

```go filename="Function signature"
run() playwright
```

Initializes the Playwright engine and returns a Playwright instance.

```go filename="Example"
>>> playwright.run()
playwright()
```

### install

```go filename="Function signature"
install()
```

Installs the Playwright engine and browsers.

```go filename="Example"
>>> playwright.install()
```

### goto

```go filename="Function signature"
goto(url string, callback function, config ...map) any
```

Navigates to a URL, executes a callback with a page, and handles setup/cleanup automatically.

- `url`: The URL to navigate to.
- `callback`: A function that receives a page object.
- `config` (optional): A map containing configuration options:
  - `playwright`: An existing Playwright instance to use.
  - `browser`: An existing browser instance to use.

```go filename="Example"
>>> playwright.goto("https://example.com", func(page) {
...     element := page.locator("#my-element")
...     return element.text_content()
... })
"Example content"
```

## Types

### PlaywrightInstance

The PlaywrightInstance object represents a Playwright engine instance and provides access to browser types.

#### Methods

##### chromium

```go filename="Property"
chromium browserType
```

The Chromium browser type.

```go filename="Example"
>>> pw := playwright.run()
>>> pw.chromium
playwright.browserType()
```

##### firefox

```go filename="Property"
firefox browserType
```

The Firefox browser type.

```go filename="Example"
>>> pw := playwright.run()
>>> pw.firefox
playwright.browserType()
```

##### webkit

```go filename="Property"
webkit browserType
```

The WebKit browser type.

```go filename="Example"
>>> pw := playwright.run()
>>> pw.webkit
playwright.browserType()
```

##### stop

```go filename="Method signature"
stop()
```

Stops the Playwright engine and releases resources.

```go filename="Example"
>>> pw := playwright.run()
>>> pw.stop()
nil
```

### BrowserType

The BrowserType object represents a browser type and provides methods to launch browser instances.

#### Methods

##### launch

```go filename="Method signature"
launch() browser
```

Launches a browser instance.

```go filename="Example"
>>> pw := playwright.run()
>>> browser := pw.chromium.launch()
playwright.browser()
```

### Browser

The Browser object represents a browser instance and provides methods to create pages and manage the browser.

#### Methods

##### new_page

```go filename="Method signature"
new_page() page
```

Creates a new browser page.

```go filename="Example"
>>> pw := playwright.run()
>>> browser := pw.chromium.launch()
>>> page := browser.new_page()
playwright.page()
```

##### close

```go filename="Method signature"
close()
```

Closes the browser and releases resources.

```go filename="Example"
>>> pw := playwright.run()
>>> browser := pw.chromium.launch()
>>> browser.close()
nil
```

### Page

The Page object represents a browser page and provides methods to interact with web content.

#### Methods

##### goto

```go filename="Method signature"
goto(url string)
```

Navigates the page to the given URL.

```go filename="Example"
>>> pw := playwright.run()
>>> browser := pw.chromium.launch()
>>> page := browser.new_page()
>>> page.goto("https://example.com")
nil
```

##### locator

```go filename="Method signature"
locator(selector string) locator
```

Returns a locator for the given CSS or XPath selector.

```go filename="Example"
>>> pw := playwright.run()
>>> browser := pw.chromium.launch()
>>> page := browser.new_page()
>>> page.goto("https://example.com")
>>> element := page.locator("h1")
playwright.locator()
```

##### close

```go filename="Method signature"
close()
```

Closes the page.

```go filename="Example"
>>> pw := playwright.run()
>>> browser := pw.chromium.launch()
>>> page := browser.new_page()
>>> page.close()
nil
```

### Locator

The Locator object represents an element on the page and provides methods to interact with it.

#### Methods

##### text_content

```go filename="Method signature"
text_content() string
```

Returns the text content of the element.

```go filename="Example"
>>> pw := playwright.run()
>>> browser := pw.chromium.launch()
>>> page := browser.new_page()
>>> page.goto("https://example.com")
>>> element := page.locator("h1")
>>> element.text_content()
"Example Domain"
```

##### all

```go filename="Method signature"
all() locatorArray
```

Returns all elements matching the locator as a locator array.

```go filename="Example"
>>> pw := playwright.run()
>>> browser := pw.chromium.launch()
>>> page := browser.new_page()
>>> page.goto("https://example.com")
>>> elements := page.locator("a").all()
playwright.locatorArray()
```

## Examples

### Scraping Hacker News Headlines

```go filename="Example"
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
    printf("%d: %s\n", i+1, title)
}

// Clean up resources
browser.close()
pw.stop()
```

### Using goto for Simplified Script

```go filename="Example"
// Simplified approach with goto
playwright.goto("https://news.ycombinator.com", func(page) {
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
