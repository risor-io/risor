package yaml

import (
	"context"

	"github.com/risor-io/risor/builtins"
	"github.com/risor-io/risor/object"
	"gopkg.in/yaml.v3"
)

func init() {
	builtins.RegisterCodec("yaml", &builtins.Codec{
		Encode: encodeYAML,
		Decode: decodeYAML,
	})
}

func encodeYAML(ctx context.Context, obj object.Object) object.Object {
	nativeObject := obj.Interface()
	if nativeObject == nil {
		return object.Errorf("value error: encode() does not support %T", obj)
	}
	yamlBytes, err := yaml.Marshal(nativeObject)
	if err != nil {
		return object.NewError(err)
	}
	return object.NewString(string(yamlBytes))
}

func decodeYAML(ctx context.Context, obj object.Object) object.Object {
	data, err := object.AsBytes(obj)
	if err != nil {
		return err
	}
	var result interface{}
	if err := yaml.Unmarshal([]byte(data), &result); err != nil {
		return object.NewError(err)
	}
	return object.FromGoType(result)
}
