# os

Module `os` provides a platform-independent interface to operating system
functionality.

By default, this module interacts with the host operating system normally by
calling the underlying Go `os` package. However, alternative OS abstraction
layers may be used via the Go [WithOS](https://pkg.go.dev/github.com/risor-io/risor@v1.2.0/os#WithOS) function. This assists with sandboxing scripts and providing
access to object storage like AWS S3 via a filesystem-like interface.

## Attributes

### stdin

`stdin` is an open file pointing to the standard input for the process.

```go copy filename="Example"
>>> os.stdin.read()
byte_slice("hello world")
```

### stdout

`stdout` is an open file pointing to the standard output for the process.

```go copy filename="Example"
>>> os.stdout.write("hello world")
11
```

### err_not_exist

`err_not_exist` is an error indicating that a file or directory does not exist.

```go copy filename="Example"
>>> if errors.is(err, os.err_not_exist) { print("file does not exist") }
```

### err_exist

`err_exist` is an error indicating that a file or directory already exists.

```go copy filename="Example"
>>> if errors.is(err, os.err_exist) { print("file already exists") }
```

### err_permission

`err_permission` is an error indicating that permission is denied.

```go copy filename="Example"
>>> if errors.is(err, os.err_permission) { print("permission denied") }
```

### err_closed

`err_closed` is an error indicating that the file is already closed.

```go copy filename="Example"
>>> if errors.is(err, os.err_closed) { print("file already closed") }
```

### err_invalid

`err_invalid` is an error indicating that the operation is invalid.

```go copy filename="Example"
>>> if errors.is(err, os.err_invalid) { print("invalid operation") }
```

### err_no_deadline

`err_no_deadline` is an error indicating that no deadline is set.

```go copy filename="Example"
>>> if errors.is(err, os.err_no_deadline) { print("no deadline set") }
```

### err_deadline_exceeded

`err_deadline_exceeded` is an error indicating that the deadline has been exceeded.

```go copy filename="Example"
>>> if errors.is(err, os.err_deadline_exceeded) { print("deadline exceeded") }
```

## Functions

### chdir

```go filename="Function signature"
chdir(dir string)
```

Changes the working directory to dir.

```go copy filename="Example"
>>> os.chdir("/tmp")
>>> os.getwd()
"/tmp"
```

### create

```go filename="Function signature"
create(name string) File
```

Creates or truncates the named file.

```go copy filename="Example"
>>> f := os.create("foo.txt")
>>> f.write("hello world")
11
>>> f.close()
```

### environ

```go filename="Function signature"
environ() list
```

Returns a copy of strings representing the environment, in the form "key=value".

```go copy filename="Example"
>>> os.environ()
["TERM=xterm-256color", "SHELL=/bin/bash", "USER=alice", ...]
```

### exit

```go filename="Function signature"
exit(code int)
```

Terminates the program with the given exit code.

```go copy filename="Example"
>>> os.exit(0)
```

### getenv

```go filename="Function signature"
getenv(key string) string
```

Returns the value of the environment variable key.

```go copy filename="Example"
>>> os.getenv("USER")
"alice"
```

### getpid

```go filename="Function signature"
getpid() int
```

Returns the current process ID.

```go copy filename="Example"
>>> os.getpid()
1234
```

### getuid

```go filename="Function signature"
getuid() int
```

Returns the current user ID.

```go copy filename="Example"
>>> os.getuid()
501
```

### getwd

```go filename="Function signature"
getwd() string
```

Returns the current working directory.

```go copy filename="Example"
>>> os.getwd()
"/home/alice"
```

### hostname

```go filename="Function signature"
hostname() string
```

Returns the host name reported by the kernel.

```go copy filename="Example"
>>> os.hostname()
"alice-macbook-pro-1.local"
```

### mkdir_all

```go filename="Function signature"
mkdir_all(path string, perm int)
```

Creates a directory named path, along with any necessary parent directories.

```go copy filename="Example"
>>> os.mkdir_all("/tmp/foo/bar", 0755)
```

### mkdir_temp

```go filename="Function signature"
mkdir_temp(dir, prefix string) string
```

Creates a new temporary directory in the directory dir, using prefix to generate
its name.

```go copy filename="Example"
>>> os.mkdir_temp("/tmp", "foo")
"/tmp/foo4103914411"
```

### mkdir

```go filename="Function signature"
mkdir(path string, perm int)
```

Creates a new directory with the specified name and permission bits. If
a permissions value is not specified, 0755 is used.

```go copy filename="Example"
>>> os.mkdir("/tmp/foo", 0755)
```

### open

```go filename="Function signature"
open(name string) File
```

Opens the named file.

```go copy filename="Example"
>>> f := os.open("foo.txt")
>>> f.read()
byte_slice("hello world")
>>> f.close()
```

### read_dir

```go filename="Function signature"
read_dir(name string) list
```

Returns a list of directory entries sorted by filename. If a name is not
specified, the current directory is used.

```go copy filename="Example"
>>> os.read_dir("/tmp")
[dir_entry(name=foo.txt, type=regular), dir_entry(name=bar.txt, type=regular)]
```

### read_file

```go filename="Function signature"
read_file(name string) byte_slice
```

Reads the named file and returns its contents.

```go copy filename="Example"
>>> os.read_file("/tmp/foo.txt")
byte_slice("hello world")
```

### remove

```go filename="Function signature"
remove(name string)
```

Removes the named file or empty directory.

```go copy filename="Example"
>>> os.remove("/tmp/old/junk.txt")
```

### remove_all

```go filename="Function signature"
remove_all(name string)
```

Removes path and any children it contains.

```go copy filename="Example"
>>> os.remove_all("/tmp/junk")
```

### rename

```go filename="Function signature"
rename(old, new string)
```

Renames (moves) old to new.

```go copy filename="Example"
>>> os.rename("old.txt", "new.txt")
```

### setenv

```go filename="Function signature"
setenv(key, value string)
```

Sets the value of the environment variable key to value.

```go copy filename="Example"
>>> os.setenv("USER", "bob")
>>> os.getenv("USER")
"bob"
```

### stat

```go filename="Function signature"
stat(name string) FileInfo
```

Returns a FileInfo describing the named file.

```go copy filename="Example"
>>> os.stat("example.txt")
file_info(name=example.txt, mode=-rw-r--r--, size=84, mod_time=2023-08-06T08:45:56-04:00)
```

### symlink

```go filename="Function signature"
symlink(old, new string)
```

Creates a symbolic link new pointing to old.

```go copy filename="Example"
>>> os.symlink("foo.txt", "bar.txt")
```

### temp_dir

```go filename="Function signature"
temp_dir() string
```

Returns the default directory to use for temporary files.

```go copy filename="Example"
>>> os.temp_dir()
"/tmp"
```

### unsetenv

```go filename="Function signature"
unsetenv(key string)
```

Unsets the environment variable key.

```go copy filename="Example"
>>> os.unsetenv("USER")
>>> os.getenv("USER")
""
```

### user_cache_dir

```go filename="Function signature"
user_cache_dir() string
```

Returns the default root directory to use for user-specific non-essential data.

```go copy filename="Example"
>>> os.user_cache_dir()
"/home/alice/.cache"
```

### user_config_dir

```go filename="Function signature"
user_config_dir() string
```

Returns the default root directory to use for user-specific configuration data.

```go copy filename="Example"
>>> os.user_config_dir()
"/home/alice/.config"
```

### user_home_dir

```go filename="Function signature"
user_home_dir() string
```

Returns the current user's home directory.

```go copy filename="Example"
>>> os.user_home_dir()
"/home/alice"
```

### write_file

```go filename="Function signature"
write_file(name string, data byte_slice / string)
```

Writes the given byte_slice or string to the named file.

```go copy filename="Example"
>>> os.write_file("example.txt", "hey!")
>>> os.read_file("example.txt")
byte_slice("hey!")
```

### current_user

```go filename="Function signature"
current_user() map
```

Returns a map representing the current user.

```go copy filename="Example"
>>> os.current_user()
{"gid": "20", "home_dir": "/Users/alice", "name": "Alice", "uid": "501", "username": "alice"}
```

### lookup_user

```go filename="Function signature"
lookup_user(name string) map
```

Looks up a user by name and returns a map representation.

```go copy filename="Example"
>>> os.lookup_user("bob")
{"gid": "20", "home_dir": "/Users/bob", "name": "Bob", "uid": "502", "username": "bob"}
```

### lookup_uid

```go filename="Function signature"
lookup_uid(uid string) map
```

Looks up a user by user ID and returns a map representation.

```go copy filename="Example"
>>> os.lookup_uid("501")
{"gid": "20", "home_dir": "/Users/alice", "name": "Alice", "uid": "501", "username": "alice"}
```

### lookup_group

```go filename="Function signature"
lookup_group(name string) map
```

Looks up a group by name and returns a map representation.

```go copy filename="Example"
>>> os.lookup_group("staff")
{"gid": "20", "name": "staff"}
```

### lookup_gid

```go filename="Function signature"
lookup_gid(gid string) map
```

Looks up a group by group ID and returns a map representation.

```go copy filename="Example"
>>> os.lookup_gid("20")
{"gid": "20", "name": "staff"}
```
