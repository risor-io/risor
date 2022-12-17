package object

// ReturnValue wraps Object and implements Object interface.
type ReturnValue struct {
	// Value is the object that is to be returned
	Value Object
}

func (rv *ReturnValue) Type() Type {
	return RETURN_VALUE
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

func (rv *ReturnValue) GetAttr(name string) (Object, bool) {
	return nil, false
}

func (rv *ReturnValue) ToInterface() interface{} {
	return nil
}

func (rv *ReturnValue) Equals(other Object) Object {
	if other.Type() == RETURN_VALUE && rv == other.(*ReturnValue) {
		return True
	}
	return False
}
