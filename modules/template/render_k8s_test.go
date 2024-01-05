//go:build k8s
// +build k8s

package template

import (
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestRenderK8s(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     any
		expected string
	}{
		{
			name: "lookup ConfigMap",
			template: `
{{- $cm := k8sLookup "ConfigMap.v1" "default" "my-config-map" -}}
{{- $cfg := index $cm.data "conf.yaml" | fromYaml -}}
key1 is {{ $cm.data.key1 }}
foobar key is {{ $cfg.foobar.key }}`,
			expected: "key1 is val1\nfoobar key is val",
		},
		{
			name:     "ConfigMap json path",
			template: `key is {{ k8sLookup "ConfigMap.v1" "default" "my-config-map" | jsonPath ".data.key1" }}`,
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
	}

	jcm := `{
  "apiVersion": "v1",
  "kind": "ConfigMap",
  "metadata": {
    "name": "my-config-map",
    "namespace": "default"
  },
  "data": {
    "key1": "val1",
    "conf.yaml": "foo: bar\nbar: foo\nfoobar:\n  key: val\n"
  }
}`

	cm, _, err := unstructured.UnstructuredJSONScheme.Decode([]byte(jcm), nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sClient = fake.NewClientBuilder().WithObjects(cm.(*unstructured.Unstructured)).Build()
			buf := new(strings.Builder)
			err := Render(buf, tt.template, tt.data)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if buf.String() != tt.expected {
				t.Errorf("unexpected result\n\twanted: %s\n\tgot: %s", tt.expected, buf.String())
			}
		})
	}
}
