//go:build aws
// +build aws

package aws

import (
	"fmt"

	"github.com/risor-io/risor/object"
)

func mapGetStr(m *object.Map, key string) (string, bool, error) {
	value := m.GetWithDefault(key, nil)
	if value == nil {
		return "", false, nil
	}
	str, err := object.AsString(value)
	if err != nil {
		return "", true, fmt.Errorf("type error: %s must be a string (got %s)", key, value.Type())
	}
	return str, true, nil
}

func mapGetMap(m *object.Map, key string) (*object.Map, bool, error) {
	value := m.GetWithDefault(key, nil)
	if value == nil {
		return nil, false, nil
	}
	valueMap, err := object.AsMap(value)
	if err != nil {
		return nil, true, fmt.Errorf("type error: %s must be a map (got %s)", key, value.Type())
	}
	return valueMap, true, nil
}

func mapGetStrList(m *object.Map, key string) ([]string, bool, error) {
	value := m.GetWithDefault(key, nil)
	if value == nil {
		return nil, false, nil
	}
	list, err := object.AsList(value)
	if err != nil {
		return nil, true, fmt.Errorf("type error: %s must be a list (got %s)", key, value.Type())
	}
	result := make([]string, list.Size())
	for i, item := range list.Value() {
		str, err := object.AsString(item)
		if err != nil {
			return nil, true, fmt.Errorf("type error: %s[%d] must be a string (got %s)", key, i, item.Type())
		}
		result[i] = str
	}
	return result, true, nil
}
