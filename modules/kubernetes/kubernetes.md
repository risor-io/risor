import { Callout } from 'nextra/components';

# kubernetes

<Callout type="info" emoji="ℹ️">
  This module requires that Risor has been compiled with the `k8s` Go build tag.
  When compiling **manually**, [make sure you specify `-tags k8s`](https://github.com/risor-io/risor#build-and-install-the-cli-from-source).
</Callout>

Module `k8s` provides methods for getting, listing, deleting and updating resources using the Kubernetes API.

## Functions

### get

```go filename="Function signature"
get(kind string, options object) object
```

Can be used to get a single object or a list of objects from the Kubernetes API.

```go filename="Example"
# get all namespaces
k8s.get("namespace.v1")

# get a specific namespace
k8s.get("namespace.v1", {"name": "default"})

# get pods by labels
k8s.get("pod.v1", {"namespace": "kube-system", "selector": "k8s-app=kube-apiserver" })

# get a specific pod by name
k8s.get("pod.v1", {"namespace": "kube-system", "name": "kube-apiserver-i-036c4abb51cd79a10"})

# get pods not running
k8s.get("pod.v1", {"fieldSelector": "status.phase!=Running"})
```

### delete

```go filename="Function signature"
delete(kind string, options object) object
```

Can be used to delete a single object or a list of objects from the Kubernetes API.

```go filename="Example"
# delete pods not running
k8s.delete("pod.v1", {"fieldSelector": "status.phase!=Running"})
```

### apply

```go filename="Function signature"
apply(manifest string, options object)
```

Can be used to apply (create or update) a kubernetes object from a JSON or YAML manifest

```go filename="Example"
manifest := string(os.read_file("/tmp/foo.yaml"))
apply(manifest, {"namespace": "my-namespace"})
```
