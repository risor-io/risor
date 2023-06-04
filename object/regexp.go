package object

import (
	"fmt"
	"regexp"

	"github.com/cloudcmds/tamarin/v2/op"
)

// Regexp wraps regexp.Regexp and implements the Object interface.
type Regexp struct {
	value *regexp.Regexp
}

func (r *Regexp) Type() Type {
	return REGEXP
}

func (r *Regexp) Value() *regexp.Regexp {
	return r.value
}

func (r *Regexp) Inspect() string {
	return r.value.String()
}

func (r *Regexp) String() string {
	return fmt.Sprintf("regexp(%s)", r.value.String())
}

func (r *Regexp) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (r *Regexp) Interface() interface{} {
	return r.value
}

func (r *Regexp) Compare(other Object) (int, error) {
	typeComp := CompareTypes(r, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherRegex := other.(*Regexp)
	thisStr := r.value.String()
	otherStr := otherRegex.value.String()
	if thisStr > otherStr {
		return 1, nil
	} else if thisStr < otherStr {
		return -1, nil
	}
	return 0, nil
}

func (r *Regexp) Equals(other Object) Object {
	if other.Type() != REGEXP {
		return False
	}
	otherRegex := other.(*Regexp)
	thisStr := r.value.String()
	otherStr := otherRegex.value.String()
	if thisStr == otherStr {
		return True
	}
	return False
}

func (r *Regexp) IsTruthy() bool {
	return true
}

func (r *Regexp) RunOperation(opType op.BinaryOpType, right Object) Object {
	return NewError(fmt.Errorf("unsupported operation for regexp: %v", opType))
}

func NewRegexp(re *regexp.Regexp) *Regexp {
	return &Regexp{value: re}
}
