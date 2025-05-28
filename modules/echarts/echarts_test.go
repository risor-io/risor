package echarts

import (
	"context"
	"os"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/assert"
)

func TestBarChart(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Series A": object.NewList([]object.Object{
			object.NewInt(10),
			object.NewInt(20),
			object.NewInt(30),
		}),
		"Series B": object.NewList([]object.Object{
			object.NewInt(15),
			object.NewInt(25),
			object.NewInt(35),
		}),
	})

	options := object.NewMap(map[string]object.Object{
		"title":    object.NewString("Test Bar Chart"),
		"subtitle": object.NewString("Test Subtitle"),
		"xlabels": object.NewList([]object.Object{
			object.NewString("A"),
			object.NewString("B"),
			object.NewString("C"),
		}),
	})

	chart := Bar(context.Background(), data, options)
	assert.False(t, object.IsError(chart))
	assert.Equal(t, CHART, chart.Type())

	chartObj := chart.(*Chart)
	assert.Equal(t, "bar", chartObj.chartType)
}

func TestLineChart(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Series A": object.NewList([]object.Object{
			object.NewInt(10),
			object.NewInt(20),
			object.NewInt(30),
		}),
	})

	chart := Line(context.Background(), data)
	assert.False(t, object.IsError(chart))
	assert.Equal(t, CHART, chart.Type())

	chartObj := chart.(*Chart)
	assert.Equal(t, "line", chartObj.chartType)
}

func TestScatterChart(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Series A": object.NewList([]object.Object{
			object.NewList([]object.Object{
				object.NewInt(1),
				object.NewInt(2),
			}),
			object.NewList([]object.Object{
				object.NewInt(3),
				object.NewInt(4),
			}),
		}),
	})

	chart := Scatter(context.Background(), data)
	assert.False(t, object.IsError(chart))
	assert.Equal(t, CHART, chart.Type())

	chartObj := chart.(*Chart)
	assert.Equal(t, "scatter", chartObj.chartType)
}

func TestChartOverlap(t *testing.T) {
	barData := object.NewMap(map[string]object.Object{
		"Bars": object.NewList([]object.Object{
			object.NewInt(10),
			object.NewInt(20),
			object.NewInt(30),
		}),
	})

	lineData := object.NewMap(map[string]object.Object{
		"Lines": object.NewList([]object.Object{
			object.NewInt(15),
			object.NewInt(25),
			object.NewInt(35),
		}),
	})

	barChart := Bar(context.Background(), barData)
	lineChart := Line(context.Background(), lineData)

	assert.False(t, object.IsError(barChart))
	assert.False(t, object.IsError(lineChart))

	overlapMethod, ok := barChart.GetAttr("overlap")
	assert.True(t, ok)

	result := overlapMethod.(*object.Builtin).Call(context.Background(), lineChart)
	assert.Equal(t, object.Nil, result)
}

func TestChartRender(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Test": object.NewList([]object.Object{
			object.NewInt(1),
			object.NewInt(2),
			object.NewInt(3),
		}),
	})

	chart := Bar(context.Background(), data)
	assert.False(t, object.IsError(chart))

	renderMethod, ok := chart.GetAttr("render")
	assert.True(t, ok)

	filename := "test_chart.html"
	result := renderMethod.(*object.Builtin).Call(context.Background(), object.NewString(filename))
	assert.Equal(t, object.Nil, result)

	_, err := os.Stat(filename)
	assert.NoError(t, err)

	os.Remove(filename)
}

func TestPieChart(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Apple":  object.NewInt(30),
		"Orange": object.NewInt(20),
		"Banana": object.NewInt(25),
	})

	filename := "test_pie.html"
	result := Pie(context.Background(), object.NewString(filename), data)
	assert.Equal(t, object.Nil, result)

	_, err := os.Stat(filename)
	assert.NoError(t, err)

	os.Remove(filename)
}

