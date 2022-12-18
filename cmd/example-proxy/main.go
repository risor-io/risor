// This program demonstrates using Tamarin as a library
// to run a simple script.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cloudcmds/tamarin/exec"
	"github.com/cloudcmds/tamarin/object"
	"github.com/cloudcmds/tamarin/scope"
)

type MyServiceOpts struct {
	Foo string
	Bar int
}

type MyService struct {
	Name    string
	Age     int
	IsAdult bool
	Blergh  []string
}

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

func (svc *MyService) RunCtx(ctx context.Context, opts MyServiceOpts) (string, error) {
	select {
	case <-ctx.Done():
		return "", errors.New("deadline exceeded")
	case <-time.After(100 * time.Millisecond):
	}
	return svc.Run(opts)
}

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	svc := &MyService{
		Name:    "This is a Test",
		Age:     21,
		IsAdult: true,
		// Blergh:  []string{"foo", "bar"},
	}

	input := `
	// First this should print HELLO
	print("svc.ToUpper() result:", svc.ToUpper("hello"))

	// Then ParseInt should return an *object.Result since it has an error
	// in its return signature
	result := svc.ParseInt("234")
	assert(result.is_ok())
	print("svc.ParseInt() result:", result.unwrap())

	// Call a method that accepts a complex type
	result = svc.RunCtx({"foo": "fish", "bar": 42})
	print("svc.RunCtx() result:", result)

	print("svc.Name:", svc.Name)
	print("svc.Age:", svc.Age)
	print("svc.IsAdult:", svc.IsAdult)
	print("svc.Blergh:", svc.Blergh)

	svc.Name
	`

	registry, err := object.NewTypeRegistry()
	if err != nil {
		fmt.Println("Type registry error:", err)
		os.Exit(1)
	}

	p, err := object.NewProxy(registry, svc)
	if err != nil {
		fmt.Println("Proxy error:", err)
		os.Exit(1)
	}

	s := scope.New(scope.Opts{})
	s.Declare("svc", p, true)

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
