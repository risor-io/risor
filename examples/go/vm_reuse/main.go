package main

import (
	"context"
	"fmt"
	"log"

	"github.com/risor-io/risor"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/vm"
)

func getCustomModule() *object.Module {
	return object.NewBuiltinsModule(
		"simplemath",
		map[string]object.Object{
			"add": object.NewBuiltin("add", func(ctx context.Context, args ...object.Object) object.Object {
				if len(args) != 2 {
					return object.Errorf("add takes 2 arguments")
				}
				a, err := object.AsInt(args[0])
				if err != nil {
					return object.Errorf("add expected an integer, got %s", args[0].Type())
				}
				b, err := object.AsInt(args[1])
				if err != nil {
					return object.Errorf("add expected an integer, got %s", args[1].Type())
				}
				return object.NewInt(a + b)
			}),
		},
	)
}

func main() {
	ctx := context.Background()
	customModule := getCustomModule()

	vm, err := vm.NewEmpty()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		fmt.Printf("==== execution %d ====\n", i)

		result, err := risor.Eval(ctx,
			"simplemath.add(1, 2)",
			risor.WithVM(vm),
			risor.WithGlobals(map[string]any{"simplemath": customModule}),
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
}
