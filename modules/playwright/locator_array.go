package playwright

import (
	"context"
	"fmt"

	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
)

// LocatorArray wraps a slice of playwright Locators
type LocatorArray struct {
	locators []*Locator
}

func (la *LocatorArray) Type() object.Type {
	return "playwright.locator_array"
}

func (la *LocatorArray) Inspect() string {
	return fmt.Sprintf("playwright.locator_array(len=%d)", len(la.locators))
}

func (la *LocatorArray) Interface() interface{} {
	return la.locators
}

func (la *LocatorArray) Equals(other object.Object) object.Object {
	return object.NewBool(la == other)
}

func (la *LocatorArray) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (la *LocatorArray) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("cannot set attribute %q on playwright.locator_array object", name))
}

func (la *LocatorArray) IsTruthy() bool {
	return len(la.locators) > 0
}

func (la *LocatorArray) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("operation %v not supported on playwright.locator_array object", opType))
}

func (la *LocatorArray) Cost() int {
	return 0
}

func (la *LocatorArray) GetItem(key object.Object) (object.Object, *object.Error) {
	idx, err := object.AsInt(key)
	if err != nil {
		return nil, err
	}
	if idx < 0 || int(idx) >= len(la.locators) {
		return nil, object.NewError(fmt.Errorf("index out of range: %d", idx))
	}
	return la.locators[idx], nil
}

func (la *LocatorArray) GetSlice(s object.Slice) (object.Object, *object.Error) {
	start := 0
	end := len(la.locators)

	if s.Start != nil {
		startIdx, err := object.AsInt(s.Start)
		if err != nil {
			return nil, err
		}
		start = int(startIdx)
		if start < 0 {
			start = 0
		}
	}

	if s.Stop != nil {
		stopIdx, err := object.AsInt(s.Stop)
		if err != nil {
			return nil, err
		}
		end = int(stopIdx)
		if end > len(la.locators) {
			end = len(la.locators)
		}
	}

	if start > end {
		start = end
	}

	result := &LocatorArray{
		locators: la.locators[start:end],
	}
	return result, nil
}

func (la *LocatorArray) SetItem(key, value object.Object) *object.Error {
	return object.NewError(fmt.Errorf("cannot modify playwright.locator_array object"))
}

func (la *LocatorArray) DelItem(key object.Object) *object.Error {
	return object.NewError(fmt.Errorf("cannot modify playwright.locator_array object"))
}

func (la *LocatorArray) Contains(item object.Object) *object.Bool {
	for _, locator := range la.locators {
		if locator.Equals(item) == object.True {
			return object.True
		}
	}
	return object.False
}

func (la *LocatorArray) Len() *object.Int {
	return object.NewInt(int64(len(la.locators)))
}

func (la *LocatorArray) Iter() object.Iterator {
	return &LocatorArrayIterator{array: la, index: -1}
}

type LocatorArrayIterator struct {
	array *LocatorArray
	index int
	entry *LocatorArrayIteratorEntry
}

func (it *LocatorArrayIterator) Type() object.Type {
	return "playwright.locator_array_iterator"
}

func (it *LocatorArrayIterator) Inspect() string {
	return fmt.Sprintf("playwright.locator_array_iterator(index=%d)", it.index)
}

func (it *LocatorArrayIterator) Interface() interface{} {
	return it
}

func (it *LocatorArrayIterator) Equals(other object.Object) object.Object {
	return object.NewBool(it == other)
}

func (it *LocatorArrayIterator) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (it *LocatorArrayIterator) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("cannot set attribute %q on playwright.locator_array_iterator object", name))
}

func (it *LocatorArrayIterator) IsTruthy() bool {
	return true
}

func (it *LocatorArrayIterator) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("operation %v not supported on playwright.locator_array_iterator object", opType))
}

func (it *LocatorArrayIterator) Cost() int {
	return 0
}

func (it *LocatorArrayIterator) Next(ctx context.Context) (object.Object, bool) {
	it.index++
	if it.index >= len(it.array.locators) {
		return nil, false
	}
	it.entry = &LocatorArrayIteratorEntry{
		idx:     it.index,
		locator: it.array.locators[it.index],
	}
	return it.entry.Value(), true
}

func (it *LocatorArrayIterator) Entry() (object.IteratorEntry, bool) {
	if it.entry == nil {
		return nil, false
	}
	return it.entry, true
}

type LocatorArrayIteratorEntry struct {
	idx     int
	locator *Locator
}

func (e *LocatorArrayIteratorEntry) Type() object.Type {
	return "playwright.locator_array_iterator_entry"
}

func (e *LocatorArrayIteratorEntry) Inspect() string {
	return fmt.Sprintf("playwright.locator_array_iterator_entry(idx=%d)", e.idx)
}

func (e *LocatorArrayIteratorEntry) Interface() interface{} {
	return e
}

func (e *LocatorArrayIteratorEntry) Equals(other object.Object) object.Object {
	return object.NewBool(e == other)
}

func (e *LocatorArrayIteratorEntry) GetAttr(name string) (object.Object, bool) {
	return nil, false
}

func (e *LocatorArrayIteratorEntry) SetAttr(name string, value object.Object) error {
	return object.NewError(fmt.Errorf("cannot set attribute %q on playwright.locator_array_iterator_entry object", name))
}

func (e *LocatorArrayIteratorEntry) IsTruthy() bool {
	return true
}

func (e *LocatorArrayIteratorEntry) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.NewError(fmt.Errorf("operation %v not supported on playwright.locator_array_iterator_entry object", opType))
}

func (e *LocatorArrayIteratorEntry) Cost() int {
	return 0
}

func (e *LocatorArrayIteratorEntry) Key() object.Object {
	return object.NewInt(int64(e.idx))
}

func (e *LocatorArrayIteratorEntry) Value() object.Object {
	return e.locator
}

func (e *LocatorArrayIteratorEntry) Primary() object.Object {
	return e.Value()
}
