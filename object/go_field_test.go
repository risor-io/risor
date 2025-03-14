package object

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Enabled bool   `json:"enabled"`
}

func TestGoField(t *testing.T) {
	// Create a test instance and get its type
	s := &testStruct{}
	typ := reflect.TypeOf(s).Elem()

	// Test string field
	nameField, ok := typ.FieldByName("Name")
	require.True(t, ok)
	field, err := newGoField(nameField)
	require.Nil(t, err)

	require.Equal(t, "Name", field.Name())
	require.Equal(t, "string", field.GoType().Name())
	require.Equal(t, `json:"name"`, string(field.Tag()))

	// Test int field
	ageField, ok := typ.FieldByName("Age")
	require.True(t, ok)
	field, err = newGoField(ageField)
	require.Nil(t, err)

	require.Equal(t, "Age", field.Name())
	require.Equal(t, "int", field.GoType().Name())
	require.Equal(t, `json:"age"`, string(field.Tag()))

	// Test bool field
	enabledField, ok := typ.FieldByName("Enabled")
	require.True(t, ok)
	field, err = newGoField(enabledField)
	require.Nil(t, err)

	require.Equal(t, "Enabled", field.Name())
	require.Equal(t, "bool", field.GoType().Name())
	require.Equal(t, `json:"enabled"`, string(field.Tag()))
}

type complexStruct struct {
	Data    map[string]interface{} `json:"data"`
	Numbers []int                  `json:"numbers"`
	Ptr     *string                `json:"ptr"`
}

func TestGoFieldComplexTypes(t *testing.T) {
	s := complexStruct{}
	typ := reflect.TypeOf(s)

	// Test map field
	dataField, ok := typ.FieldByName("Data")
	require.True(t, ok)
	field, err := newGoField(dataField)
	require.Nil(t, err)

	require.Equal(t, "Data", field.Name())
	require.Equal(t, "map[string]interface {}", field.GoType().Name())
	require.Equal(t, `json:"data"`, string(field.Tag()))

	// Test slice field
	numbersField, ok := typ.FieldByName("Numbers")
	require.True(t, ok)
	field, err = newGoField(numbersField)
	require.Nil(t, err)

	require.Equal(t, "Numbers", field.Name())
	require.Equal(t, "[]int", field.GoType().Name())
	require.Equal(t, `json:"numbers"`, string(field.Tag()))

	// Test pointer field
	ptrField, ok := typ.FieldByName("Ptr")
	require.True(t, ok)
	field, err = newGoField(ptrField)
	require.Nil(t, err)

	require.Equal(t, "Ptr", field.Name())
	require.Equal(t, "*string", field.GoType().Name())
	require.Equal(t, `json:"ptr"`, string(field.Tag()))
}

type nestedStruct struct {
	Inner struct {
		Value string `json:"value"`
	} `json:"inner"`
}

func TestGoFieldNestedStruct(t *testing.T) {
	s := &nestedStruct{}
	typ := reflect.TypeOf(s).Elem()

	// Test nested struct field
	innerField, ok := typ.FieldByName("Inner")
	require.True(t, ok)
	field, err := newGoField(innerField)
	require.Nil(t, err)

	require.Equal(t, "Inner", field.Name())
	require.Equal(t, "*struct { Value string \"json:\\\"value\\\"\" }", field.GoType().Name())
	require.Equal(t, `json:"inner"`, string(field.Tag()))
}

func TestGoFieldGetAttr(t *testing.T) {
	s := &testStruct{}
	typ := reflect.TypeOf(s).Elem()
	nameField, ok := typ.FieldByName("Name")
	require.True(t, ok)
	field, err := newGoField(nameField)
	require.Nil(t, err)

	// Test GetAttr for name
	nameAttr, ok := field.GetAttr("name")
	require.True(t, ok)
	require.Equal(t, "Name", nameAttr.(*String).value)

	// Test GetAttr for type
	typeAttr, ok := field.GetAttr("type")
	require.True(t, ok)
	require.Equal(t, "string", typeAttr.(*GoType).Name())

	// Test GetAttr for tag
	tagAttr, ok := field.GetAttr("tag")
	require.True(t, ok)
	require.Equal(t, `json:"name"`, tagAttr.(*String).value)

	// Test GetAttr for non-existent attribute
	_, ok = field.GetAttr("nonexistent")
	require.False(t, ok)
}
