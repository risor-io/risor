package object

// BreakValue is an implementation detail used to handle break statements
type BreakValue struct{}

func (bv *BreakValue) Type() Type {
	return BREAK_VALUE
}

func (bv *BreakValue) Inspect() string {
	return "break"
}

func (bv *BreakValue) String() string {
	return "break()"
}

func (bv *BreakValue) Interface() interface{} {
	return nil
}

func (bv *BreakValue) Equals(other Object) Object {
	if other.Type() == BREAK_VALUE {
		return True
	}
	return False
}

func (bv *BreakValue) GetAttr(name string) (Object, bool) {
	return nil, false
}
