//go:build !aws
// +build !aws

package aws

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return object.NewBuiltinsModule("aws", map[string]object.Object{})
}
