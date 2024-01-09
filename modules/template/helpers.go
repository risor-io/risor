package template

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

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
