package echarts

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/risor-io/risor/errz"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

const (
	CHART object.Type = "echarts.chart"
)

type Chart struct {
	chart interface {
		Render(w io.Writer) error
	}
	chartType string
}

func (c *Chart) Type() object.Type {
	return CHART
}

func (c *Chart) Inspect() string {
	return fmt.Sprintf("echarts.%s()", c.chartType)
}

func (c *Chart) Interface() interface{} {
	return c.chart
}

func (c *Chart) IsTruthy() bool {
	return true
}

func (c *Chart) Cost() int {
	return 0
}

func (c *Chart) MarshalJSON() ([]byte, error) {
	return nil, errz.TypeErrorf("type error: unable to marshal %s", CHART)
}

func (c *Chart) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for %s: %v", CHART, opType)
}

func (c *Chart) Equals(other object.Object) object.Object {
	return object.NewBool(c == other)
}

func (c *Chart) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "overlap":
		return object.NewBuiltin("echarts.chart.overlap",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.Errorf("overlap() takes exactly 1 argument (%d given)", len(args))
				}

				otherChart, ok := args[0].(*Chart)
				if !ok {
					return object.Errorf("overlap() argument must be a chart object")
				}

				overlaper, ok := otherChart.chart.(charts.Overlaper)
				if !ok {
					return object.Errorf("chart does not support overlapping")
				}

				overlapable, ok := c.chart.(interface {
					Overlap(chart ...charts.Overlaper)
				})
				if !ok {
					return object.Errorf("this chart type does not support overlapping")
				}

				overlapable.Overlap(overlaper)
				return object.Nil
			}), true
	case "render":
		return object.NewBuiltin("echarts.chart.render",
			func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 1 {
					return object.Errorf("render() takes exactly 1 argument (%d given)", len(args))
				}

				filename, err := object.AsString(args[0])
				if err != nil {
					return err
				}

				f, cerr := os.Create(filename)
				if cerr != nil {
					return object.NewError(cerr)
				}
				defer f.Close()

				if renderErr := c.chart.Render(f); renderErr != nil {
					return object.NewError(renderErr)
				}

				return object.Nil
			}), true
	}
	return nil, false
}

func (c *Chart) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: %s object has no settable attribute %q", CHART, name)
}

func Module() *object.Module {
	return object.NewBuiltinsModule(
		"echarts", map[string]object.Object{
			"bar":     object.NewBuiltin("bar", Bar),
			"line":    object.NewBuiltin("line", Line),
			"pie":     object.NewBuiltin("pie", Pie),
			"liquid":  object.NewBuiltin("liquid", Liquid),
			"heatmap": object.NewBuiltin("heatmap", Heatmap),
			"scatter": object.NewBuiltin("scatter", Scatter),
		},
	)
}

func Line(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.Errorf("missing arguments, at least 1 required")
	}

	data, err := object.AsMap(args[0])
	if err != nil {
		return err
	}

	series := map[string][]opts.LineData{}
	for k, v := range data.Value() {
		items := make([]opts.LineData, 0)
		i, err := object.AsList(v)
		if err != nil {
			return err
		}

		title := object.NewString(k).String()
		for _, v := range i.Value() {
			items = append(items, opts.LineData{Value: v})
		}
		series[title] = items
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 1 {
		options, err = object.AsMap(args[1])
		if err != nil {
			return err
		}
	}

	title, err := strValue(options, "title", "Line Chart")
	if err != nil {
		return err
	}

	subtitle, err := strValue(options, "subtitle", "Line Chart Example")
	if err != nil {
		return err
	}

	xAxis, err := strListValue(options, "xlabels", []string{})
	if err != nil {
		return err
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
		}),
		charts.WithLegendOpts(opts.Legend{Orient: "horizontal", Left: "right", Top: "bottom"}),
	)

	line.SetXAxis(xAxis)
	for t, i := range series {
		line.AddSeries(t, i)
	}

	return &Chart{chart: line, chartType: "line"}
}

func Bar(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.Errorf("missing arguments, at least 1 required")
	}

	data, err := object.AsMap(args[0])
	if err != nil {
		return err
	}

	series := map[string][]opts.BarData{}
	for k, v := range data.Value() {
		items := make([]opts.BarData, 0)
		i, err := object.AsList(v)
		if err != nil {
			return err
		}

		title := object.NewString(k).String()
		for _, v := range i.Value() {
			items = append(items, opts.BarData{Value: v})
		}
		series[title] = items
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 1 {
		options, err = object.AsMap(args[1])
		if err != nil {
			return err
		}
	}

	title, err := strValue(options, "title", "Bar Chart")
	if err != nil {
		return err
	}

	subtitle, err := strValue(options, "subtitle", "Bar Chart Example")
	if err != nil {
		return err
	}

	xAxis, err := strListValue(options, "xlabels", []string{})
	if err != nil {
		return err
	}

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
		}),
		charts.WithLegendOpts(opts.Legend{Orient: "horizontal", Left: "right", Top: "bottom"}),
	)

	bar.SetXAxis(xAxis)
	for t, i := range series {
		bar.AddSeries(t, i)
	}

	return &Chart{chart: bar, chartType: "bar"}
}

