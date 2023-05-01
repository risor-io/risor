package importer

import (
	"context"

	"github.com/cloudcmds/tamarin/object"
)

// Importer is an interface used to import Tamarin code modules
type Importer interface {

	// Import a module by name
	Import(ctx context.Context, name string) (object.Module, error)
}
