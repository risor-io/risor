package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetScriptArguments(t *testing.T) {

	var test = []struct {
		args      []string
		expOS     []string
		expScript []string
	}{
		{
			[]string{"--foo", "--bar", "1"},
			[]string{"--foo", "--bar", "1"},
			[]string{},
		},
		{
			[]string{"--foo", "1", "--", "/path/to/script"},
			[]string{"/path/to/script"},
			[]string{"/path/to/script"},
		},
		{
			[]string{"--foo", "1", "--"},
			[]string{},
			[]string{},
		},
		{
			[]string{"--", "/path/to/script", "2", "-h", "--bar"},
			[]string{"/path/to/script"},
			[]string{"/path/to/script", "2", "-h", "--bar"},
		},
		{
			[]string{"--"},
			[]string{},
			[]string{},
		},
		{
			[]string{},
			[]string{},
			[]string{},
		},
		{
			[]string{"--no-repl", "1", "--", "/path/to/script", "1", "-f"},
			[]string{"/path/to/script"},
			[]string{"/path/to/script", "1", "-f"},
		},
	}

	for _, tt := range test {
		t.Run(strings.Join(tt.args, "-"), func(t *testing.T) {
			origArgs := os.Args
			testArgs := []string{"risor"}
			testArgs = append(testArgs, tt.args...)
			os.Args = testArgs
			defer func() { os.Args = origArgs }()
			osArgs, scriptArgs := getScriptArgs(tt.args)
			require.Equal(t, tt.expOS, osArgs)
			require.Equal(t, tt.expScript, scriptArgs)
		})
	}
}
