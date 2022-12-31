package object

import (
	"fmt"
	"strings"
)

type String struct {
	value string
}

func (s *String) Type() Type {
	return STRING
}

func (s *String) Value() string {
	return s.value
}

func (s *String) Inspect() string {
	sLen := len(s.value)
	if sLen >= 2 {
		if s.value[0] == '"' && s.value[sLen-1] == '"' {
			if strings.Count(s.value, "\"") == 2 {
				return fmt.Sprintf("'%s'", s.value)
			}
		}
	}
	return fmt.Sprintf("%q", s.value)
}

func (s *String) String() string {
	return fmt.Sprintf("string(%s)", s.value)
}

func (s *String) HashKey() HashKey {
	return HashKey{Type: s.Type(), StrValue: s.value}
}

func (s *String) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (s *String) Interface() interface{} {
	return s.value
}

func (s *String) Compare(other Object) (int, error) {
	typeComp := CompareTypes(s, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherStr := other.(*String)
	if s.value == otherStr.value {
		return 0, nil
	}
	if s.value > otherStr.value {
		return 1, nil
	}
	return -1, nil
}

func (s *String) Equals(other Object) Object {
	if other.Type() == STRING && s.value == other.(*String).value {
		return True
	}
	return False
}

func (s *String) IsTruthy() bool {
	return s.value != ""
}

func (s *String) Reversed() *String {
	runes := []rune(s.value)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return NewString(string(runes))
}

func (s *String) GetItem(key Object) (Object, *Error) {
	indexObj, ok := key.(*Int)
	if !ok {
		return nil, Errorf("index error: string index must be an int (got %s)", key.Type())
	}
	runes := []rune(s.value)
	index, err := ResolveIndex(indexObj.value, int64(len(runes)))
	if err != nil {
		return nil, Errorf(err.Error())
	}
	return NewString(string(runes[index])), nil
}

func (s *String) GetSlice(slice Slice) (Object, *Error) {
	runes := []rune(s.value)
	start, stop, err := ResolveIntSlice(slice, int64(len(runes)))
	if err != nil {
		return nil, Errorf(err.Error())
	}
	resultRunes := runes[start:stop]
	return NewString(string(resultRunes)), nil
}

func (s *String) SetItem(key, value Object) *Error {
	return Errorf("eval error: string does not support set item")
}

func (s *String) DelItem(key Object) *Error {
	return Errorf("eval error: string does not support del item")
}

func (s *String) Contains(key Object) *Bool {
	strObj, ok := key.(*String)
	if !ok {
		return False
	}
	return NewBool(strings.Contains(s.value, strObj.value))
}

func (s *String) Len() *Int {
	return NewInt(int64(len(s.value)))
}

func (s *String) Iter() Iterator {
	return NewStringIter(s)
}

func (s *String) Runes() []Object {
	runes := []rune(s.value)
	result := make([]Object, len(runes))
	for i, r := range runes {
		result[i] = NewString(string(r))
	}
	return result
}

func NewString(s string) *String {
	return &String{value: s}
}
