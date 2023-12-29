# image

Module `image` provides 2-D image encoding and decoding. Image pixels are
represented as RGBA values.

Supported image formats: `bmp`, `jpg`, `png`.

## Functions

### decode

```go filename="Function signature"
decode(b byte_slice) image
```

Returns an image object that is decoded from the given bytes. If a byte_buffer
or io.Reader is given, it is automatically converted to a byte_slice.

```go copy filename="Example"
>>> img := image.decode(open("/path/to/test.png"))
>>> img.width
3440
>>> img.height
1416
>>> img.dimensions()
{"height": 1416, "width": 3440}
>>> img.bounds()
{"max": {"x": 3440, "y": 1416}, "min": {"x": 0, "y": 0}}
>>> img.at(0, 0)
color(r=52428 g=31611 b=7967 a=65535)
```

### encode

```go filename="Function signature"
encode(img image, format string) byte_slice
```

Encodes the given image object into the given format, returning the encoded bytes.

```go copy filename="Example"
>>> img
image(width=256, height=256)
>>> image.encode(img, "png")
byte_slice(...)
```

## Types

### image

The `image` type represents a 2-D image as a rectangular grid of color values.

#### Attributes

| Name       | Type           | Description                                 |
| ---------- | -------------- | ------------------------------------------- |
| width      | int            | The width of the image in pixels            |
| height     | int            | The height of the image in pixels           |
| dimensions | func() map     | The width and height of the image in pixels |
| bounds     | func() map     | The bounds of the image                     |
| at         | func(x, y int) | Returns the color at the given coordinates  |

### color

The `color` type represents a color as an RGBA value.

#### Attributes

| Name | Type | Description                 |
| ---- | ---- | --------------------------- |
| rgba | list | The RGBA value of the color |
