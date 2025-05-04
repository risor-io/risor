package main

import (
	"context"
	"log"

	"github.com/risor-io/risor"
	"github.com/solarlune/tetra3d"
)

type Engine struct{}

func (e *Engine) NewVector(x, y, z float32) tetra3d.Vector3 {
	return tetra3d.NewVector3(x, y, z)
}

func main() {
	src := `
	a := Engine.NewVector(4,5,6)
	b := Engine.NewVector(1,2,3)
	c := a.Add(b) // This works now, which is great!
	c.X = 15 // However, this fails with "type error: cannot set field X"
	print("c.X =", c.X) // Print the value to verify it was set correctly
	`
	_, err := risor.Eval(context.Background(), src, risor.WithGlobal("Engine", &Engine{}))
	if err != nil {
		log.Fatal(err)
	}
}
