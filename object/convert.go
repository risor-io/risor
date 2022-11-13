package object

func AsString(obj Object) (result string, err *Error) {
	s, ok := obj.(*String)
	if !ok {
		return "", NewError("type error: expected a string (got %v)", obj.Type())
	}
	return s.Value, nil
}

func AsInteger(obj Object) (int64, *Error) {
	i, ok := obj.(*Integer)
	if !ok {
		return 0, NewError("type error: expected an integer (got %v)", obj.Type())
	}
	return i.Value, nil
}

func AsFloat(obj Object) (float64, *Error) {
	switch obj := obj.(type) {
	case *Integer:
		return float64(obj.Value), nil
	case *Float:
		return obj.Value, nil
	default:
		return 0.0, NewError("type error: expected a number (got %v)", obj.Type())
	}
}

func AsArray(obj Object) (*Array, *Error) {
	arr, ok := obj.(*Array)
	if !ok {
		return nil, NewError("type error: expected an array (got %v)", obj.Type())
	}
	return arr, nil
}
