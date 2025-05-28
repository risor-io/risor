# echarts

The `echarts` module provides a Risor interface for creating charts using the [go-echarts](https://github.com/go-echarts/go-echarts) library. It supports creating various types of charts including bar charts, line charts, scatter plots, pie charts, liquid charts, and heatmaps.

## Key Features

- Support for multiple chart types (bar, line, scatter, pie, liquid, heatmap)
- Consistent API across all chart types
- Chart overlapping functionality for creating composite visualizations
- Configurable chart options (titles, legends, axes)
- Rendering to HTML files via `.render()` method

## Functions

All chart functions follow a consistent pattern: they take data as the first parameter, optional configuration as the second parameter, and return a chart object that can be manipulated and rendered.

### bar

```go filename="Function signature"
bar(data map, options map) chart
```

Creates a new bar chart object.

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
pie(data map, options map) chart
```

Creates a new pie chart object.

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

chart := echarts.pie(data, options)
chart.render("browser_share.html")
```

### liquid

```go filename="Function signature"
liquid(value number, options map) chart
```

Creates a new liquid fill chart object.

```risor copy filename="Example"
chart := echarts.liquid(0.75, {
  title: "Project Completion",
  subtitle: "75% Complete"
})
chart.render("progress.html")
```

The `value` parameter should be a number between 0 and 1 representing the fill percentage.

### heatmap

```go filename="Function signature"
heatmap(data map, options map) chart
```

Creates a new heatmap chart object.

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

chart := echarts.heatmap(data, options)
chart.render("correlation.html")
```

For heatmaps, the `values` array contains [x, y, value] triplets where x and y are coordinates and value is the intensity.

## Chart Objects

All chart functions return chart objects that implement the following methods:

### render

```go filename="Method signature"
chart.render(filename string)
```

Renders the chart to an HTML file.

```risor copy filename="Example"
chart := echarts.bar(data, options)
chart.render("output.html")
```

### overlap

```go filename="Method signature"
chart.overlap(other_chart chart)
```

Overlaps another chart onto this chart, creating a composite visualization. Note that not all chart types support overlapping (e.g., pie and liquid charts cannot be overlapped).

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

**Chart Overlap Compatibility:**
- ✅ Bar charts can overlap and be overlapped
- ✅ Line charts can overlap and be overlapped  
- ✅ Scatter charts can overlap and be overlapped
- ❌ Pie charts cannot overlap or be overlapped
- ❌ Liquid charts cannot overlap or be overlapped
- ❌ Heatmap charts cannot overlap or be overlapped

## Options

All chart functions accept an optional `options` map with the following supported keys:

| Option    | Type     | Description                    | Default Value     |
|-----------|----------|--------------------------------|-------------------|
| `title`   | string   | Main title of the chart        | Chart type + " Chart" |
| `subtitle`| string   | Subtitle text                  | Chart type + " Example" |
| `xlabels` | []string | X-axis category labels         | Empty array       |
| `ylabels` | []string | Y-axis category labels (heatmap only) | Empty array |

## Complete Examples

### Basic Chart Creation

```risor copy filename="Simple chart creation"
// Create a basic bar chart
data := {"Revenue": [100, 150, 200, 175, 225]}
chart := echarts.bar(data, {title: "Monthly Revenue"})
chart.render("revenue.html")

// Create a pie chart
pie_data := {"Product A": 45, "Product B": 30, "Product C": 25}
pie_chart := echarts.pie(pie_data, {title: "Product Distribution"})
pie_chart.render("products.html")

// Create a liquid chart
progress := echarts.liquid(0.68, {title: "Project Progress"})
progress.render("progress.html")
```

### Advanced Overlapping Chart

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

### Multi-Chart Dashboard

```risor copy filename="Creating multiple charts"
// Different chart types for a dashboard
sales := {"Q1": 100, "Q2": 150, "Q3": 200, "Q4": 175}
temperatures := {"Week 1": [20, 22, 25], "Week 2": [18, 21, 24]}
completion := 0.85

// Create individual charts
sales_pie := echarts.pie(sales, {title: "Quarterly Sales"})
temp_line := echarts.line(temperatures, {title: "Temperature Trends"})
progress_liquid := echarts.liquid(completion, {title: "Project Status"})

// Render each chart
sales_pie.render("dashboard_sales.html")
temp_line.render("dashboard_temperature.html")
progress_liquid.render("dashboard_progress.html")
```

This demonstrates the power and flexibility of the echarts module for creating rich, interactive data visualizations with a simple and consistent API.