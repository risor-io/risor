package yaml

import (
	"fmt"

	"github.com/risor-io/risor/object"
	"gopkg.in/yaml.v3"
)

//risor:generate

//risor:export
func unmarshal(data []byte) object.Object {
	var obj interface{}
	if err := yaml.Unmarshal(data, &obj); err != nil {
		return object.Errorf("value error: yaml.unmarshal failed with: %s", err.Error())
	}
	scriptObj := object.FromGoType(obj)
	if scriptObj == nil {
		return object.Errorf("type error: yaml.unmarshal failed")
	}
	return scriptObj
}

//risor:export
func marshal(value object.Object) (string, error) {
	b, err := yaml.Marshal(value.Interface())
	if err != nil {
		return "", fmt.Errorf("value error: yaml.marshal failed: %w", err)
	}
	return string(b), nil
}

//risor:export
func valid(data []byte) bool {
	var v any
	return yaml.Unmarshal(data, &v) == nil
}
