package main

import (
	"context"
	"fmt"

	"github.com/risor-io/risor"
)

type Object struct {
	A int
}

func (o Object) Set(other Object) Object {
	o.A = other.A
	fmt.Println("set", o.A)
	return o
}

type Creator struct{}

func (c *Creator) NewObject() Object {
	return Object{}
}

func main() {

	src := `
	obj := Creator.NewObject()
	obj2 := Creator.NewObject()
	obj3 := obj.Set(obj2)
	`

	_, err := risor.Eval(context.Background(), src, risor.WithGlobal("Creator", &Creator{}))
	if err != nil {
		fmt.Println(err)
	}

}
