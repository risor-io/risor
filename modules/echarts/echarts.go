package echarts

import (
	"context"
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return object.NewBuiltinsModule(
		"echarts", map[string]object.Object{
			"bar":     object.NewBuiltin("bar", Bar),
			"line":    object.NewBuiltin("line", Line),
			"pie":     object.NewBuiltin("pie", Pie),
			"liquid":  object.NewBuiltin("liquid", Liquid),
			"heatmap": object.NewBuiltin("heatmap", Heatmap),
		},
	)
}

func require(funcName string, count int, args []object.Object) *object.Error {
	nArgs := len(args)
	if nArgs != count {
		if count == 1 {
			return object.Errorf(
				fmt.Sprintf("type error: %s() takes exactly 1 argument (%d given)",
					funcName, nArgs))
		}
		return object.Errorf(
			fmt.Sprintf("type error: %s() takes exactly %d arguments (%d given)",
				funcName, count, nArgs))
	}
	return nil
}

func Line(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.Errorf("missing arguments, 2 required")
	}

	data, err := object.AsMap(args[1])
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

	file, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 2 {
		options, err = object.AsMap(args[2])
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

	f, ferr := os.Create(file)
	if ferr != nil {
		return object.NewError(ferr)
	}
	defer f.Close()

	nErr := line.Render(f)
	if nErr != nil {
		return object.NewError(nErr)
	}
	return nil
}

func Bar(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.Errorf("missing arguments, 2 required")
	}

	data, err := object.AsMap(args[1])
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

	file, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 2 {
		options, err = object.AsMap(args[2])
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

	f, cerr := os.Create(file)
	if cerr != nil {
		return object.NewError(cerr)
	}
	defer f.Close()

	nErr := bar.Render(f)
	if nErr != nil {
		return object.NewError(nErr)
	}
	return nil
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
	if len(args) < 2 {
		return object.Errorf("missing arguments, 2 required")
	}

	data, err := object.AsMap(args[1])
	if err != nil {
		return err
	}

	items := make([]opts.PieData, 0)
	for k, v := range data.Value() {
		name := object.NewString(k).String()
		value := v
		items = append(items, opts.PieData{Name: name, Value: value})
	}

	file, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 2 {
		options, err = object.AsMap(args[2])
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

	f, cerr := os.Create(file)
	if cerr != nil {
		return object.NewError(cerr)
	}
	defer f.Close()

	nErr := pie.Render(f)
	if nErr != nil {
		return object.NewError(nErr)
	}
	return nil
}

func Liquid(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.Errorf("missing arguments, 2 required")
	}

	value := args[1]

	file, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 2 {
		options, err = object.AsMap(args[2])
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

	f, cerr := os.Create(file)
	if cerr != nil {
		return object.NewError(cerr)
	}
	defer f.Close()

	nErr := liquid.Render(f)
	if nErr != nil {
		return object.NewError(nErr)
	}
	return nil
}

func Heatmap(ctx context.Context, args ...object.Object) object.Object {
	if len(args) < 2 {
		return object.Errorf("missing arguments, 2 required")
	}

	data, err := object.AsMap(args[1])
	if err != nil {
		return err
	}

	file, err := object.AsString(args[0])
	if err != nil {
		return err
	}

	options := object.NewMap(map[string]object.Object{})
	if len(args) > 2 {
		options, err = object.AsMap(args[2])
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

	// Extract heatmap data from the input map
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

	f, cerr := os.Create(file)
	if cerr != nil {
		return object.NewError(cerr)
	}
	defer f.Close()

	nErr := heatmap.Render(f)
	if nErr != nil {
		return object.NewError(nErr)
	}
	return nil
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