func TestLiquidChart(t *testing.T) {
	value := object.NewFloat(0.75)

	filename := "test_liquid.html"
	result := Liquid(context.Background(), object.NewString(filename), value)
	assert.Equal(t, object.Nil, result)

	_, err := os.Stat(filename)
	assert.NoError(t, err)

	os.Remove(filename)
}

func TestHeatmapChart(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"values": object.NewList([]object.Object{
			object.NewList([]object.Object{
				object.NewInt(0),
				object.NewInt(0),
				object.NewInt(10),
			}),
			object.NewList([]object.Object{
				object.NewInt(1),
				object.NewInt(1),
				object.NewInt(20),
			}),
		}),
	})

	filename := "test_heatmap.html"
	result := Heatmap(context.Background(), object.NewString(filename), data)
	assert.Equal(t, object.Nil, result)

	_, err := os.Stat(filename)
	assert.NoError(t, err)

	os.Remove(filename)
}

func TestChartEquals(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Test": object.NewList([]object.Object{object.NewInt(1)}),
	})

	chart1 := Bar(context.Background(), data)
	chart2 := Bar(context.Background(), data)

	assert.True(t, chart1.Equals(chart1).(*object.Bool).Value())
	assert.False(t, chart1.Equals(chart2).(*object.Bool).Value())
}

func TestInvalidOperations(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Test": object.NewList([]object.Object{object.NewInt(1)}),
	})

	chart := Bar(context.Background(), data)
	
	result := chart.SetAttr("invalid", object.NewString("test"))
	assert.Error(t, result)

	_, ok := chart.GetAttr("nonexistent")
	assert.False(t, ok)
}

func TestInvalidChartArguments(t *testing.T) {
	result := Bar(context.Background())
	assert.True(t, object.IsError(result))

	result = Line(context.Background())
	assert.True(t, object.IsError(result))

	result = Scatter(context.Background())
	assert.True(t, object.IsError(result))

	result = Pie(context.Background(), object.NewString("test.html"))
	assert.True(t, object.IsError(result))

	result = Liquid(context.Background(), object.NewString("test.html"))
	assert.True(t, object.IsError(result))

	result = Heatmap(context.Background(), object.NewString("test.html"))
	assert.True(t, object.IsError(result))
}

func TestOverlapWithIncompatibleChart(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Test": object.NewList([]object.Object{object.NewInt(1)}),
	})

	chart := Bar(context.Background(), data)
	otherObject := object.NewString("not a chart")

	overlapMethod, ok := chart.GetAttr("overlap")
	assert.True(t, ok)

	result := overlapMethod.(*object.Builtin).Call(context.Background(), otherObject)
	assert.True(t, object.IsError(result))
}

func TestChartWithOptions(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Sales": object.NewList([]object.Object{
			object.NewInt(100),
			object.NewInt(150),
			object.NewInt(200),
		}),
	})

	options := object.NewMap(map[string]object.Object{
		"title":    object.NewString("Monthly Sales"),
		"subtitle": object.NewString("Q1 2024"),
		"xlabels": object.NewList([]object.Object{
			object.NewString("Jan"),
			object.NewString("Feb"),
			object.NewString("Mar"),
		}),
	})

	chart := Bar(context.Background(), data, options)
	assert.False(t, object.IsError(chart))
	assert.Equal(t, CHART, chart.Type())
}

func TestChartInspect(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Test": object.NewList([]object.Object{object.NewInt(1)}),
	})

	barChart := Bar(context.Background(), data)
	assert.Equal(t, "echarts.bar()", barChart.Inspect())

	lineChart := Line(context.Background(), data)
	assert.Equal(t, "echarts.line()", lineChart.Inspect())

	scatterChart := Scatter(context.Background(), data)
	assert.Equal(t, "echarts.scatter()", scatterChart.Inspect())
}

func TestChartProperties(t *testing.T) {
	data := object.NewMap(map[string]object.Object{
		"Test": object.NewList([]object.Object{object.NewInt(1)}),
	})

	chart := Bar(context.Background(), data)
	assert.True(t, chart.IsTruthy())
	assert.Equal(t, 0, chart.Cost())
}