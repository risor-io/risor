//go:build !sched
// +build !sched

package sched

import (
	"github.com/risor-io/risor/object"
)

func Module() *object.Module {
	return nil
}
