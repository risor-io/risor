package template

import (
	"context"
	"io"
	"strings"

	"github.com/risor-io/risor/os"
)

const (
	DelimStart = "{{"
	DelimEnd   = "}}"
)

// Render is a go template rendering function it includes all the sprig lib functions
// as well as some extras like a k8sLookup function to get values from k8s objects
// you can access environment variables from the template under .Env
// The passed values will be available under .Values in the templates
func Render(ctx context.Context, out io.Writer, tmpl string, values any) error {
	t, err := newTemplate("tmpl").Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(out, map[string]any{
		"Env":    envMap(ctx),
		"Values": values,
	})
}

func envMap(ctx context.Context) map[string]string {
	envMap := make(map[string]string)

	for _, v := range os.GetDefaultOS(ctx).Environ() {
		kv := strings.SplitN(v, "=", 2)
		envMap[kv[0]] = kv[1]
	}
	return envMap
}
