//go:build !google
// +build !google

package google

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
