//go:build !k8s
// +build !k8s

package template

func k8sLookup(kind, namespace, name string) (map[string]any, error) {
	return nil, nil
}
