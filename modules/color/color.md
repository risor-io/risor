# color

The `color` module provides functions to colorize text in the
terminal.

The core functionality is provided by
[github.com/fatih/color](https://github.com/fatih/color).

Configuring colors and other text attributes is done using the various
[constants](#constants) described below. Multiple constants can be passed
to functions in this module to combine attributes.

## Module

```go copy filename="Function signature"
color(options ...int) color.color
```

The `color` module object itself is callable, which is a shorthand for
`color.color()`.

```go copy filename="Example"
>>> color(color.fg_green).printf("hello!\n")
// colorized output
```

## Functions

### color

```go filename="Function signature"
color(options ...int) color.color
```

Creates a new color object with the specified options.

```go copy filename="Example"
>>> c := color.color(color.fg_red, color.bold)
```

### set

```go filename="Function signature"
set(options ...int)
```

Sets the terminal color to the specified options.

```go copy filename="Example"
>>> color.set(color.fg_green, color.bold)
```

### unset

```go filename="Function signature"
unset()
```

Resets the terminal color to the default.

```go copy filename="Example"
>>> color.unset()
```

## Types

### color

Defines a custom color object.

#### Methods

##### sprintf

```go filename="Method signature"
sprintf(format string, a ...object) string
```

Formats the string with the color object.

```go copy filename="Example"
>>> c := color.color(color.fg_red, color.bold)
>>> c.sprintf("Hello, %s!", "world")
"\x1b[31;1mHello, world!\x1b[0;22m"
>>> print(c.sprintf("Hello, %s!", "world"))
// colorized output
```

##### fprintf

```go filename="Method signature"
fprintf(w io.writer, format string, a ...object)
```

Writes the colorized, formatted string to the given writer.

```go copy filename="Example"
>>> b := buffer()
>>> c.fprintf(b, "Hello, %s!\n", "world")
>>> b
buffer("\x1b[31;1mHello, world!\n\x1b[0m")
```

##### printf

```go filename="Method signature"
printf(format string, a ...object)
```

Writes the colorized, formatted string to the standard output. Note
you may need to include a trailing newline to flush the output.

```go copy filename="Example"
>>> c.printf("Hello, %s!\n", "world")
// colorized output
```

## Constants

| Name         | Description                             |
| ------------ | --------------------------------------- |
| reset        | Reset all attributes                    |
| bold         | Bold text                               |
| dim          | Dim text                                |
| italic       | Italic text                             |
| underline    | Underlined text                         |
| blinkslow    | Blinking text (slow)                    |
| blinkrapid   | Blinking text (rapid)                   |
| reversevideo | Reverse video                           |
| concealed    | Concealed text                          |
| crossedout   | Crossed-out text                        |
| bg_black     | Black background color                  |
| bg_blue      | Blue background color                   |
| bg_cyan      | Cyan background color                   |
| bg_green     | Green background color                  |
| bg_hiblack   | High-intensity black background color   |
| bg_hiblue    | High-intensity blue background color    |
| bg_hicyan    | High-intensity cyan background color    |
| bg_higreen   | High-intensity green background color   |
| bg_himagenta | High-intensity magenta background color |
| bg_hired     | High-intensity red background color     |
| bg_hiwhite   | High-intensity white background color   |
| bg_hiyellow  | High-intensity yellow background color  |
| bg_magenta   | Magenta background color                |
| bg_red       | Red background color                    |
| bg_white     | White background color                  |
| bg_yellow    | Yellow background color                 |
| fg_black     | Black foreground color                  |
| fg_blue      | Blue foreground color                   |
| fg_cyan      | Cyan foreground color                   |
| fg_green     | Green foreground color                  |
| fg_hiblack   | High-intensity black foreground color   |
| fg_hiblue    | High-intensity blue foreground color    |
| fg_hicyan    | High-intensity cyan foreground color    |
| fg_higreen   | High-intensity green foreground color   |
| fg_himagenta | High-intensity magenta foreground color |
| fg_hired     | High-intensity red foreground color     |
| fg_hiwhite   | High-intensity white foreground color   |
| fg_hiyellow  | High-intensity yellow foreground color  |
| fg_magenta   | Magenta foreground color                |
| fg_red       | Red foreground color                    |
| fg_white     | White foreground color                  |
| fg_yellow    | Yellow foreground color                 |
