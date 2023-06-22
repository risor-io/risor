//go:build !aws
// +build !aws

package aws

import (
	"github.com/cloudcmds/tamarin/v2/object"
)

func Module() *object.Module {
	return object.NewBuiltinsModule("aws", map[string]object.Object{})
}
