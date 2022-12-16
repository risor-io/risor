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
	return fmt.Sprintf("String(%s)", s.Value)
}

func (s *String) HashKey() Key {
	return Key{Type: s.Type(), StrValue: s.Value}
}

func (s *String) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (s *String) InvokeMethod(method string, args ...Object) Object {
	return NewError("type error: %s object has no method %s", s.Type(), method)
}

func (s *String) ToInterface() interface{} {
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

func NewString(s string) *String {
	return &String{Value: s}
}
