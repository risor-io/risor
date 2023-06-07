# Executable Scripts

Tamarin scripts can easily be made directly executable on MacOS and Linux.
In the steps below, the Tamarin script filename is `myscript`. The full path
to the script is `/path/to/myscript`.

1\. Add a shebang line at the start of the Tamarin script.

```
#!/usr/bin/env tamarin
```

2\. Allow execution as a program by running `chmod` against the script:

```
chmod +x /path/to/myscript
```

3\. Optionally, update your `PATH` variable so that your shell can find the script:

```
export PATH=/path/to/:$PATH
```

Having done that, you should be able to run `myscript` as a program from your shell.
You should add the `export PATH` statement to `~/.bashrc` or `~/.zshrc` to persist the
modified `PATH` variable for future sessions.

## Example Script

```go
#!/usr/bin/env tamarin

print("just a test")
```
