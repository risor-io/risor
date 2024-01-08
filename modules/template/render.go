package template

import (
	"context"
	"io"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"

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

func newTemplate(name string) *template.Template {
	tpl := template.New(name).Delims(DelimStart, DelimEnd)
	funcMap := sprig.TxtFuncMap()
	funcMap["toYaml"] = toYaml
	funcMap["fromYaml"] = fromYaml
	funcMap["jsonPath"] = jsonPath
	funcMap["k8sLookup"] = k8sLookup
	funcMap["include"] = func(name string, data any) (string, error) {
		buf := new(strings.Builder)
		if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	funcMap["tpl"] = func(tmpl string, data any) (string, error) {
		t, err := template.New("").Funcs(funcMap).Delims(DelimStart, DelimEnd).Parse(tmpl)
		if err != nil {
			return "", err
		}
		buf := new(strings.Builder)
		if err := t.Execute(buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return tpl.Funcs(funcMap)
}

func toYaml(v any) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(b), "\n"), nil
}

func fromYaml(str string) (map[string]any, error) {
	m := map[string]any{}
	err := yaml.Unmarshal([]byte(str), &m)
	return m, err
}
