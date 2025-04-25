# qrcode

The `qrcode` module provides a Risor wrapper for the [go-qrcode](https://github.com/yeqown/go-qrcode) library.
It allows Risor scripts to generate QR codes.

## Module

```go copy filename="Function signature"
qrcode(content, options=nil)
```

The `qrcode` module object itself is callable to provide a shorthand for
creating a new QR code.

```go copy filename="Example"
>>> qrcode("https://example.com")
qrcode.qrcode(...)
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

```go filename="Example"
>>> qrcode.create("https://example.com")
qrcode.qrcode(...)
>>> qrcode.create("12345", {"encoding_mode": "numeric", "error_correction": "high"})
qrcode.qrcode(...)
```

### save

```go filename="Function signature"
save(qr, path, width=256)
```

Saves the QR code to a PNG file at the specified path. An optional width parameter can be specified 
(defaults to 256 pixels).

```go filename="Example"
>>> qr := qrcode.create("https://example.com")
>>> qrcode.save(qr, "example.png")
nil
>>> qrcode.save(qr, "example-large.png", 200)
nil
```

### to_base64

```go filename="Function signature"
to_base64(qr, width=256) string
```

Converts the QR code to a base64-encoded PNG image string. An optional width parameter can be specified
(defaults to 256 pixels). This is useful for embedding QR codes in HTML or other formats.

```go filename="Example"
>>> qr := qrcode.create("https://example.com")
>>> base64_str := qrcode.to_base64(qr)
>>> base64_str[0:20]  // Show first 20 characters
"iVBORw0KGgoAAAANSU"
```

## Types

### qrcode

The qrcode object represents a generated QR code that can be saved to a file or converted to other formats.

The qrcode object has the following methods:

#### save

```go filename="Method signature"
qrcode.save(path) nil
```

Saves the QR code to a PNG file at the specified path.

```go filename="Example"
>>> qr := qrcode.create("https://example.com")
>>> qr.save("example.png")
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

In addition to these methods, you can also use the module functions `save` and `to_base64` to work with the generated QR code.

## Complete Example

Here's a complete example showing how to create a QR code and save it or convert it to base64:

```go
// Create a QR code for a website
qr := qrcode.create("https://example.com", {
    "encoding_mode": "byte",
    "error_correction": "high"
})

// Save to a file using the object method
qr.save("example.png")

// Or using the module function
qrcode.save(qr, "example2.png")

// Get the dimension
dim := qr.dimension()
print("QR code dimension:", dim)

// Get base64 representation
base64_data := qrcode.to_base64(qr, 200)
print("data:image/png;base64," + base64_data)
```
