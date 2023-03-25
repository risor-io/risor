package object

import "fmt"

type Cell struct {
	*DefaultImpl
	value *Object
}

func (c *Cell) Value() Object {
	if c.value == nil {
		return nil
	}
	return *c.value
}

func (c *Cell) Set(value Object) {
	*c.value = value
}

func (c *Cell) Type() Type {
	return CELL
}

func (c *Cell) Interface() interface{} {
	if c.value == nil {
		return nil
	}
	return (*c.value).Interface()
}

func (c *Cell) Equals(other Object) Object {
	if c == other {
		return True
	}
	return False
}

func NewCell(value *Object) *Cell {
	fmt.Println("NewCell", value)
	return &Cell{DefaultImpl: &DefaultImpl{}, value: value}
}
