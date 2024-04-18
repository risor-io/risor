package tablewriter

import (
	"bytes"
	"context"
	"testing"

	"github.com/olekukonko/tablewriter"
	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	bufObj := object.NewBuffer(buf)
	w := CreateWriter(context.Background(), bufObj)
	require.NotNil(t, w)

	rows := object.NewList(
		[]object.Object{
			object.NewList([]object.Object{
				object.NewString("a"),
				object.NewString("b"),
				object.NewString("c"),
			}),
			object.NewList([]object.Object{
				object.NewString("1"),
				object.NewString("2"),
				object.NewString("3"),
			}),
		},
	)

	opts := object.NewMap(map[string]object.Object{
		"header": object.NewList([]object.Object{
			object.NewString("H1"),
			object.NewString("H2"),
			object.NewString("H3"),
		}),
		"footer": object.NewList([]object.Object{
			object.NewString("F1"),
			object.NewString("F2"),
			object.NewString("F3"),
		}),
		"alignment": object.NewInt(int64(tablewriter.ALIGN_RIGHT)),
	})

	r := Render(context.Background(), rows, opts, bufObj)
	require.NotNil(t, r)

	require.Equal(t, `+----+----+----+
| H1 | H2 | H3 |
+----+----+----+
|  a |  b |  c |
|  1 |  2 |  3 |
+----+----+----+
| F1 | F2 | F3 |
+----+----+----+
`, buf.String())
}