func Scatter(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.Errorf("missing arguments, at least 1 required")
	}

	data, err := object.AsMap(args[0])
	if err != nil {
		return err
	}

	series := map[string][]opts.ScatterData{}
	for k, v := range data.Value() {
		items := make([]opts.ScatterData, 0)
		i, err := object.AsList(v)
		if err != nil {
			return err
		}

		title := object.NewString(k).String()
		for _, v := range i.Value() {
			switch item := v.(type) {
			case *object.List:
				if len(item.Value()) >= 2 {
					x := item.Value()[0]
					y := item.Value()[1]
					items = append(items, opts.ScatterData{Value: []interface{}{x, y}})
				}
			default:
				items = append(items, opts.ScatterData{Value: v})
			}
		}
		series[title] = items
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 1 {
		options, err = object.AsMap(args[1])
		if err != nil {
			return err
		}
	}

	title, err := strValue(options, "title", "Scatter Chart")
	if err != nil {
		return err
	}

	subtitle, err := strValue(options, "subtitle", "Scatter Chart Example")
	if err != nil {
		return err
	}

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
		}),
		charts.WithLegendOpts(opts.Legend{Orient: "horizontal", Left: "right", Top: "bottom"}),
	)

	for t, i := range series {
		scatter.AddSeries(t, i)
	}

	return &Chart{chart: scatter, chartType: "scatter"}
}

func strValue(opts *object.Map, key, def string) (string, *object.Error) {
	omap := opts.Value()
	if _, ok := omap[key]; !ok {
		return def, nil
	}

	v, err := object.AsString(omap[key])
	if err != nil {
		return "", err
	}
	if v == "" {
		return def, nil
	}
	return v, nil
}

func Pie(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.Errorf("missing arguments, at least 1 required")
	}

	data, err := object.AsMap(args[0])
	if err != nil {
		return err
	}

	items := make([]opts.PieData, 0)
	for k, v := range data.Value() {
		name := object.NewString(k).String()
		value := v
		items = append(items, opts.PieData{Name: name, Value: value})
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 1 {
		options, err = object.AsMap(args[1])
		if err != nil {
			return err
		}
	}

	title, err := strValue(options, "title", "Pie Chart")
	if err != nil {
		return err
	}

	subtitle, err := strValue(options, "subtitle", "Pie Chart Example")
	if err != nil {
		return err
	}

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
		}),
		charts.WithLegendOpts(opts.Legend{Orient: "vertical", Left: "left"}),
	)

	pie.AddSeries("pie", items)

	return &Chart{chart: pie, chartType: "pie"}
}

func Liquid(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.Errorf("missing arguments, at least 1 required")
	}

	value := args[0]

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 1 {
		var err *object.Error
		options, err = object.AsMap(args[1])
		if err != nil {
			return err
		}
	}

	title, err := strValue(options, "title", "Liquid Chart")
	if err != nil {
		return err
	}

	subtitle, err := strValue(options, "subtitle", "Liquid Chart Example")
	if err != nil {
		return err
	}

	liquid := charts.NewLiquid()
	liquid.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
		}),
	)

	liquid.AddSeries("liquid", []opts.LiquidData{{Value: value}})

	return &Chart{chart: liquid, chartType: "liquid"}
}

func Heatmap(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 1 {
		return object.Errorf("missing arguments, at least 1 required")
	}

	data, err := object.AsMap(args[0])
	if err != nil {
		return err
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 1 {
		options, err = object.AsMap(args[1])
		if err != nil {
			return err
		}
	}

	title, err := strValue(options, "title", "Heatmap")
	if err != nil {
		return err
	}

	subtitle, err := strValue(options, "subtitle", "Heatmap Example")
	if err != nil {
		return err
	}

	xAxis, err := strListValue(options, "xlabels", []string{})
	if err != nil {
		return err
	}

	yAxis, err := strListValue(options, "ylabels", []string{})
	if err != nil {
		return err
	}

	heatmapData := make([]opts.HeatMapData, 0)
	dataMap := data.Value()

	if values, ok := dataMap["values"]; ok {
		valuesList, err := object.AsList(values)
		if err != nil {
			return err
		}

		for _, item := range valuesList.Value() {
			itemList, err := object.AsList(item)
			if err != nil {
				return err
			}

			if len(itemList.Value()) >= 3 {
				x := itemList.Value()[0]
				y := itemList.Value()[1]
				value := itemList.Value()[2]
				heatmapData = append(heatmapData, opts.HeatMapData{Value: [3]interface{}{x, y, value}})
			}
		}
	}

	heatmap := charts.NewHeatMap()
	heatmap.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subtitle,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "category",
			Data:      xAxis,
			SplitArea: &opts.SplitArea{Show: opts.Bool(true)},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "category",
			Data:      yAxis,
			SplitArea: &opts.SplitArea{Show: opts.Bool(false)},
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: opts.Bool(true),
			Orient:     "horizontal",
			Left:       "center",
			Bottom:     "15%",
		}),
	)

	heatmap.AddSeries("heatmap", heatmapData)

	return &Chart{chart: heatmap, chartType: "heatmap"}
}

func strListValue(opts *object.Map, key string, def []string) ([]string, *object.Error) {
	omap := opts.Value()
	if _, ok := omap[key]; !ok {
		return def, nil
	}

	v, err := object.AsStringSlice(omap[key])
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return def, nil
	}
	return v, nil
}
