# goquery

The `goquery` module provides a convenient way to parse and query HTML documents using CSS selectors. It is a wrapper around the [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery) library, which implements features similar to jQuery for parsing and manipulating HTML.

## Module

```go copy filename="Function signature"
goquery(input ...object)
```

The `goquery` module object itself is callable in order to provide a shorthand for creating a goquery document from various sources:

- A string input: Parses the given HTML string into a document.
- A byte slice input: Parses the HTML content from the given bytes.
- A file or reader input: Creates a document from the content of the file or reader.

This is equivalent to calling `goquery.parse()` with the same arguments.

```go copy filename="Example"
>>> doc := goquery("<html><body><h1>Hello World</h1></body></html>")
goquery.document()
>>> doc := goquery(bytes("<html><body><h1>Hello</h1></body></html>"))
goquery.document()
>>> f := file.open("page.html")
>>> doc := goquery(f)
goquery.document()
```

## Functions

### parse

```go filename="Function signature"
parse(input) document
```

Creates a goquery document by parsing the given input. The input can be:
- A string containing HTML
- A byte slice containing HTML
- A file or reader object

```go filename="Example"
>>> doc := goquery.parse("<html><body><h1>Hello World</h1></body></html>")
goquery.document()
>>> doc := goquery.parse(bytes("<html><body><h1>Hello</h1></body></html>"))
goquery.document()
>>> f := open("page.html")
>>> doc := goquery.parse(f)
goquery.document()
```

## Types

### document

The document object represents an HTML document and provides methods for finding and manipulating its contents.

#### Methods

##### find

```go filename="Method signature"
find(selector string) selection
```

Finds elements in the document that match the given CSS selector and returns them as a selection.

```go filename="Example"
>>> doc := goquery("<html><body><h1>Hello</h1><p>World</p></body></html>")
>>> doc.find("h1")
goquery.selection()
```

##### html

```go filename="Method signature"
html() string
```

Returns the HTML content of the document.

```go filename="Example"
>>> doc := goquery("<html><body><h1>Hello</h1></body></html>")
>>> doc.html()
"<html><head></head><body><h1>Hello</h1></body></html>"
```

##### text

```go filename="Method signature"
text() string
```

Returns the text content of the document, with all HTML tags removed.

```go filename="Example"
>>> doc := goquery("<html><body><h1>Hello</h1><p>World</p></body></html>")
>>> doc.text()
"HelloWorld"
```

### selection

The selection object represents a set of HTML elements selected from a document. It provides methods for filtering, traversing, and manipulating these elements.

#### Methods

##### find

```go filename="Method signature"
find(selector string) selection
```

Finds elements within the current selection that match the given CSS selector.

```go filename="Example"
>>> doc := goquery("<div><p>First</p><p>Second</p></div>")
>>> div := doc.find("div")
>>> div.find("p")
goquery.selection()
```

##### attr

```go filename="Method signature"
attr(name string) string
```

Returns the value of the specified attribute for the first element in the selection. Returns nil if the attribute doesn't exist.

```go filename="Example"
>>> doc := goquery("<a href='https://example.com'>Example</a>")
>>> doc.find("a").attr("href")
"https://example.com"
```

##### html

```go filename="Method signature"
html() string
```

Returns the HTML content of the first element in the selection.

```go filename="Example"
>>> doc := goquery("<div><p>Hello <b>World</b></p></div>")
>>> doc.find("p").html()
"Hello <b>World</b>"
```

##### text

```go filename="Method signature"
text() string
```

Returns the combined text content of all elements in the selection.

```go filename="Example"
>>> doc := goquery("<div><p>Hello <b>World</b></p><p>Example</p></div>")
>>> doc.find("p").text()
"Hello WorldExample"
```

##### each

```go filename="Method signature"
each(func(index int, selection selection) object) object
```

Iterates over each element in the selection and calls the provided function with the index and a selection containing only that element.

```go filename="Example"
>>> doc := goquery("<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul>")
>>> doc.find("li").each(func(i, s) { print(i, s.text()) })
0 Item 1
1 Item 2
2 Item 3
```

##### eq

```go filename="Method signature"
eq(index int) selection
```

Returns a selection containing only the element at the specified index within the current selection.

```go filename="Example"
>>> doc := goquery("<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul>")
>>> doc.find("li").eq(1).text()
"Item 2"
```

##### length

```go filename="Method signature"
length() int
```

Returns the number of elements in the selection.

```go filename="Example"
>>> doc := goquery("<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul>")
>>> doc.find("li").length()
3
```

##### first

```go filename="Method signature"
first() selection
```

Returns a selection containing only the first element in the current selection.

```go filename="Example"
>>> doc := goquery("<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul>")
>>> doc.find("li").first().text()
"Item 1"
```

##### last

```go filename="Method signature"
last() selection
```

Returns a selection containing only the last element in the current selection.

```go filename="Example"
>>> doc := goquery("<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul>")
>>> doc.find("li").last().text()
"Item 3"
```

##### parent

```go filename="Method signature"
parent() selection
```

Returns a selection containing the parent elements of each element in the current selection.

```go filename="Example"
>>> doc := goquery("<div><p><span>Text</span></p></div>")
>>> span := doc.find("span")
>>> span.parent().html()
"<span>Text</span>"
```

##### children

```go filename="Method signature"
children(selector ...string) selection
```

Returns a selection containing the children of each element in the current selection, optionally filtered by a selector.

```go filename="Example"
>>> doc := goquery("<div><p>First</p><span>Second</span><p>Third</p></div>")
>>> div := doc.find("div")
>>> div.children().length()
3
>>> div.children("p").length()
2
```

##### filter

```go filename="Method signature"
filter(selector string) selection
```

Returns a filtered selection containing only the elements that match the given selector.

```go filename="Example"
>>> doc := goquery("<div><p class='a'>First</p><p>Second</p><p class='a'>Third</p></div>")
>>> paragraphs := doc.find("p")
>>> paragraphs.filter(".a").length()
2
```

##### not

```go filename="Method signature"
not(selector string) selection
```

Returns a filtered selection containing only the elements that do not match the given selector.

```go filename="Example"
>>> doc := goquery("<div><p class='a'>First</p><p>Second</p><p class='a'>Third</p></div>")
>>> paragraphs := doc.find("p")
>>> paragraphs.not(".a").length()
1
```

##### has_class

```go filename="Method signature"
has_class(class string) bool
```

Returns whether the first element in the selection has the given class.

```go filename="Example"
>>> doc := goquery("<div class='container main'><p>Content</p></div>")
>>> div := doc.find("div")
>>> div.has_class("container")
true
>>> div.has_class("sidebar")
false
``` 