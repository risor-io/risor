# echarts

The `echarts` module provides a Risor interface for creating charts using the [go-echarts](https://github.com/go-echarts/go-echarts) library. It supports creating various types of charts including bar charts, line charts, scatter plots, pie charts, liquid charts, and heatmaps.

## Key Features

- Support for multiple chart types (bar, line, scatter, pie, liquid, heatmap)
- Chart overlapping functionality for creating composite visualizations
- Configurable chart options (titles, legends, axes)
- Direct rendering to HTML files

## Functions

### bar

```go filename="Function signature"
bar(data map, options map) chart
```

Creates a new bar chart object. Returns a chart object that can be rendered or overlapped with other charts.

```risor copy filename="Example"
data := {
  "Sales": [120, 200, 150, 80, 70, 110, 130],
  "Revenue": [100, 190, 140, 70, 60, 100, 120]
}

options := {
  title: "Monthly Sales Report",
  subtitle: "Q1 2024 Performance",
  xlabels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul"]
}

chart := echarts.bar(data, options)
chart.render("sales_chart.html")
```

The `data` parameter is a map where keys are series names and values are lists of numeric data.

### line

```go filename="Function signature"
line(data map, options map) chart
```

Creates a new line chart object.

```risor copy filename="Example"
data := {
  "Temperature": [20, 22, 25, 28, 30, 32, 29],
  "Humidity": [60, 65, 70, 68, 72, 75, 73]
}

options := {
  title: "Weather Data",
  subtitle: "Daily measurements",
  xlabels: ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]
}

chart := echarts.line(data, options)
chart.render("weather_chart.html")
```

### scatter

```go filename="Function signature"
scatter(data map, options map) chart
```

Creates a new scatter plot chart object.

```risor copy filename="Example"
data := {
  "Dataset A": [[1, 4], [2, 6], [3, 8], [4, 10], [5, 12]],
  "Dataset B": [[1, 2], [2, 5], [3, 3], [4, 8], [5, 7]]
}

options := {
  title: "Scatter Plot Analysis",
  subtitle: "Correlation study"
}

chart := echarts.scatter(data, options)
chart.render("scatter_chart.html")
```

For scatter plots, data values should be lists of [x, y] coordinates.

### pie

```go filename="Function signature"
pie(filename string, data map, options map)
```

Creates and renders a pie chart directly to a file.

```risor copy filename="Example"
data := {
  "Chrome": 60,
  "Firefox": 20,
  "Safari": 15,
  "Edge": 5
}

options := {
  title: "Browser Market Share",
  subtitle: "2024 Statistics"
}

echarts.pie("browser_share.html", data, options)
```

### liquid

```go filename="Function signature"
liquid(filename string, value number, options map)
```

Creates and renders a liquid fill chart directly to a file.

```risor copy filename="Example"
echarts.liquid("progress.html", 0.75, {
  title: "Project Completion",
  subtitle: "75% Complete"
})
```

### heatmap

```go filename="Function signature"
heatmap(filename string, data map, options map)
```

Creates and renders a heatmap chart directly to a file.

```risor copy filename="Example"
data := {
  values: [
    [0, 0, 5], [0, 1, 1], [0, 2, 0],
    [1, 0, 1], [1, 1, 3], [1, 2, 0],
    [2, 0, 0], [2, 1, 2], [2, 2, 4]
  ]
}

options := {
  title: "Correlation Matrix",
  subtitle: "Feature relationships",
  xlabels: ["Feature A", "Feature B", "Feature C"],
  ylabels: ["Metric X", "Metric Y", "Metric Z"]
}

echarts.heatmap("correlation.html", data, options)
```

For heatmaps, the `values` array contains [x, y, value] triplets.

## Chart Objects

Chart objects returned by `bar()`, `line()`, and `scatter()` functions have the following methods:

### overlap

```go filename="Method signature"
chart.overlap(other_chart chart)
```

Overlaps another chart onto this chart, creating a composite visualization.

```risor copy filename="Example"
// Create a bar chart
bar_data := {"Sales": [120, 200, 150, 80, 70]}
bar_chart := echarts.bar(bar_data, {
  title: "Sales and Trends",
  xlabels: ["Q1", "Q2", "Q3", "Q4", "Q5"]
})

// Create a line chart
line_data := {"Trend": [100, 180, 160, 90, 85]}
line_chart := echarts.line(line_data)

// Overlap the line chart onto the bar chart
bar_chart.overlap(line_chart)

// Render the combined chart
bar_chart.render("combined_chart.html")
```

### render

```go filename="Method signature"
chart.render(filename string)
```

Renders the chart to an HTML file.

```risor copy filename="Example"
chart := echarts.bar(data, options)
chart.render("output.html")
```

## Options

All chart functions accept an optional `options` map with the following supported keys:

| Option    | Type     | Description                    | Default Value     |
|-----------|----------|--------------------------------|-------------------|
| `title`   | string   | Main title of the chart        | Chart type + "Chart" |
| `subtitle`| string   | Subtitle text                  | Chart type + "Example" |
| `xlabels` | []string | X-axis category labels         | Empty array       |
| `ylabels` | []string | Y-axis category labels (heatmap only) | Empty array |

## Complete Example

```risor copy filename="Complex overlapping chart example"
// Sales data
sales_data := {
  "Actual Sales": [120, 200, 150, 80, 70, 110, 130],
  "Target Sales": [100, 180, 160, 90, 80, 120, 140]
}

// Create bar chart for sales
sales_chart := echarts.bar(sales_data, {
  title: "Sales Performance Dashboard",
  subtitle: "Q1 2024 - Actual vs Target",
  xlabels: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul"]
})

// Trend line data
trend_data := {
  "Growth Trend": [110, 190, 155, 85, 75, 115, 135]
}

// Create line chart for trend
trend_chart := echarts.line(trend_data)

// Overlay trend on sales chart
sales_chart.overlap(trend_chart)

// Render the combined chart
sales_chart.render("sales_dashboard.html")
```

This creates a comprehensive chart combining bar and line visualizations, demonstrating the power of the overlap functionality for creating rich, multi-layered data visualizations.
