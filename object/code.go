package object

import (
	"github.com/cloudcmds/tamarin/internal/op"
)

type Loop struct {
	ContinuePos []int
	BreakPos    []int
}

type Code struct {
	Name         string
	IsNamed      bool
	Parent       *Code
	Symbols      *SymbolTable
	Instructions []op.Code
	Constants    []Object
	Loops        []*Loop
	Names        []string
	Source       string
}

func (c *Code) AddName(name string) uint16 {
	c.Names = append(c.Names, name)
	return uint16(len(c.Names) - 1)
}

func (c *Code) SymbolCount() uint16 {
	return c.Symbols.Size()
}

func NewCode(name string) *Code {
	return &Code{Name: name, Symbols: NewSymbolTable()}
}
