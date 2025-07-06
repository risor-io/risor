# ssh

The SSH module provides functionality to establish SSH connections and execute commands on remote servers.

## Functions

### connect

```go filename="Function signature"
connect(host string, user string, password string, timeout int, port int) ssh_client
```

Establishes an SSH connection using password authentication.

**Parameters:**
- `host` (string): The hostname or IP address of the SSH server
- `user` (string): The username for authentication
- `password` (string): The password for authentication
- `timeout` (int, optional): Connection timeout in seconds (default: 30)
- `port` (int, optional): SSH port (default: 22)

**Returns:** SSH client object for use with other SSH functions

```go copy filename="Example"
client := ssh.connect("example.com", "myuser", "mypassword")
```

### connect_with_key

```go filename="Function signature"
connect_with_key(host string, user string, private_key string, timeout int, port int) ssh_client
```

Establishes an SSH connection using private key authentication.

**Parameters:**
- `host` (string): The hostname or IP address of the SSH server
- `user` (string): The username for authentication
- `private_key` (string): The private key content (PEM format)
- `timeout` (int, optional): Connection timeout in seconds (default: 30)
- `port` (int, optional): SSH port (default: 22)

**Returns:** SSH client object for use with other SSH functions

```go copy filename="Example"
key := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----`
client := ssh.connect_with_key("example.com", "myuser", key)
```

### execute

```go filename="Function signature"
execute(client ssh_client, command string) string
```

Executes a command on the SSH connection and returns the output.

**Parameters:**
- `client` (ssh_client): SSH client object from `connect()` or `connect_with_key()`
- `command` (string): The command to execute

**Returns:** Command output as a string

```go copy filename="Example"
client := ssh.connect("example.com", "myuser", "mypassword")
output := ssh.execute(client, "ls -la")
print(output)
ssh.close(client)
```

### new_session

```go filename="Function signature"
new_session(client ssh_client) ssh_session
```

Creates a new SSH session for executing multiple commands.

**Parameters:**
- `client` (ssh_client): SSH client object from `connect()` or `connect_with_key()`

**Returns:** SSH session object

```go copy filename="Example"
client := ssh.connect("example.com", "myuser", "mypassword")
session := ssh.new_session(client)
```

### session_run

```go filename="Function signature"
session_run(session ssh_session, command string) string
```

Executes a command on an SSH session.

**Parameters:**
- `session` (ssh_session): SSH session object from `new_session()`
- `command` (string): The command to execute

**Returns:** Command output as a string

```go copy filename="Example"
client := ssh.connect("example.com", "myuser", "mypassword")
session := ssh.new_session(client)
output := ssh.session_run(session, "whoami")
print(output)
ssh.session_close(session)
ssh.close(client)
```

### session_close

```go filename="Function signature"
session_close(session ssh_session) null
```

Closes an SSH session.

**Parameters:**
- `session` (ssh_session): SSH session object to close

**Returns:** null

```go copy filename="Example"
session := ssh.new_session(client)
// ... use session ...
ssh.session_close(session)
```

### close

```go filename="Function signature"
close(client ssh_client) null
```

Closes an SSH connection.

**Parameters:**
- `client` (ssh_client): SSH client object to close

**Returns:** null

```go copy filename="Example"
client := ssh.connect("example.com", "myuser", "mypassword")
// ... use client ...
ssh.close(client)
```

## Usage Examples

### Basic Password Authentication

```go copy filename="Example"
// Connect with password
client := ssh.connect("example.com", "user", "password", 30, 22)

// Execute a command
output := ssh.execute(client, "uptime")
print("Server uptime:", output)

// Close the connection
ssh.close(client)
```

### Key-Based Authentication

```go copy filename="Example"
// Load private key
private_key := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----`

// Connect with key
client := ssh.connect_with_key("example.com", "user", private_key)

// Execute commands
output := ssh.execute(client, "df -h")
print("Disk usage:", output)

ssh.close(client)
```

### Session Management

```go copy filename="Example"
// Connect and create session
client := ssh.connect("example.com", "user", "password")
session := ssh.new_session(client)

// Execute multiple commands on the same session
output1 := ssh.session_run(session, "pwd")
output2 := ssh.session_run(session, "ls -la")

print("Current directory:", output1)
print("Directory listing:", output2)

// Clean up
ssh.session_close(session)
ssh.close(client)
```

## Security Notes

- The current implementation uses `ssh.InsecureIgnoreHostKey()` for host key verification. In production environments, you should implement proper host key verification.
- Always close SSH connections and sessions when done to free resources.
- Store private keys securely and never include them directly in your code in production environments.