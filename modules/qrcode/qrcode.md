# qrcode

The `qrcode` module supports easily creating and saving QR codes.

Wraps the [go-qrcode](https://github.com/yeqown/go-qrcode) library.

## Module

```go copy filename="Function signature"
qrcode(content, options=nil)
```

The `qrcode` module object itself is callable to provide a shorthand for
creating a new QR code.

```go copy filename="Example"
>>> qrcode("https://example.com")
qrcode.qrcode()
```

## Functions

### create

```go filename="Function signature"
create(content, options=nil) qrcode
```

Creates a new QR code with the given content. An options map may be provided to configure the QR code.

Available options:
- `encoding_mode`: String specifying the encoding mode - "numeric", "alphanumeric", or "byte"
- `error_correction`: String specifying the error correction level - "low", "medium", "high", or "highest"
- `width`: Integer specifying the QR code width in pixels (default: 40, range: 1-255)

```go filename="Example"
>>> qrcode.create("https://example.com")
qrcode.qrcode()
>>> qrcode.create("12345", {"encoding_mode": "numeric", "error_correction": "high", "width": 80})
qrcode.qrcode()
```

### save

```go filename="Function signature"
save(qr, path, style_options=nil)
```

Saves the QR code to a file at the specified path. A style_options map may be provided to configure the appearance of the QR code.

Available style options:
- `bg_transparent`: (bool) make the background transparent
- `bg_color_hex`: (string) set background color using hex color code (e.g. "#FFFFFF")
- `fg_color_hex`: (string) set foreground color using hex color code (e.g. "#000000")
- `logo_image`: (image) a Risor image object to use as a logo in the center
- `shape`: (string) "circle" or "rectangle" (default: "rectangle")
- `border_width`: (int) width of the border around the QR code
- `format`: (string) "png" or "jpeg" (default: "png")
- `width`: (int) override the width setting from code creation

```go filename="Example"
>>> qr := qrcode.create("https://example.com")
>>> qrcode.save(qr, "example.png")
nil
>>> qrcode.save(qr, "example-large.png", {"width": 200})
nil
>>> qrcode.save(qr, "example-styled.png", {"bg_color_hex": "#F5F5F5", "fg_color_hex": "#0078D7"})
nil
```

## Types

### qrcode

The qrcode object represents a generated QR code that can be saved to a file or converted to other formats.

The qrcode object has the following methods and properties:

#### save

```go filename="Method signature"
qrcode.save(path, style_options=nil) nil
```

Saves the QR code to a file at the specified path. Accepts the same style_options as the module-level save function.

```go filename="Example"
>>> qr := qrcode.create("https://example.com")
>>> qr.save("example.png")
nil
>>> qr.save("example_styled.png", {
...     "bg_color_hex": "#F5F5F5",
...     "fg_color_hex": "#0078D7",
...     "shape": "circle",
...     "border_width": 10
... })
nil
```

#### dimension

```go filename="Method signature"
qrcode.dimension() int
```

Returns the dimension (width/height) of the QR code in modules (the small squares that make up a QR code).

```go filename="Example"
>>> qr := qrcode.create("https://example.com")
>>> qr.dimension()
25
```

#### bytes

```go filename="Method signature"
qrcode.bytes(style_options=nil) byte_slice
```

Returns the QR code as a byte slice. A style_options map may be provided to configure the appearance (same options as save).

```go filename="Example"
>>> qr := qrcode.create("https://example.com")
>>> qr.bytes()
byte_slice([...])
>>> qr.bytes({"bg_transparent": true})
byte_slice([...])
```

#### base64

```go filename="Method signature"
qrcode.base64(style_options=nil) string
```

Returns the QR code as a base64-encoded image string. A style_options map may be provided to configure the appearance.
This is useful for embedding QR codes in HTML or other formats.

```go filename="Example"
>>> qr := qrcode.create("https://example.com")
>>> base64_str := qr.base64()
>>> base64_str[0:20]  // Show first 20 characters
"iVBORw0KGgoAAAANSU"
>>> transparent_base64 := qr.base64({"bg_transparent": true})
```

#### width

```go filename="Property"
qrcode.width int
```

Returns the configured width of the QR code in pixels.

```go filename="Example"
>>> qr := qrcode.create("https://example.com", {"width": 100})
>>> qr.width
100
```

## Examples

### Basic QR Code

```risor
qr := qrcode.create("https://example.com")
qr.save("example.png")
```

### QR Code with Options

```risor
qr := qrcode.create("https://example.com", {
    "encoding_mode": "byte",
    "error_correction": "high",
    "width": 80
})
qr.save("example.png")
```

### QR Code with Styling

```risor
qr := qrcode.create("https://example.com")
qr.save("example.png", {
    "bg_color_hex": "#F5F5F5",
    "fg_color_hex": "#0078D7",
    "shape": "circle",
    "border_width": 10
})
```

### QR Code with Logo

```risor
logo := image.decode(open("logo.png"))
qr.save("example_with_image_logo.png", {
    "logo_image": logo
})
```

### Get QR Code as Base64

```risor
qr := qrcode.create("https://example.com")
base64Data := qr.base64({
    "bg_transparent": true
})
print(base64Data)
```

## Complete Example

Here's a complete example showing how to create a QR code and save it or convert it to base64:

```go
// Create a QR code for a website
qr := qrcode.create("https://example.com", {
    "encoding_mode": "byte",
    "error_correction": "high",
    "width": 80
})

// Save to a file using the object method
qr.save("example.png")

// Save with styling options
qr.save("example_styled.png", {
    "bg_color_hex": "#F5F5F5",
    "fg_color_hex": "#0078D7",
    "border_width": 5
})

// Get the dimension
dim := qr.dimension()
print("QR code dimension:", dim)

// Get the configured width
width := qr.width
print("QR code width:", width)

// Get raw bytes
raw_bytes := qr.bytes()
print("QR code bytes length:", len(raw_bytes))

// Get base64 representation
base64_data := qr.base64()
print("data:image/png;base64," + base64_data)
```
