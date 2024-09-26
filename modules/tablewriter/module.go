package tablewriter

import (
	"context"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/os"
)

func CreateWriter(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("tablewriter.writer", 0, 1, args); err != nil {
		return err
	}
	writer, err := getDestination(ctx, args...)
	if err != nil {
		return err
	}
	return NewWriter(tablewriter.NewWriter(writer))
}

func Render(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.RequireRange("tablewriter", 1, 3, args); err != nil {
		return err
	}
	// First argument: list of lists of strings
	rows, err := object.AsList(args[0])
	if err != nil {
		return err
	}
	// Optional second argument: map of options
	var opts map[string]object.Object
	if len(args) > 1 {
		m, err := object.AsMap(args[1])
		if err != nil {
			return err
		}
		opts = m.Value()
	}
	// Optional third argument: destination
	var ioWriter io.Writer
	if len(args) > 2 {
		ioWriter, err = getDestination(ctx, args[2])
		if err != nil {
			return err
		}
	} else {
		ioWriter = os.GetDefaultOS(ctx).Stdout()
	}
	obj := NewWriter(tablewriter.NewWriter(ioWriter))
	for k, v := range opts {
		switch k {
		case "header":
			if err, ok := obj.setHeader(v).(*object.Error); ok {
				return err
			}
		case "footer":
			if err, ok := obj.setFooter(v).(*object.Error); ok {
				return err
			}
		case "center_separator":
			if err, ok := obj.setCenterSeparator(v).(*object.Error); ok {
				return err
			}
		case "border":
			if err, ok := obj.setBorder(v).(*object.Error); ok {
				return err
			}
		case "auto_wrap_text":
			if err, ok := obj.setAutoWrapText(v).(*object.Error); ok {
				return err
			}
		case "auto_format_headers":
			if err, ok := obj.setAutoFormatHeaders(v).(*object.Error); ok {
				return err
			}
		case "header_alignment":
			if err, ok := obj.setHeaderAlignment(v).(*object.Error); ok {
				return err
			}
		case "header_line":
			if err, ok := obj.setHeaderLine(v).(*object.Error); ok {
				return err
			}
		case "alignment":
			if err, ok := obj.setAlignment(v).(*object.Error); ok {
				return err
			}
		case "row_separator":
			if err, ok := obj.setRowSeparator(v).(*object.Error); ok {
				return err
			}
		default:
			return object.Errorf("argument error: unknown option %s", k)
		}
	}
	obj.appendBulk(rows)
	obj.Render()
	return object.Nil
}

func getDestination(ctx context.Context, args ...object.Object) (io.Writer, *object.Error) {
	if len(args) == 0 {
		return os.GetDefaultOS(ctx).Stdout(), nil
	}
	return object.AsWriter(args[0])
}

func Module() *object.Module {
	return object.NewBuiltinsModule("tablewriter", map[string]object.Object{
		"writer":        object.NewBuiltin("writer", CreateWriter),
		"align_center":  object.NewInt(int64(tablewriter.ALIGN_CENTER)),
		"align_left":    object.NewInt(int64(tablewriter.ALIGN_LEFT)),
		"align_right":   object.NewInt(int64(tablewriter.ALIGN_RIGHT)),
		"align_default": object.NewInt(int64(tablewriter.ALIGN_DEFAULT)),
	}, Render)
}
