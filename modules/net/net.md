# net

The `net` module provides a set of functions for working with network addresses
and performing network lookups.

This module does not yet provide a way to create network connections or servers.
Feel free to open a GitHub issue if you would like to see this functionality
added.

The core functionality is provided by the Go standard library's
[`net`](https://pkg.go.dev/net) package.

## Functions

### interface_addrs

```go filename="Function signature"
interface_addrs() []string
```

Returns a list of the system's network interfaces.

```go filename="Example"
>>> net.interface_addrs()
["127.0.0.1/8", "::1/128", "fe80::1/64", "etc."]
```

### join_host_port

```go filename="Function signature"
join_host_port(host, port string) string
```

Joins the host and port together.

```go filename="Example"
>>> net.join_host_port("localhost", "8080")
"localhost:8080"
```

### lookup_addr

```go filename="Function signature"
lookup_addr(addr string) (string, error)
```

Looks up the host name of the specified address.

```go filename="Example"
>>> net.lookup_addr("127.0.0.1")
["localhost"]
```

### lookup_host

```go filename="Function signature"
lookup_host(host string) []string
```

Looks up the IP addresses of the specified host.

```go filename="Example"
>>> net.lookup_host("google.com")
["172.253.62.113", "172.253.62.102", "172.253.62.139", "etc."]
```

### lookup_ip

```go filename="Function signature"
lookup_ip(host string) []net.ip
```

Looks up the IP addresses of the specified host.

```go filename="Example"
>>> net.lookup_ip("google.com")
[net.ip(172.253.62.113),
 net.ip(172.253.62.102),
 net.ip(172.253.62.139),
 etc.]
```
