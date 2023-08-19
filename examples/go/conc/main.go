package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/object"
)

type State struct {
	Name    string
	Running bool
}

type Service struct {
	name       string
	running    bool
	startCount int
	stopCount  int
}

func (s *State) IsRunning() bool {
	return s.Running
}

func (s *Service) Start() error {
	if s.running {
		return fmt.Errorf("service %s already running", s.name)
	}
	s.running = true
	s.startCount++
	return nil
}

func (s *Service) Stop() error {
	if !s.running {
		return fmt.Errorf("service %s not running", s.name)
	}
	s.running = false
	s.stopCount++
	return nil
}

func (s *Service) SetName(name string) {
	s.name = name
}

func (s *Service) GetName() string {
	return s.name
}

func (s *Service) PrintState() {
	fmt.Printf("printing state... name: %s running %t\n", s.name, s.running)
}

func (s *Service) GetState() *State {
	return &State{
		Name:    s.name,
		Running: s.running,
	}
}

const defaultExample = `
runFunc('svc-{i}', func(svc) {
	i := i
	svc.SetName('svc-{i}')
	state := svc.GetState()
	print("STATE:", state, type(state))
	state.IsRunning()
})
`

var red = color.New(color.FgRed).SprintfFunc()

func GetServiceRunFunc(svc *Service) object.BuiltinFunction {
	return func(ctx context.Context, args ...object.Object) object.Object {
		if len(args) != 2 {
			return object.Errorf("expected 2 arguments (got %d)", len(args))
		}
		strObj, ok := args[0].(*object.String)
		if !ok {
			return object.Errorf("expected a string")
		}
		name := strObj.String()
		fn, ok := args[1].(*object.Function)
		if !ok {
			return object.Errorf("expected a function")
		}
		proxy, err := object.NewProxy(svc)
		if err != nil {
			return object.NewError(err)
		}
		callFunc, ok := object.GetCallFunc(ctx)
		if !ok {
			object.Errorf("no call func")
		}
		result, err := callFunc(ctx, fn, []object.Object{proxy})
		if err != nil {
			return object.Errorf("failed to run func for %q: %s", name, err)
		}
		return result
	}
}

func main() {
	var code string
	flag.StringVar(&code, "code", defaultExample, "Code to evaluate")
	flag.Parse()

	ctx := context.Background()

	svc := &Service{}

	runFunc := object.NewBuiltin("runFunc", GetServiceRunFunc(svc))

	taskCount := 10
	var wg sync.WaitGroup
	wg.Add(taskCount)

	for i := 0; i < taskCount; i++ {
		go func(i int) {
			defer wg.Done()
			result, err := risor.Eval(ctx, code, risor.WithGlobal("i", i), risor.WithGlobal("runFunc", runFunc))
			if err != nil {
				fmt.Println(red(err.Error()))
				os.Exit(1)
			}
			fmt.Println("RESULT:", result)
		}(i)
	}
	wg.Wait()
}
