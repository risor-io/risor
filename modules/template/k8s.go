//go:build k8s
// +build k8s

package template

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var k8sClient client.Reader

func k8sLookup(kind, namespace, name string) (map[string]any, error) {
	if k8sClient == nil {
		kfg, err := config.GetConfig()
		if err != nil {
			return nil, err
		}
		if c, err := client.New(kfg, client.Options{}); err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		} else {
			k8sClient = c
		}
	}
	gvk, gk := schema.ParseKindArg(kind)
	if gvk == nil {
		// this looks strange but it should make sense if you read the ParseKindArg docs
		gvk = &schema.GroupVersionKind{
			Kind:    gk.Kind,
			Version: gk.Group,
		}
	}
	if name != "" {
		// fetching a single resource by name
		u := unstructured.Unstructured{}
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   gvk.Group,
			Kind:    gvk.Kind,
			Version: gvk.Version,
		})
		if err := k8sClient.Get(context.Background(), client.ObjectKey{
			Namespace: namespace,
			Name:      name,
		}, &u); err != nil {
			return nil, fmt.Errorf("failed to get: %w", err)
		}
		return u.UnstructuredContent(), nil
	}
	ul := &unstructured.UnstructuredList{}
	ul.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind + "List", // TODO: is there a better way?
	})
	opts := &client.ListOptions{
		Namespace: namespace,
	}
	if err := k8sClient.List(context.Background(), ul, opts); err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return ul.UnstructuredContent(), nil
}
