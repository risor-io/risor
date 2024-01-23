import { Callout } from 'nextra/components';

# vault

<Callout type="info" emoji="ℹ️">
  This module requires that Risor has been compiled with the `vault` Go build tag.
  When compiling **manually**, [make sure you specify `-tags vault`](https://github.com/risor-io/risor#build-and-install-the-cli-from-source).
</Callout>

Module `vault` provides a client to interact with Hashicorp Vault

## Functions

### connect

```go filename="Function signature"
connect(address string)
```

Instanciates a new client for the given Vault address

```go filename="Example"
client := vault.connect("http://127.0.0.1:8200")
client.token = "hvs.AJ71UUKBsv1jiW7pJljTz4BN"
```

### write

```go filename="Function signature"
write(data object, path string)
```

Writes the given object to the given path.

```go filename="Example"
client.write({"data": {"password1": "t0p$ecret", "password2": "cl@$$ified"}}, "/secret/data/foo")
```

### write_raw

```go filename="Function signature"
write_raw(data byte_slice, path string)
```

Writes the given data to the given path.

```go filename="Example"
client.write_raw(byte_slice('{"data": {"password1": "t0p$ecret", "password2": "cl@$$ified"}}'), "/secret/data/foo")
```

### read

```go filename="Function signature"
read(path string) object
```

Returns an object from the given path

```go filename="Example"
client.read("/secret/data/foo")
```

### read_raw

```go filename="Function signature"
read_raw(path string) object
```

Returns an HTTP response object from Vault from the given path

```go filename="Example"
client.read("/secret/data/foo")
```

### delete

```go filename="Function signature"
delete(path string) object
```

Deletes an object from the given path

```go filename="Example"
client.delete("/secret/data/foo")
```

### list

```go filename="Function signature"
list(path string) object
```

Returns a list objects from the given path

```go filename="Example"
client.list("/secret/data/foo")
```
