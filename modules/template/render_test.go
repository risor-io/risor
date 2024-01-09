package template

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     any
		expected string
	}{
		{
			name:     "use env",
			template: `var FOO is {{ .Env.FOO }}`,
			expected: "var FOO is BAR",
		},
		{
			name: "use env and values",
			template: `var FOO is {{ .Env.FOO }}
value BAR is {{ .Values.BAR }}`,
			expected: "var FOO is BAR\nvalue BAR is FOO",
			data: map[string]string{
				"BAR": "FOO",
			},
		},
		{
			name: "jsonpath",
			data: map[string]any{
				"data": map[string]string{
					"key1": "val1",
				},
			},
			template: `key is {{ .Values | jsonPath ".data.key1" }}`,
			expected: "key is val1",
		},
		{
			name: "to yaml",
			data: struct {
				Foo    string
				Bar    string
				Foobar map[string]string
			}{
				Foo: "bar",
				Bar: "foo",
				Foobar: map[string]string{
					"Key": "val",
				},
			},
			template: `{{ toYaml .Values }}`,
			expected: "Bar: foo\nFoo: bar\nFoobar:\n  Key: val",
		},
		{
			name:     "static",
			template: "static string",
			expected: "static string",
		},
	}

	os.Setenv("FOO", "BAR")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(strings.Builder)
			err := Render(context.TODO(), buf, tt.template, tt.data)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if buf.String() != tt.expected {
				t.Errorf("unexpected result\n\twanted: %s\n\tgot: %s", tt.expected, buf.String())
			}
		})
	}
}
