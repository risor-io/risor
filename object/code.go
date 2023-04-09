package object

import (
	"github.com/cloudcmds/tamarin/internal/op"
)

type Loop struct {
	ContinuePos []uint16
	BreakPos    []uint16
}

type Code struct {
	Name         string
	IsNamed      bool
	Parent       *Code
	Children     []*Code
	Symbols      *SymbolTable
	Instructions []op.Code
	Constants    []Object
	Loops        []*Loop
	Names        []string
}

func (s *Code) AddName(name string) uint16 {
	s.Names = append(s.Names, name)
	return uint16(len(s.Names) - 1)
}

func NewCode(name string) *Code {
	return &Code{Name: name, Symbols: NewSymbolTable()}
}
