//go:build !vault
// +build !vault

package vault

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
