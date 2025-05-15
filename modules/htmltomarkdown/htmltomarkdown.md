# htmltomarkdown

The `htmltomarkdown` module supports converting HTML to Markdown.

This module uses the [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) Go library.

## Module

```go copy filename="Function signature"
htmltomarkdown(html string) string
```

The `htmltomarkdown` module object itself is callable in order to provide a shorthand for converting HTML to Markdown:

This is equivalent to calling `htmltomarkdown.convert()` with the same arguments.

## Functions

### convert

```go filename="Function signature"
convert(html string) string
```

Convert HTML content to Markdown format.

```go copy filename="Example"
import htmltomarkdown

html := "<strong>Bold Text</strong>"
md := htmltomarkdown.convert(html)
print(md)  // Outputs: "**Bold Text**"

html = "<h1>Heading</h1><p>Paragraph with <a href='https://example.com'>link</a></p>"
md = htmltomarkdown.convert(html)
print(md)  // Outputs: "# Heading\n\nParagraph with [link](https://example.com)"
```
