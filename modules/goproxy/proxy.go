package goproxy

import (
	"fmt"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

var _ object.Object = (*Proxy)(nil)

const PROXY object.Type = "goproxy.proxy"

type Proxy struct {
	value *goproxy.ProxyHttpServer
}

func (p *Proxy) IsTruthy() bool {
	return true
}

func (p *Proxy) Type() object.Type {
	return PROXY
}

func (p *Proxy) Inspect() string {
	return fmt.Sprintf("%s()", PROXY)
}

func (p *Proxy) Value() *goproxy.ProxyHttpServer {
	return p.value
}

func (p *Proxy) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("attribute error: cannot set %q on %s object", name, PROXY)
}

func (p *Proxy) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "sprintf":
		return object.NewString("OK"), true
	default:
		return nil, false
	}
}

func (p *Proxy) Interface() interface{} {
	return p.value
}

func (p *Proxy) Equals(other object.Object) object.Object {
	return object.NewBool(p == other)
}

func (p *Proxy) Cost() int {
	return 0
}

func (p *Proxy) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("eval error: unsupported operation for %s: %v", PROXY, opType)
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.value.ServeHTTP(w, r)
}

func NewProxy(v *goproxy.ProxyHttpServer) *Proxy {
	return &Proxy{value: v}
}
