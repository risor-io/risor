# gha

Module `gha` provides utility functions when working inside GitHub Actions.

## Functions

### is_debug

```go filename="Function signature"
is_debug() bool
```

Gets whether Actions Step Debug is on or not
(by checking the `RUNNER_DEBUG` environment variable)

```go copy filename="Example"
>>> gha.is_debug()
false
```

### log_debug

```go filename="Function signature"
log_debug(msg string) bool
```

Writes debug message to log. This message is only shown when Action Step Debug
is on. See also: [`gha.is_debug()`](#is_debug)

```go copy filename="Example"
>>> gha.log_debug("Hello world")
::debug::Hello world
```

### log_notice

```go filename="Function signature"
log_notice(msg string, props={})
```

Adds a notice issue.

The [`props` parameter](#annotation-properties) can specify thing like a file
path and line number to add the issue to.

```go copy filename="Example"
>>> gha.log_notice("Hello world")
::notice::Hello world
>>> gha.log_notice("Hello world", {title: "Risor", file: "somefile.txt", line: 5})
::notice file=somefile.txt,title=Risor,line=5::Hello world
```

### log_warning

```go filename="Function signature"
log_warning(msg string, props={})
```

Adds a warning issue.

The [`props` parameter](#annotation-properties) can specify thing like a file
path and line number to add the issue to.

```go copy filename="Example"
>>> gha.log_warning("Hello world")
::warning::Hello world
>>> gha.log_warning("Hello world", {title: "Risor", file: "somefile.txt", line: 5})
::warning file=somefile.txt,title=Risor,line=5::Hello world
```

### log_error

```go filename="Function signature"
log_error(msg string, props={})
```

Adds a error issue.

The [`props` parameter](#annotation-properties) can specify thing like a file
path and line number to add the issue to.

```go copy filename="Example"
>>> gha.log_error("Hello world")
::error::Hello world
>>> gha.log_error("Hello world", {title: "Risor", file: "somefile.txt", line: 5})
::error file=somefile.txt,title=Risor,line=5::Hello world
```

### start_group

```go filename="Function signature"
start_group(name string)
```

Begins an output group. Output until the next `end_group` will be foldable
in this group.

```go copy filename="Example"
>>> gha.start_group("My group")
::group::My group
```

### end_group

```go filename="Function signature"
end_group()
```

End an output group.

```go copy filename="Example"
>>> gha.end_group()
::endgroup::
```

### set_output

```go filename="Function signature"
set_output(name string, value any)
```

Sets a GitHub Action output variable.

This function makes use of the `GITHUB_OUTPUT` environment variable, if set.
Otherwise it falls back to the ([deprecated](https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/))
workflow command of `::set-output::`.

```go copy filename="Example"
>>> gha.set_output("my-var", "some value")
::set-output name=my-var::some value
```

### set_env

```go filename="Function signature"
set_env(name string, value any)
```

Sets a GitHub Action environment variable for this action and future actions
in the same job.

This function makes use of the `GITHUB_ENV` environment variable, if set.
Otherwise it falls back to the ([deprecated](https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/))
workflow command of `::set-env::`.

```go copy filename="Example"
>>> gha.set_output("MY_VAR", "some value")
::set-env name=MY_VAR::some value
>>> os.getenv("MY_VAR")
"some value"
```

### add_path

```go filename="Function signature"
add_path(dir string)
```

Prepends directory to the PATH (for this action and future actions)

This function makes use of the `GITHUB_PATH` environment variable, if set.
Otherwise it falls back to the ([deprecated](https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/))
workflow command of `::add-path::`.

```go copy filename="Example"
>>> gha.add_path("/some/new/dir")
::add-path::/some/new/dir
>>> os.getenv("PATH")
"/some/new/dir:/usr/local/bin:/usr/bin"
```

## Types

### Annotation properties

The following keys can be supplied when creating an issue with the
[`log_notice`](#log_notice), [`log_warning`](#log_warning),
or [`log_error`](#log_error) functions.
All keys are optional.

| Key        | Type   | Description                                                                                                                                        |
| ---------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| title      | string | A title for the annotation.                                                                                                                        |
| file       | string | The path of the file for which the annotation should be created (relative to the repository root directory)                                        |
| line       | int    | The start line for the annotation.                                                                                                                 |
| column     | int    | The start column for the annotation. Cannot be sent when `line` and `end_line` are different values.                                               |
| end_line   | int    | The end line for the annotation. Defaults to `line` when `line` is provided.                                                                       |
| end_column | int    | The end column for the annotation. Cannot be sent when `line` and `end_line` are different values. Defaults to `column` when `column` is provided. |
