## Proxying Calls to Go Objects

You can expose arbitrary Go objects to Tamarin code in order to enable method
calls on those objects. This allows you to expose existing structs in your
application as Tamarin objects that scripts can be written against. Tamarin
automatically discovers public methods on your Go types and converts inputs and
outputs for primitive types and for structs that you register.

Input and output values are type-converted automatically, for a variety of types.
Go structs are mapped to Tamarin map objects. Go `context.Context` and `error`
values are handled automatically.

```go title="proxy_service.go" linenums="1"
	// Create a registry that tracks proxied Go types and their attributes
	registry, err := object.NewTypeRegistry()
	if err != nil {
		return err
	}

	// This is the Go service we will expose in Tamarin
	svc := &MyService{}

	// Wrap the service in a Tamarin Proxy
	proxy, err := object.NewProxy(registry, svc)
	if err != nil {
		return err
	}

	// Add the proxy to a Tamarin execution scope
	s := scope.New(scope.Opts{})
	s.Declare("svc", proxy, true)

	// Execute Tamarin code against that scope. By doing this, the Tamarin
	// code can call public methods on `svc` and retrieve its public fields.
	result, err := exec.Execute(ctx, exec.Opts{
		Input: string(scriptSourceCode),
		Code: s,
	})
```

See [example-proxy](../cmd/example-proxy/main.go) for a complete example.
