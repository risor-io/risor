package goproxy

import (
	"context"

	"github.com/elazarl/goproxy"
	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func CreateProxy(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("goproxy.create", 0, args); err != nil {
		return err
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	return NewProxy(proxy)
}

func Module() *object.Module {
	return object.NewBuiltinsModule(
		"goproxy",
		map[string]object.Object{},
		CreateProxy)
}
