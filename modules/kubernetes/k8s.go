//go:build k8s
// +build k8s

package kubernetes

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	konfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"

	"github.com/risor-io/risor/object"
)

var k8sClient client.Client

func buildK8sClient(kubeconfig, context string) (client.Client, error) {
	var (
		kfg *rest.Config
		err error
	)
	if kubeconfig != "" || context != "" {
		kfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
			&clientcmd.ConfigOverrides{
				CurrentContext: context,
			},
		).ClientConfig()
	} else {
		kfg, err = konfig.GetConfig()
	}
	if err != nil {
		return nil, err
	}
	return client.New(kfg, client.Options{})
}

func getStrOption(params *object.Map, name string) (string, *object.Error) {
	var (
		res    string
		errObj *object.Error
	)
	if params != nil {
		if o := params.GetWithDefault(name, nil); o != nil {
			res, errObj = object.AsString(o)
			if errObj != nil {
				return res, errObj
			}
		}
	}

	return res, nil
}

func getBoolOption(params *object.Map, name string) (bool, *object.Error) {
	var (
		res    bool
		errObj *object.Error
	)
	if params != nil {
		if o := params.GetWithDefault(name, nil); o != nil {
			res, errObj = object.AsBool(o)
			if errObj != nil {
				return res, errObj
			}
		}
	}

	return res, nil
}

func getLabelSelector(params *object.Map) (labels.Selector, *object.Error) {
	selectorStr, errObj := getStrOption(params, "selector")
	if errObj != nil {
		return nil, errObj
	}
	if selectorStr == "" {
		return nil, nil
	}
	sel, err := labels.Parse(selectorStr)
	if err != nil {
		return nil, object.NewError(err)
	}

	return sel, nil
}

func getFieldSelector(params *object.Map) (fields.Selector, *object.Error) {
	selectorStr, errObj := getStrOption(params, "fieldSelector")
	if errObj != nil {
		return nil, errObj
	}
	if selectorStr == "" {
		return nil, nil
	}
	sel, err := fields.ParseSelector(selectorStr)
	if err != nil {
		return nil, object.NewError(err)
	}

	return sel, nil
}

func getObjectMeta(kind string, params *object.Map) (*schema.GroupVersionKind, client.ObjectKey, *object.Error) {
	name, errObj := getStrOption(params, "name")
	if errObj != nil {
		return nil, client.ObjectKey{}, errObj
	}
	namespace, errObj := getStrOption(params, "namespace")
	if errObj != nil {
		return nil, client.ObjectKey{}, errObj
	}
	if namespace == "all" {
		namespace = ""
	}
	if k8sClient == nil {
		if c, err := buildK8sClient("", ""); err != nil {
			return nil, client.ObjectKey{}, object.NewError(err)
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

	return gvk, client.ObjectKey{Namespace: namespace, Name: name}, nil
}

func Get(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsRangeError("k8s.get", 1, 2, len(args))
	}
	kind, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var (
		errObj *object.Error
		params *object.Map
	)
	if numArgs == 2 {
		params, errObj = object.AsMap(args[1])
		if errObj != nil {
			return errObj
		}
	}
	gvk, key, err := getObjectMeta(kind, params)
	if err != nil {
		return err
	}
	if key.Name != "" {
		// fetching a single resource by name
		u := unstructured.Unstructured{}
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   gvk.Group,
			Kind:    gvk.Kind,
			Version: gvk.Version,
		})
		if err := k8sClient.Get(context.Background(), key, &u); err != nil {
			return object.NewError(err)
		}
		return object.FromGoType(u.UnstructuredContent())
	}
	ul := &unstructured.UnstructuredList{}
	ul.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind + "List", // TODO: is there a better way?
	})
	opts := &client.ListOptions{
		Namespace: key.Namespace,
	}
	if sel, err := getLabelSelector(params); err != nil {
		return err
	} else {
		opts.LabelSelector = sel
	}
	if sel, err := getFieldSelector(params); err != nil {
		return err
	} else {
		opts.FieldSelector = sel
	}
	if err := k8sClient.List(context.Background(), ul, opts); err != nil {
		return object.NewError(err)
	}
	return object.FromGoType(ul.UnstructuredContent())
}

