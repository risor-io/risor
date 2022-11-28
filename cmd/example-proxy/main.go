// This program demonstrates using Tamarin as a library
// to run a simple script.
package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

type MyServiceOpts struct {
	Foo string
	Bar int
}

type MyService struct{}

func (svc *MyService) ToUpper(s string) string {
	return strings.ToUpper(s)
}

func (svc *MyService) ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func (svc *MyService) Run(opts MyServiceOpts) (string, error) {
	// Could run some I/O or computation here
	return fmt.Sprintf("foo=%s, bar=%d", opts.Foo, opts.Bar), nil
}

func main() {

	ctx := context.Background()

	svc := &MyService{}

	input := `
	// First this should print HELLO
	print(svc.ToUpper("hello"))

	// Then ParseInt should return an *object.Result since it has an error
	// in its return signature
	result := svc.ParseInt("234")
	assert(result.is_ok())
	print("Result was Ok!", result.unwrap())
	
	// Call a method that accepts a complex type
	svc.Run({"foo": "fish", "bar": 42})
	`

	proxyMgr, err := object.NewProxyManager(object.ProxyManagerOpts{
		Types: []any{
			&MyService{},
			MyServiceOpts{},
		},
	})
	if err != nil {
		fmt.Println("Proxy manager error:", err)
		os.Exit(1)
	}

	s := scope.New(scope.Opts{})
	s.Declare("svc", object.NewProxy(proxyMgr, svc), true)

	result, err := exec.Execute(ctx, exec.Opts{
		Input: string(input),
		Scope: s,
	})
	if err != nil {
		fmt.Println("Execution error:", err)
		os.Exit(1)
	}

	fmt.Printf("script result: %s (type %v)\n",
		result.Inspect(), reflect.TypeOf(result))
}
