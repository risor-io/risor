package compiler

import (
	"github.com/cloudcmds/tamarin/internal/op"
	"github.com/cloudcmds/tamarin/internal/symbol"
	"github.com/cloudcmds/tamarin/object"
)

type Scope struct {
	Name         string
	IsNamed      bool
	Parent       *Scope
	Children     []*Scope
	Symbols      *symbol.Table
	Instructions []op.Code
	Constants    []object.Object
	Loops        []*Loop
	Names        []string
}

func (s *Scope) AddName(name string) uint16 {
	s.Names = append(s.Names, name)
	return uint16(len(s.Names) - 1)
}

func NewScope(name string) *Scope {
	return &Scope{Name: name, Symbols: symbol.NewTable()}
}
