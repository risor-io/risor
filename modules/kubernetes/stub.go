//go:build !k8s
// +build !k8s

package kubernetes

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
