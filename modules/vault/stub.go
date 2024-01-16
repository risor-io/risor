//go:build !vault
// +build !vault

package kubernetes

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
