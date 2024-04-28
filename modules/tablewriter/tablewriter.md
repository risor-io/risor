# tablewriter

The `tablewriter` module is used to print rows of data in a
tabular format.

The core functionality is provided by
[github.com/olekukonko/tablewriter](github.com/olekukonko/tablewriter).

There are two ways to use the `tablewriter` module:

- Use the `writer` function to create a tablewriter object and
    use its methods to configure and render the table.
- Call the tablewriter module directly to pass in rows of data and
    render the table in one action.

## Module

```go copy filename="Function signature"
tablewriter(rows [][]string, options map, writer io.writer = os.stdout)
```

Renders a table with the given rows of data. The options map can be used to set
properties of the table, such as the header, footer, and alignment. If a writer
is not provided, it defaults to stdout.

```go copy filename="Example"
>>> tablewriter([["Name", "Age"], ["Alice", 25], ["Bob", 30]])
+-------+-----+
| Name  | Age |
| Alice |  25 |
| Bob   |  30 |
+-------+-----+
```

```go copy filename="Example"
>>> tablewriter([["Alice", 25], ["Bob", 30]], {header: ["Name", "Age"], border: false, row_separator: "="})
  NAME  | AGE  
========+======
  Alice |  25  
  Bob   |  30
```

Available options:

| Option              | Type     | Description                                              |
| ------------------- | -------- | -------------------------------------------------------- |
| header              | []string | The header names for the table.                          |
| footer              | []string | The footer of the table.                                 |
| center_separator    | string   | The center separator character used to separate columns. |
| border              | bool     | Whether to draw a border around the table.               |
| auto_wrap_text      | bool     | Whether to automatically wrap text in the cells.         |
| auto_format_headers | bool     | Whether to automatically format the headers.             |
| header_alignment    | int      | The alignment of the header text.                        |
| header_line         | bool     | Whether to draw a line under the header.                 |
| alignment           | int      | The alignment of the text in the cells.                  |
| row_separator       | string   | The separator character used to separate rows.           |

## Functions

### writer

```go filename="Function signature"
tablewriter.writer(writer io.writer = os.stdout) tablewriter.writer
```

Returns a new tablewriter.writer object that writes to the given output writer.
If an output writer is not provided, it defaults to stdout.

```go copy filename="Example"
>>> w := tablewriter.writer()
>>> w.append_bulk([["Name", "Age"], ["Alice", 25], ["Bob", 30]])
>>> w.render()
+-------+-----+
| Name  | Age |
| Alice |  25 |
| Bob   |  30 |
+-------+-----+
```

## Types

### writer

A table writer object.

#### Methods

##### set_header

```go filename="Method signature"
set_header(header []string)
```

Sets the header text for columns in the table.

##### set_footer

```go filename="Method signature"
set_footer(footer []string)
```

Sets the footer text for columns in the table.

##### set_center_separator

```go filename="Method signature"
set_center_separator(separator string)
```

Sets the character used to separate columns.

##### set_border

```go filename="Method signature"
set_border(border bool)
```

Sets whether to draw a border around the table.

##### append

```go filename="Method signature"
append(row []string)
```

Appends a row to the table.

##### append_bulk

```go filename="Method signature"

append_bulk(rows [][]string)
```

Appends multiple rows to the table.

##### set_auto_wrap_text

```go filename="Method signature"
set_auto_wrap_text(auto_wrap bool)
```

Sets whether to automatically wrap text in the cells.

##### set_auto_format_headers

```go filename="Method signature"
set_auto_format_headers(auto_format bool)
```

Sets whether to automatically format the headers.

##### set_header_alignment

```go filename="Method signature"
set_header_alignment(align int)
```

Sets the alignment of the header text. See the constants section below for
acceptable values.

##### set_header_line

```go filename="Method signature"
set_header_line(header_line bool)
```

Sets whether to draw a line under the header.

##### set_alignment

```go filename="Method signature"
set_alignment(align int)
```

Sets the alignment of the text in the cells. See the constants section below for
acceptable values.

##### set_row_separator

```go filename="Method signature"
set_row_separator(separator string)
```

Sets the character used to separate rows.

##### render

```go filename="Method signature"
render()
```

Renders the table to the writer.

## Constants

### align_default

Specifies the default alignment for the text.

```go filename="Constant"
>>> tablewriter.align_default
0
```

### align_center

Specifies that text should be centered in the cell.

```go filename="Constant"
>>> tablewriter.align_center
1
```

### align_right

Specifies that text should be right-aligned in the cell.

```go filename="Constant"
>>> tablewriter.align_right
2
```

### align_left

Specifies that text should be left-aligned in the cell.

```go filename="Constant"
>>> tablewriter.align_left
3
```
