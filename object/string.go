package object

import (
	"fmt"
	"strings"
)

// String wraps string and implements Object and Hashable interfaces.
type String struct {
	// Value holds the string wrapped by this object.
	Value string
}

func (s *String) Type() Type {
	return STRING
}

func (s *String) Inspect() string {
	sLen := len(s.Value)
	if sLen >= 2 {
		if s.Value[0] == '"' && s.Value[sLen-1] == '"' {
			if strings.Count(s.Value, "\"") == 2 {
				return fmt.Sprintf("'%s'", s.Value)
			}
		}
	}
	return fmt.Sprintf("%q", s.Value)
}

func (s *String) String() string {
	return fmt.Sprintf("string(%s)", s.Value)
}

func (s *String) HashKey() HashKey {
	return HashKey{Type: s.Type(), StrValue: s.Value}
}

func (s *String) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (s *String) Interface() interface{} {
	return s.Value
}

func (s *String) Compare(other Object) (int, error) {
	typeComp := CompareTypes(s, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherStr := other.(*String)
	if s.Value == otherStr.Value {
		return 0, nil
	}
	if s.Value > otherStr.Value {
		return 1, nil
	}
	return -1, nil
}

func (s *String) Equals(other Object) Object {
	if other.Type() == STRING && s.Value == other.(*String).Value {
		return True
	}
	return False
}

func (s *String) Reversed() *String {
	runes := []rune(s.Value)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return NewString(string(runes))
}

func (s *String) GetItem(key Object) (Object, *Error) {
	indexObj, ok := key.(*Int)
	if !ok {
		return nil, NewError("index error: string index must be an int (got %s)", key.Type())
	}
	runes := []rune(s.Value)
	index, err := ResolveIndex(indexObj.Value, int64(len(runes)))
	if err != nil {
		return nil, NewError(err.Error())
	}
	return NewString(string(runes[index])), nil
}

func (s *String) GetSlice(slice Slice) (Object, *Error) {
	runes := []rune(s.Value)
	start, stop, err := ResolveIntSlice(slice, int64(len(runes)))
	if err != nil {
		return nil, NewError(err.Error())
	}
	resultRunes := runes[start:stop]
	return NewString(string(resultRunes)), nil
}

func (s *String) SetItem(key, value Object) *Error {
	return NewError("eval error: string does not support set item")
}

func (s *String) DelItem(key Object) *Error {
	return NewError("eval error: string does not support del item")
}

func (s *String) Contains(key Object) *Bool {
	strObj, ok := key.(*String)
	if !ok {
		return False
	}
	return NewBool(strings.Contains(s.Value, strObj.Value))
}

func (s *String) Len() *Int {
	return NewInt(int64(len(s.Value)))
}

func NewString(s string) *String {
	return &String{Value: s}
}
