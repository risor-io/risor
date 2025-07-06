# ssh

The `ssh` module provides functions for establishing SSH connections and executing commands remotely on servers.

## Module

```go copy filename="Function signature"
ssh.connect(user string, host string, options map) ssh.client
```

Establishes an SSH connection to a remote host. The module is callable, so `ssh(...)` is equivalent to `ssh.connect(...)`.

**Parameters:**
- `user` (string): Username for authentication
- `host` (string): Hostname or IP address to connect to
- `options` (map): Configuration options

**Options:**
- `password` (string): Password for authentication (optional)
- `private_key` (string): Private key content for key-based authentication (optional)
- `port` (int): Port number (optional, defaults to 22)
- `timeout` (int): Connection timeout in seconds (optional, defaults to 30)
- `insecure` (bool): Use insecure host key callback (optional, defaults to false)

**Note:** You must provide either `password` or `private_key` in the options map for authentication.

```go copy filename="Example"
>>> client := ssh.connect("user", "example.com", {
...   password: "secret123",
...   port: 22,
...   timeout: 10
... })
>>> client.execute("ls -la")
"total 24\ndrwxr-xr-x 3 user user 4096 Jan 15 10:30 .\ndrwxr-xr-x 5 user user 4096 Jan 15 10:29 ..\n-rw-r--r-- 1 user user  220 Jan 15 10:29 .bash_logout\n"

>>> // Shorter form using callable module
>>> client := ssh("user", "example.com", {
...   password: "secret123"
... })
>>> client.execute("whoami")
"user"
```

## Client

The SSH client provides methods for executing commands on remote servers.

### execute

```go filename="Method signature"
execute(command string) string
```

Executes a command on the SSH connection. This method automatically creates a session, runs the command, and closes the session.

```go filename="Example"
>>> client := ssh("user", "example.com", {
...   password: "secret123"
... })
>>> output := client.execute("ls -la")
>>> print(output)
total 24
drwxr-xr-x 3 user user 4096 Jan 15 10:30 .
drwxr-xr-x 5 user user 4096 Jan 15 10:29 ..
-rw-r--r-- 1 user user  220 Jan 15 10:29 .bash_logout

>>> uptime := client.execute("uptime")
>>> print("Server uptime:", uptime)
Server uptime:  10:30:42 up 1 day,  2:45,  1 user,  load average: 0.08, 0.02, 0.01
```

### close

```go filename="Method signature"
close()
```

Closes the SSH connection and frees associated resources.

```go filename="Example"
>>> client := ssh("user", "example.com", {
...   password: "secret123"
... })
>>> client.execute("whoami")
"user"
>>> client.close()
```

## Authentication Examples

### Password Authentication

```go copy filename="Example"
>>> client := ssh("ubuntu", "192.168.1.100", {
...   password: "mypassword",
...   port: 22,
...   timeout: 30
... })
>>> client.execute("uname -a")
"Linux server 5.4.0-88-generic #99-Ubuntu SMP Thu Sep 23 17:29:00 UTC 2021 x86_64 x86_64 x86_64 GNU/Linux"
```

### Private Key Authentication

```go copy filename="Example"
>>> private_key := `-----BEGIN RSA PRIVATE KEY-----
... MIIEpAIBAAKCAQEA...
... -----END RSA PRIVATE KEY-----`
>>> client := ssh("user", "example.com", {
...   private_key: private_key,
...   port: 2222
... })
>>> client.execute("hostname")
"example.com"
```

### Insecure Host Key

```go copy filename="Example"
>>> // Not recommended for production use
>>> client := ssh("user", "example.com", {
...   password: "secret123",
...   insecure: true
... })
>>> client.execute("echo 'Connected with insecure host key verification'")
"Connected with insecure host key verification"
```

## Types

### ssh.client

The SSH client represents an active SSH connection to a remote host and provides
methods for executing commands.

## Security Notes

- By default, the module uses the system's `~/.ssh/known_hosts` file for host key verification
- Set `insecure: true` in the options to disable host key verification (not recommended for production)
- Always close connections when finished to free resources
- Private keys should be stored securely and not hardcoded in scripts
- Each call to `execute()` creates a new session that is automatically closed after the command completes
