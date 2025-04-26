package slack

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

// NextFn is a function type that produces the next item in an iteration
type NextFn func(ctx context.Context) (value object.Object, hasNext bool, err error)

// GenericIterator implements a Risor iterator that can be used for any paginated results
type GenericIterator struct {
	iterType     string
	nextFn       NextFn
	currentIndex int64
	currentEntry *object.Entry
}

// NewGenericIterator creates a new generic iterator with the given type and next function
func NewGenericIterator(iterType string, nextFn NextFn) *GenericIterator {
	return &GenericIterator{
		iterType:     iterType,
		nextFn:       nextFn,
		currentIndex: -1,
	}
}

func (i *GenericIterator) Type() object.Type {
	return object.Type(i.iterType)
}

func (i *GenericIterator) Inspect() string {
	return fmt.Sprintf("%s()", i.iterType)
}

func (i *GenericIterator) Interface() interface{} {
	return i
}

func (i *GenericIterator) Equals(other object.Object) object.Object {
	if i == other {
		return object.True
	}
	return object.False
}

func (i *GenericIterator) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (i *GenericIterator) SetAttr(name string, value object.Object) error {
	return fmt.Errorf("cannot set attribute on %s", i.iterType)
}

func (i *GenericIterator) IsTruthy() bool {
	return true
}

func (i *GenericIterator) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.Errorf("operation not supported on %s", i.iterType)
}

func (i *GenericIterator) Cost() int {
	return 1
}

// Next method returns the next value in the iteration, along with a boolean
// indicating whether there are more items to iterate over.
func (i *GenericIterator) Next(ctx context.Context) (object.Object, bool) {
	// Call the next function to get the next value
	value, hasNext, err := i.nextFn(ctx)

	// If there was an error, return it
	if err != nil {
		return object.NewError(err), false
	}

	// If there are no more items, return nil and false to terminate iteration
	if !hasNext {
		return nil, false
	}

	// Update the current index and entry
	i.currentIndex++
	i.currentEntry = object.NewEntry(object.NewInt(i.currentIndex), value)

	return value, true
}

func (i *GenericIterator) Entry() (object.IteratorEntry, bool) {
	if i.currentEntry == nil {
		return nil, false
	}
	return i.currentEntry, true
}
