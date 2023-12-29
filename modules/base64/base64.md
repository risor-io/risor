# base64

## Functions

### decode

```go filename="Function signature"
decode(s string, pad bool) byte_slice
```

Decode base64 string s to a byte_slice, with padding if pad is true.
If not provided, pad defaults to false.

```go copy filename="Example"
>>> base64.decode("aGVsbG8=")
byte_slice("hello")
```

### encode

```go filename="Function signature"
encode(b byte_slice, pad bool) string
```

Encode byte_slice b to a base64 string. The encoded string is padded
if pad is true. If not provided, pad defaults to false.

```go copy filename="Example"
>>> base64.encode("hello")
"aGVsbG8="
```

### url_decode

```go filename="Function signature"
url_decode(s string, pad bool) byte_slice
```

Decode base64 string s to a byte slice using the alternate base64 codec.
The string is understood to be padded if pad is true. If not provided, pad
defaults to false.

```go copy filename="Example"
>>> base64.url_decode("YWJjK2Zvbz9iYXI9YmF6")
byte_slice("abc+foo?bar=baz")
```

### url_encode

```go filename="Function signature"
url_encode(b byte_slice, pad bool) string
```

Encode byte slice b to a base64 string using the alternate base64 codec.
The encoded string is padded if pad is true. If not provided, pad defaults
to false. The encoded string is safe for use in URLs and file names.

```go
>>> base64.url_encode("abc+foo?bar=baz")
"YWJjK2Zvbz9iYXI9YmF6"
```
