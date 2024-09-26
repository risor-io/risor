package tablewriter

import (
	"context"
	"fmt"

	"github.com/olekukonko/tablewriter"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

var _ object.Object = (*Writer)(nil)

const WRITER object.Type = "tablewriter.writer"

type Writer struct {
	value *tablewriter.Table
}

func (w *Writer) IsTruthy() bool {
	return true
}

func (w *Writer) Type() object.Type {
	return WRITER
}

func (w *Writer) Value() *tablewriter.Table {
	return w.value
}

func (w *Writer) Inspect() string {
	return fmt.Sprintf("%s()", WRITER)
}

func (w *Writer) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: cannot set %q on %s object", name, WRITER)
}

func (w *Writer) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "set_header":
		return object.NewBuiltin("set_header", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setHeader(args...)
		}), true
	case "set_footer":
		return object.NewBuiltin("set_footer", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setFooter(args...)
		}), true
	case "set_center_separator":
		return object.NewBuiltin("set_center_separator", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setCenterSeparator(args...)
		}), true
	case "set_border":
		return object.NewBuiltin("set_border", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setBorder(args...)
		}), true
	case "append":
		return object.NewBuiltin("append", func(ctx context.Context, args ...object.Object) object.Object {
			return w.append(args...)
		}), true
	case "append_bulk":
		return object.NewBuiltin("append_bulk", func(ctx context.Context, args ...object.Object) object.Object {
			return w.appendBulk(args...)
		}), true
	case "set_auto_wrap_text":
		return object.NewBuiltin("set_auto_wrap_text", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setAutoWrapText(args...)
		}), true
	case "set_auto_format_headers":
		return object.NewBuiltin("set_auto_format_headers", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setAutoFormatHeaders(args...)
		}), true
	case "set_header_alignment":
		return object.NewBuiltin("set_header_alignment", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setHeaderAlignment(args...)
		}), true
	case "set_header_line":
		return object.NewBuiltin("set_header_line", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setHeaderLine(args...)
		}), true
	case "set_alignment":
		return object.NewBuiltin("set_alignment", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setAlignment(args...)
		}), true
	case "set_row_separator":
		return object.NewBuiltin("set_row_separator", func(ctx context.Context, args ...object.Object) object.Object {
			return w.setRowSeparator(args...)
		}), true
	case "render":
		return object.NewBuiltin("render", func(ctx context.Context, args ...object.Object) object.Object {
			if err := arg.Require("tablewriter.writer.render", 0, args); err != nil {
				return err
			}
			w.value.Render()
			return object.Nil
		}), true
	default:
		return nil, false
	}
}

func (w *Writer) Render() {
	w.value.Render()
}

func (w *Writer) Interface() interface{} {
	return w.value
}

func (w *Writer) Equals(other object.Object) object.Object {
	return object.NewBool(w == other)
}

func (w *Writer) Cost() int {
	return 0
}

func (w *Writer) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.EvalErrorf("eval error: unsupported operation for %s: %v", WRITER, opType)
}

func NewWriter(v *tablewriter.Table) *Writer {
	return &Writer{value: v}
}

func stringsFromList(l *object.List) []string {
	var strs []string
	for _, item := range l.Value() {
		switch item := item.(type) {
		case fmt.Stringer:
			strs = append(strs, item.String())
		default:
			strs = append(strs, item.Inspect())
		}
	}
	return strs
}

func (w *Writer) setHeader(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_header", 1, args); err != nil {
		return err
	}
	items, err := object.AsStringSlice(args[0])
	if err != nil {
		return err
	}
	w.value.SetHeader(items)
	return object.Nil
}

func (w *Writer) setFooter(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_footer", 1, args); err != nil {
		return err
	}
	items, err := object.AsStringSlice(args[0])
	if err != nil {
		return err
	}
	w.value.SetFooter(items)
	return object.Nil
}

func (w *Writer) setCenterSeparator(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_center_separator", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	w.value.SetCenterSeparator(s)
	return object.Nil
}

func (w *Writer) setBorder(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_border", 1, args); err != nil {
		return err
	}
	b, err := object.AsBool(args[0])
	if err != nil {
		return err
	}
	w.value.SetBorder(b)
	return object.Nil
}

func (w *Writer) append(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.append", 1, args); err != nil {
		return err
	}
	items, err := object.AsList(args[0])
	if err != nil {
		return err
	}
	w.value.Append(stringsFromList(items))
	return object.Nil
}

func (w *Writer) appendBulk(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.append_bulk", 1, args); err != nil {
		return err
	}
	outer, err := object.AsList(args[0])
	if err != nil {
		return err
	}
	var strs [][]string
	for _, inner := range outer.Value() {
		switch inner := inner.(type) {
		case *object.List:
			strs = append(strs, stringsFromList(inner))
		default:
			return object.TypeErrorf("type error: expected list (got %s)", inner.Type())
		}
	}
	w.value.AppendBulk(strs)
	return object.Nil
}

func (w *Writer) setAutoWrapText(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_auto_wrap_text", 1, args); err != nil {
		return err
	}
	b, err := object.AsBool(args[0])
	if err != nil {
		return err
	}
	w.value.SetAutoWrapText(b)
	return object.Nil
}

func (w *Writer) setAutoFormatHeaders(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_auto_format_headers", 1, args); err != nil {
		return err
	}
	b, err := object.AsBool(args[0])
	if err != nil {
		return err
	}
	w.value.SetAutoFormatHeaders(b)
	return object.Nil
}

func (w *Writer) setHeaderAlignment(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_header_alignment", 1, args); err != nil {
		return err
	}
	a, err := object.AsInt(args[0])
	if err != nil {
		return err
	}
	w.value.SetHeaderAlignment(int(a))
	return object.Nil
}

func (w *Writer) setHeaderLine(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_header_line", 1, args); err != nil {
		return err
	}
	b, err := object.AsBool(args[0])
	if err != nil {
		return err
	}
	w.value.SetHeaderLine(b)
	return object.Nil
}

func (w *Writer) setAlignment(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_alignment", 1, args); err != nil {
		return err
	}
	a, err := object.AsInt(args[0])
	if err != nil {
		return err
	}
	w.value.SetAlignment(int(a))
	return object.Nil
}

func (w *Writer) setRowSeparator(args ...object.Object) object.Object {
	if err := arg.Require("tablewriter.writer.set_row_separator", 1, args); err != nil {
		return err
	}
	s, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	w.value.SetRowSeparator(s)
	return object.Nil
}
