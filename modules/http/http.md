import { Callout, Steps } from 'nextra/components';

# http

The `http` module provides functions for making HTTP requests and defines
[request](#request-1) and [response](#response) types.

<Callout type="info" emoji="ℹ️">
    Using the [fetch](/docs/builtins#fetch) built-in function instead of this
    module is encouraged. This module provides more low-level control of how
    HTTP requests are built and sent, which may be useful in some situations.
</Callout>

The approach used to make an HTTP request is as follows:

<Steps>
### Create a Request

```go copy
req := http.get("https://api.ipify.org")
```

### Send the Request

```go copy
res := req.send()
```

### Handle the Response

```go copy
print("response status:", res.status, "text:", res.text())
```

</Steps>

## Functions

### get

```go filename="Function signature"
get(url string, headers map, params map) request
```

Creates a new GET request with the given URL, headers, and query parameters. The
headers and query parameters are optional. The request that is returned can then
be executed using its `send` method. Read more about the [request](#request)
type below.

```go copy filename="Example"
>>> http.get("https://api.ipify.org").send()
http.response(status: "200 OK", content_length: 14)
```

### delete

```go filename="Function signature"
delete(url string, headers map, params map) request
```

Creates a new DELETE request with the given URL, headers, and query parameters.
The headers and query parameters are optional.

### head

```go filename="Function signature"
head(url string, headers map, params map) request
```

Creates a new HEAD request with the given URL, headers, and query parameters.
The headers and query parameters are optional.

### listen_and_serve

```go filename="Function signature"
listen_and_serve(addr string, handler func(w response_writer, r request))
```

Starts an HTTP server that listens on the specified address and calls the
handler function to handle requests. As a convenience, the handler function
may return a map or list object to be marshaled as JSON, or a string or byte
slice object which will be written as the response body as-is.

### listen_and_serve_tls

```go filename="Function signature"
listen_and_serve_tls(addr, cert_file, key_file string, handler func(w response_writer, r request))
```

Acts the same as `listen_and_serve`, but uses the provided certificate and key
files to work over HTTPS.

### patch

```go filename="Function signature"
patch(url string, headers map, body byte_slice) request
```

Creates a new PATCH request with the given URL, headers, and request body.
The headers and request body parameters are optional.

### post

```go filename="Function signature"
post(url string, headers map, body byte_slice) request
```

Creates a new POST request with the given URL, headers, and request body.
The headers and request body parameters are optional.

### put

```go filename="Function signature"
put(url string, headers map, body byte_slice) request
```

Creates a new PUT request with the given URL, headers, and request body.
The headers and request body parameters are optional.

### request

```go filename="Function signature"
request(url string, options map) request
```

Creates a new request with the given URL and options. If provided, the options
map may contain any of the following keys:

| Name    | Type                 | Description                                            |
| ------- | -------------------- | ------------------------------------------------------ |
| method  | string               | The HTTP method to use.                                |
| headers | map                  | The headers to send with the request.                  |
| params  | map                  | The query parameters to send with the request.         |
| body    | byte_slice or reader | The request body.                                      |
| timeout | int                  | Request timeout in milliseconds.                       |
| data    | object               | Object to marshal as JSON and send in the request body |

If both `body` and `data` are provided, the `body` value will be used.

## Types

### request

Represents an HTTP request that is being built and sent.

#### Attributes

| Name           | Type                           | Description                                  |
| -------------- | ------------------------------ | -------------------------------------------- |
| url            | string                         | The URL of the request.                      |
| content_length | int                            | The length of the request body.              |
| header         | map                            | The headers of the request.                  |
| send           | func()                         | Sends the request.                           |
| add_header     | func(key string, value object) | Adds a header to the request.                |
| add_cookie     | func(key string, value map)    | Adds a cookie to the request.                |
| set_body       | func(body byte_slice)          | Sets the request body.                       |
| set_data       | func(data object)              | Request body as an object to marshal to JSON |

### response

Represents an HTTP response.

#### Attributes

| Name           | Type          | Description                      |
| -------------- | ------------- | -------------------------------- |
| status         | string        | The status of the response.      |
| status_code    | int           | The status code of the response. |
| proto          | string        | The protocol of the response.    |
| content_length | int           | The length of the response body. |
| header         | map           | The headers of the response.     |
| cookies        | map           | The cookies of the response.     |
| response       | object        | The response body.               |
| json           | func() object | The response body as JSON.       |
| text           | func() string | The response body as text.       |
| close          | func()        | Closes the response body.        |

### response_writer

Represents an HTTP response writer.

#### Attributes

| Name         | Type                    | Description                                                  |
| ------------ | ----------------------- | ------------------------------------------------------------ |
| add_header   | func(key, value string) | Adds a header to the header map that will be sent.           |
| del_header   | func(key string)        | Deletes a header from the header map that will be sent.      |
| write        | func(object)            | Writes the object as the HTTP reply.                         |
| write_header | func(status_code int)   | Sends an HTTP response header with the provided status code. |
