//go:build !semver
// +build !semver

package semver

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
