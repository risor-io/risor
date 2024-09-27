package object

import "github.com/risor-io/risor/errz"

type base struct{}

func (b *base) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (b *base) SetAttr(name string, value Object) error {
	return errz.TypeErrorf("type error: object has no attribute %q", name)
}

func (b *base) IsTruthy() bool {
	return true
}

func (b *base) Cost() int {
	return 0
}