func Delete(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsRangeError("k8s.delete", 1, 2, len(args))
	}
	kind, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	var (
		errObj *object.Error
		params *object.Map
	)
	if numArgs == 2 {
		params, errObj = object.AsMap(args[1])
		if errObj != nil {
			return errObj
		}
	}
	gvk, key, err := getObjectMeta(kind, params)
	if err != nil {
		return err
	}
	if key.Name != "" {
		// deleting a single resource by name
		u := unstructured.Unstructured{}
		u.SetName(key.Name)
		u.SetNamespace(key.Namespace)
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Group:   gvk.Group,
			Kind:    gvk.Kind,
			Version: gvk.Version,
		})
		if err := k8sClient.Delete(context.Background(), &u); err != nil {
			return object.NewError(err)
		}
		return object.FromGoType("deleted")
	}
	ul := &unstructured.UnstructuredList{}
	ul.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind + "List", // TODO: is there a better way?
	})
	opts := &client.ListOptions{
		Namespace: key.Namespace,
	}
	if sel, err := getLabelSelector(params); err != nil {
		return err
	} else {
		opts.LabelSelector = sel
	}
	if sel, err := getFieldSelector(params); err != nil {
		return err
	} else {
		opts.FieldSelector = sel
	}
	deleteAll, err := getBoolOption(params, "all")
	if err != nil {
		return err
	}
	if !deleteAll && opts.LabelSelector == nil && opts.FieldSelector == nil {
		return object.NewError(errors.New(`to delete all resources you must explcitly set the "all" option`))
	}
	if err := k8sClient.List(context.Background(), ul, opts); err != nil {
		return object.NewError(err)
	}
	for _, obj := range ul.Items {
		fmt.Println(obj.GetName())
		if err := k8sClient.Delete(context.Background(), &obj); err != nil {
			return object.NewError(err)
		}
	}
	return object.FromGoType("deleted")
}

func y2u(spec string) (*unstructured.Unstructured, error) {
	j, err := yaml.YAMLToJSON([]byte(spec))
	if err != nil {
		return nil, err
	}
	u, _, err := unstructured.UnstructuredJSONScheme.Decode(j, nil, nil)
	if err != nil {
		return nil, err
	}
	return u.(*unstructured.Unstructured), nil
}

func Apply(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 2 {
		return object.NewArgsRangeError("k8s.apply", 1, 2, len(args))
	}

	manifest, objErr := object.AsString(args[0])
	if objErr != nil {
		return objErr
	}
	var (
		errObj *object.Error
		params *object.Map
	)
	if numArgs == 2 {
		params, errObj = object.AsMap(args[1])
		if errObj != nil {
			return errObj
		}
	}
	var namespace string
	namespace, errObj = getStrOption(params, "namespace")
	if errObj != nil {
		return errObj
	}

	u, err := y2u(manifest)
	if err != nil {
		return object.NewError(err)
	}

	if namespace != "" {
		u.SetNamespace(namespace)
	}

	if k8sClient == nil {
		if c, err := buildK8sClient("", ""); err != nil {
			return object.NewError(err)
		} else {
			k8sClient = c
		}
	}

	k8sObj := &unstructured.Unstructured{}
	k8sObj.SetGroupVersionKind(u.GroupVersionKind())
	k8sObj.SetName(u.GetName())
	k8sObj.SetNamespace(u.GetNamespace())
	op, err := controllerutil.CreateOrUpdate(ctx, k8sClient, k8sObj, func() error {
		k8sObj.SetAnnotations(u.GetAnnotations())
		k8sObj.SetLabels(u.GetLabels())
		// Take spec and data from the source
		// TODO: is this list complete?
		for _, key := range []string{"spec", "data", "stringData", "type"} {
			srcSpec, srcSpecOK := u.Object[key]
			if srcSpecOK {
				k8sObj.Object[key] = srcSpec
			} else {
				delete(k8sObj.Object, key)
			}
		}
		return nil
	})
	if err != nil {
		return object.NewError(err)
	}

	return object.FromGoType(op)
}

// TODO: some ideas for other useful functions
// - wait for condition
// - check if ready
func Module() *object.Module {
	return object.NewBuiltinsModule("k8a", map[string]object.Object{
		"get":    object.NewBuiltin("k8s.get", Get),
		"apply":  object.NewBuiltin("k8s.apply", Apply),
		"delete": object.NewBuiltin("k8s.delete", Delete),
	})
}
