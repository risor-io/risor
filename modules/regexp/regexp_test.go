package regexp

import (
	"context"
	"regexp"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestRegexpMatch(t *testing.T) {
	// From example: https://pkg.go.dev/regexp#MatchString
	obj := NewRegexp(regexp.MustCompile(`foo.*`))
	match, ok := obj.GetAttr("match")
	require.True(t, ok)
	result := match.(*object.Builtin).Call(context.Background(), object.NewString("seafood"))
	require.Equal(t, object.True, result)

	obj = NewRegexp(regexp.MustCompile(`bar.*`))
	match, ok = obj.GetAttr("match")
	require.True(t, ok)
	result = match.(*object.Builtin).Call(context.Background(), object.NewString("seafood"))
	require.Equal(t, object.False, result)
}

func TestRegexpFind(t *testing.T) {
	// From example: https://pkg.go.dev/regexp#Regexp.Find
	obj := NewRegexp(regexp.MustCompile(`foo.?`))
	find, ok := obj.GetAttr("find")
	require.True(t, ok)
	result := find.(*object.Builtin).Call(context.Background(), object.NewString("seafood fool"))
	require.Equal(t, object.NewString("food"), result)
}

func TestRegexpFindAll(t *testing.T) {
	// From example: https://pkg.go.dev/regexp#Regexp.FindAll
	obj := NewRegexp(regexp.MustCompile(`foo.?`))
	findAll, ok := obj.GetAttr("find_all")
	require.True(t, ok)
	result := findAll.(*object.Builtin).Call(context.Background(), object.NewString("seafood fool"))
	require.Equal(t, object.NewList([]object.Object{
		object.NewString("food"),
		object.NewString("fool"),
	}), result)
}
