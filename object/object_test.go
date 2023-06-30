package object

import (
	"errors"
	"fmt"
	"testing"

	"time"
)

func TestObjectString(t *testing.T) {

	tm, _ := time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")

	tests := []struct {
		input    Object
		expected string
	}{
		{True, "bool(true)"},
		{False, "bool(false)"},
		{Nil, "nil"},
		// {NewBreak(), "break"},
		// {NewContinue(), "continue"},
		{NewError(errors.New("kaboom")), "error(kaboom)"},
		// {NewReturn(NewInt(42)), "return(int(42))"},
		{NewFloat(3.0), "float(3)"},
		{NewInt(-3), "int(-3)"},
		{NewString("foo"), "string(foo)"},
		// {NewModule("my-scope"), "module(my-scope)"},
		{NewList([]Object{NewInt(1), NewInt(2)}), "list([int(1), int(2)])"},
		{NewSet([]Object{True, Nil}), "set(bool(true), nil)"},
		{NewMap(map[string]Object{"foo": NewInt(1), "bar": NewInt(2)}), `map("bar": int(2), "foo": int(1))`},
		{NewTime(tm), "time(2009-11-10T23:00:00Z)"},
	}

	for _, tt := range tests {
		str, ok := tt.input.(fmt.Stringer)
		if !ok {
			t.Errorf("object.String() not implemented for %T", tt.input)
			continue
		}
		if str.String() != tt.expected {
			t.Errorf("object.String() wrong. want=%q, got=%q", tt.expected, str.String())
		}
	}
}
