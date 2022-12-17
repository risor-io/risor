// The implementation of our regular-expression object.

package object

// Regexp wraps regular-expressions and implements the Object interface.
type Regexp struct {
	// Value holds the string value this object wraps.
	Value string

	// Flags holds the flags for the object
	Flags string
}

func (r *Regexp) Type() Type {
	return REGEXP
}

func (r *Regexp) Inspect() string {
	return r.Value
}

func (r *Regexp) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (r *Regexp) ToInterface() interface{} {
	return "<REGEXP>"
}

func (r *Regexp) Compare(other Object) (int, error) {
	typeComp := CompareTypes(r, other)
	if typeComp != 0 {
		return typeComp, nil
	}
	otherRegex := other.(*Regexp)
	if r.Value == otherRegex.Value {
		if r.Flags == otherRegex.Flags {
			return 0, nil
		}
		if r.Flags > otherRegex.Flags {
			return 1, nil
		}
		return -1, nil
	}
	if r.Value > otherRegex.Value {
		return 1, nil
	}
	return -1, nil
}

func (r *Regexp) Equals(other Object) Object {
	if other.Type() != REGEXP {
		return False
	}
	otherRegex := other.(*Regexp)
	if r.Value == otherRegex.Value && r.Flags == otherRegex.Flags {
		return True
	}
	return False
}
