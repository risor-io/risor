package object

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestObjectString(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
	tests := []struct {
		input    Object
		expected string
	}{
		{True, "true"},
		{False, "false"},
		{Nil, "nil"},
		{NewError(errors.New("kaboom")), "error(kaboom)"},
		{NewFloat(3.0), "3"},
		{NewInt(-3), "-3"},
		{NewString("foo"), "foo"},
		{NewList([]Object{NewInt(1), NewInt(2)}), "[1, 2]"},
		{NewSet([]Object{True, Nil}), "{true, nil}"},
		{NewMap(map[string]Object{"foo": NewInt(1), "bar": NewInt(2)}), `{"bar": 2, "foo": 1}`},
		{NewTime(tm), "time(\"2009-11-10T23:00:00Z\")"},
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

func TestComparisons(t *testing.T) {
	tests := []struct {
		left        Object
		right       Object
		expected    int
		expectedErr error
	}{
		{NewInt(1), NewInt(1), 0, nil},
		{NewInt(1), NewInt(2), -1, nil},
		{NewInt(2), NewInt(1), 1, nil},
		{NewFloat(1.0), NewFloat(1.0), 0, nil},
		{NewFloat(1.0), NewFloat(2.0), -1, nil},
		{NewFloat(2.0), NewFloat(1.0), 1, nil},
		{NewString("a"), NewString("a"), 0, nil},
		{NewString("a"), NewString("b"), -1, nil},
		{NewString("b"), NewString("a"), 1, nil},
		{True, True, 0, nil},
		{True, False, 1, nil},
		{False, True, -1, nil},
		{Nil, Nil, 0, nil},
		{Nil, True, 0, errors.New("type error: unable to compare nil and bool")},
		{NewInt(1), NewFloat(1.0), 0, nil},
		{NewInt(1), NewFloat(2.0), -1, nil},
		{NewInt(1), NewFloat(0.0), 1, nil},
		{NewFloat(1.0), NewInt(1), 0, nil},
		{NewFloat(1.0), NewInt(2), -1, nil},
		{NewFloat(1.0), NewInt(0), 1, nil},
		{NewInt(1), NewString("1"), 0, errors.New("type error: unable to compare int and string")},
		{NewString("1"), NewInt(1), 0, errors.New("type error: unable to compare string and int")},
		{NewFloat(1.0), NewString("1"), 0, errors.New("type error: unable to compare float and string")},
		{NewString("1"), NewFloat(1.0), 0, errors.New("type error: unable to compare string and float")},
		{NewByte(1), NewByte(1), 0, nil},
		{NewByte(1), NewByte(2), -1, nil},
		{NewByte(2), NewByte(1), 1, nil},
		{NewByte(1), NewInt(1), 0, nil},
		{NewByte(1), NewInt(2), -1, nil},
		{NewByte(2), NewInt(1), 1, nil},
		{NewInt(1), NewByte(1), 0, nil},
		{NewInt(1), NewByte(2), -1, nil},
		{NewInt(2), NewByte(1), 1, nil},
		{NewByte(1), NewFloat(1.0), 0, nil},
		{NewByte(1), NewFloat(2.0), -1, nil},
		{NewByte(2), NewFloat(1.0), 1, nil},
		{NewFloat(1.0), NewByte(1), 0, nil},
		{NewFloat(1.0), NewByte(2), -1, nil},
		{NewFloat(1.0), NewByte(0), 1, nil},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s %s", tt.left.Type(), tt.right.Type()), func(t *testing.T) {
			comparable, ok := tt.left.(Comparable)
			require.True(t, ok)
			cmp, cmpErr := comparable.Compare(tt.right)
			require.Equal(t, tt.expected, cmp)
			require.Equal(t, tt.expectedErr, cmpErr)
		})
	}
}
