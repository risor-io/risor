//go:build !shlex
// +build !shlex

package shlex

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
