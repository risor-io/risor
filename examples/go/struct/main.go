package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/object"
)

type Counts struct {
	Start int
	Stop  int
}

type State struct {
	Name    string
	Running bool
	Counts  Counts
}

type Service struct {
	name       string
	running    bool
	startCount int
	stopCount  int
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

func (s *Service) GetState() State {
	return State{
		Name:    s.name,
		Running: s.running,
		Counts: Counts{
			Start: s.startCount,
			Stop:  s.stopCount,
		},
	}
}

func (s *Service) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"running":     s.running,
		"start_count": s.startCount,
		"stop_count":  s.stopCount,
	}
}

const defaultExample = `
svc.SetName("My Service")
svc.Start()
svc.PrintState()
print("metrics:", svc.GetMetrics())
print("start count:", svc.GetState()["Counts"]["Start"])
`

var red = color.New(color.FgRed).SprintfFunc()

func main() {
	var code string
	flag.StringVar(&code, "code", defaultExample, "Code to evaluate")
	flag.Parse()

	ctx := context.Background()

	// Initialize the service
	svc := &Service{}

	// Create a Risor proxy for the service
	proxy, err := object.NewProxy(svc)
	if err != nil {
		fmt.Println(red(err.Error()))
		os.Exit(1)
	}

	// Build up options for Risor, including the proxy as a variable named "svc"
	opts := []risor.Option{
		risor.WithDefaultBuiltins(),
		risor.WithBuiltins(map[string]object.Object{"svc": proxy}),
	}

	// Run the Risor code which can access the service as `svc`
	if _, err = risor.Eval(ctx, code, opts...); err != nil {
		fmt.Println(red(err.Error()))
		os.Exit(1)
	}

	fmt.Println(svc.GetMetrics())
}
