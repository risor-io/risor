# Echarts 

Risor module to create charts. Wraps https://github.com/go-echarts/go-echarts.

The `echarts` module exposes a simple interface to create charts, powered by the great [go-echarts](https://github.com/go-echarts/go-echarts) library.

## Functions

### bar

```go filename="Function signature"
bar(file string, data map, opts map)
```

Creates a new bar chart.

```go copy filename="Example"
data := {
  "serie A": [1, 2, 3],
  "serie B": [3, 4, 5],
}

echarts.bar(
	"bar.html",
	data,
)
```

The `opts` argument may be a map containing any of the following keys:

| Name   | Type                          | Description                              |
| ------ | ----------------------------- | ---------------------------------------- |
| title  | string                        | The title of the chart                   |
| subtitle | string                      | The subtitle of the chart                |
| xlabels | []string                      | The labels for the x-axis                |

```go copy filename="Example"
data := {
  "serie A": [1, 2, 3],
  "serie B": [3, 4, 5],
}

echarts.bar(
	"bar.html",
	data,
	{
		title: "My awesome bar chart",
		subtitle: "this is a subtitle",
		xlabels: ["one", "two", "three"]
	},
)
```

### line

```go filename="Function signature"
line(file string, data map, opts map)
```

Creates a new line chart.

```go copy filename="Example"
data := {
  "serie A": [1, 2, 3],
  "serie B": [3, 4, 5],
}

echarts.line(
	"line.html",
	data,
)
```

The `opts` argument may be a map containing any of the following keys:

| Name   | Type                          | Description                              |
| ------ | ----------------------------- | ---------------------------------------- |
| title  | string                        | The title of the chart                   |
| subtitle | string                      | The subtitle of the chart                |
| xlabels | []string                      | The labels for the x-axis                |


```go copy filename="Example"
data := {
  "serie A": [1, 2, 3],
  "serie B": [3, 4, 5],
}

echarts.line(
	"line.html",
	data,
	{
		title: "My awesome line chart",
		subtitle: "this is a subtitle",
		xlabels: ["one", "two", "three"]
	},
)
```

### pie

```go filename="Function signature"
pie(file string, data map, opts map)
```

Creates a new pie chart.

```go copy filename="Example"
data := {
  "Apples": 30,
  "Oranges": 20,
  "Bananas": 25,
  "Grapes": 15,
}

echarts.pie(
	"pie.html",
	data,
)
```

The `opts` argument may be a map containing any of the following keys:

| Name   | Type                          | Description                              |
| ------ | ----------------------------- | ---------------------------------------- |
| title  | string                        | The title of the chart                   |
| subtitle | string                      | The subtitle of the chart                |

```go copy filename="Example"
data := {
  "Apples": 30,
  "Oranges": 20,
  "Bananas": 25,
  "Grapes": 15,
}

echarts.pie(
	"pie.html",
	data,
	{
		title: "Fruit Distribution",
		subtitle: "Sales by fruit type"
	},
)
```

### liquid

```go filename="Function signature"
liquid(file string, value number, opts map)
```

Creates a new liquid chart (also known as a liquid fill gauge).

```go copy filename="Example"
echarts.liquid(
	"liquid.html",
	0.6,
)
```

The `opts` argument may be a map containing any of the following keys:

| Name   | Type                          | Description                              |
| ------ | ----------------------------- | ---------------------------------------- |
| title  | string                        | The title of the chart                   |
| subtitle | string                      | The subtitle of the chart                |

```go copy filename="Example"
echarts.liquid(
	"liquid.html",
	0.75,
	{
		title: "Progress Indicator",
		subtitle: "75% Complete"
	},
)
```

### heatmap

```go filename="Function signature"
heatmap(file string, data map, opts map)
```

Creates a new heatmap chart.

```go copy filename="Example"
data := {
  "values": [
    [0, 0, 10],
    [0, 1, 19],
    [0, 2, 8],
    [1, 0, 12],
    [1, 1, 15],
    [1, 2, 6],
    [2, 0, 4],
    [2, 1, 7],
    [2, 2, 20],
  ]
}

echarts.heatmap(
	"heatmap.html",
	data,
)
```

The `opts` argument may be a map containing any of the following keys:

| Name   | Type                          | Description                              |
| ------ | ----------------------------- | ---------------------------------------- |
| title  | string                        | The title of the chart                   |
| subtitle | string                      | The subtitle of the chart                |
| xlabels | []string                      | The labels for the x-axis                |
| ylabels | []string                      | The labels for the y-axis                |

```go copy filename="Example"
data := {
  "values": [
    [0, 0, 10],
    [0, 1, 19],
    [0, 2, 8],
    [1, 0, 12],
    [1, 1, 15],
    [1, 2, 6],
    [2, 0, 4],
    [2, 1, 7],
    [2, 2, 20],
  ]
}

echarts.heatmap(
	"heatmap.html",
	data,
	{
		title: "Temperature Map",
		subtitle: "Daily temperature readings",
		xlabels: ["Monday", "Tuesday", "Wednesday"],
		ylabels: ["Morning", "Afternoon", "Evening"]
	},
)
```