//go:build !gh
// +build !gh

package gh

import "github.com/risor-io/risor/object"

func Module() *object.Module {
	return nil
}
