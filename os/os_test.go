package os

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanPath(t *testing.T) {
	type test struct {
		input string
		want  string
	}
	tests := []test{
		{input: "..///foo", want: "../foo"},
		{input: "../../foo", want: "../../foo"},
		{input: "../foo/../bar.txt", want: "../bar.txt"},
		{input: "foo/././bar.txt", want: "foo/bar.txt"},
		{input: "foo/../../bar.txt", want: "../bar.txt"},
	}
	for _, tc := range tests {
		got := filepath.Clean(tc.input)
		assert.Equal(t, tc.want, got)
	}
}

func TestResolvePath(t *testing.T) {
	newError := func(msg string) error {
		return &fs.PathError{
			Op:   "open",
			Path: msg,
			Err:  fs.ErrInvalid,
		}
	}
	type test struct {
		base string
		path string
		want interface{}
	}
	tests := []test{
		{base: "/", path: "..///foo", want: newError("../foo")},
		{base: "/", path: "../../foo", want: newError("../../foo")},
		{base: "/", path: "../foo/../bar.txt", want: newError("../bar.txt")},
		{base: "/", path: "/foo/././bar.txt", want: "/foo/bar.txt"},
		{base: "/", path: "/foo/../bar.txt", want: "/bar.txt"},
		{base: "/", path: "/foo/../../bar.txt", want: "/bar.txt"},
		{base: "/dir", path: "/foo/../../bar.txt", want: "/dir/bar.txt"},
		{base: "/dir/", path: "bar.txt", want: "/dir/bar.txt"},
	}
	for _, tc := range tests {
		got, err := ResolvePath(tc.base, tc.path, "open")
		switch tc.want.(type) {
		case string:
			assert.Equal(t, tc.want, got, err)
		case error:
			assert.Equal(t, tc.want, err, got)
		}
	}
}
