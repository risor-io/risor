package object

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexpMatch(t *testing.T) {
	// From example: https://pkg.go.dev/regexp#MatchString
	obj := NewRegexp(regexp.MustCompile(`foo.*`))
	match, ok := obj.GetAttr("match")
	require.True(t, ok)
	result := match.(*Builtin).Call(context.Background(), NewString("seafood"))
	require.Equal(t, True, result)

	obj = NewRegexp(regexp.MustCompile(`bar.*`))
	match, ok = obj.GetAttr("match")
	require.True(t, ok)
	result = match.(*Builtin).Call(context.Background(), NewString("seafood"))
	require.Equal(t, False, result)
}

func TestRegexpFind(t *testing.T) {
	// From example: https://pkg.go.dev/regexp#Regexp.Find
	obj := NewRegexp(regexp.MustCompile(`foo.?`))
	find, ok := obj.GetAttr("find")
	require.True(t, ok)
	result := find.(*Builtin).Call(context.Background(), NewString("seafood fool"))
	require.Equal(t, NewString("food"), result)
}

func TestRegexpFindAll(t *testing.T) {
	// From example: https://pkg.go.dev/regexp#Regexp.FindAll
	obj := NewRegexp(regexp.MustCompile(`foo.?`))
	findAll, ok := obj.GetAttr("find_all")
	require.True(t, ok)
	result := findAll.(*Builtin).Call(context.Background(), NewString("seafood fool"))
	require.Equal(t, NewList([]Object{NewString("food"), NewString("fool")}), result)
}
